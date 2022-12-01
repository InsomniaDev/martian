package hubitat

// URLs

// get data
// http://192.168.1.131/apps/api/35/devices/all?access_token=bfd56b33-a58b-4fcd-afc8-0d72e41525ac
// http://192.168.1.131/apps/api/35/devices/37?access_token=bfd56b33-a58b-4fcd-afc8-0d72e41525ac
// http://192.168.1.131/apps/api/35/devices?access_token=bfd56b33-a58b-4fcd-afc8-0d72e41525ac

// send hubitat post to this computer
// http://192.168.1.131/apps/api/35/postURL/http%3A%2F%2F192.168.1.52%3A8088%2Fhubitat?access_token=bfd56b33-a58b-4fcd-afc8-0d72e41525ac

// send hubitat post to the server
// http://192.168.1.131/apps/api/35/postURL/http%3A%2F%2F192.168.1.19%3A30919%2Fhubitat?access_token=bfd56b33-a58b-4fcd-afc8-0d72e41525ac

// update data
// curl http://192.168.1.131/apps/api/35/devices/37/off?access_token=bfd56b33-a58b-4fcd-afc8-0d72e41525ac

// mutation createHubitat {
//   updateIntegration(type: "hubitat", value: "{\"url\": \"http://192.168.1.131\",\"accessKey\": \"bfd56b33-a58b-4fcd-afc8-0d72e41525ac\"}")
// }

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/insomniadev/martian/pkg/pubsub"
	"github.com/spf13/viper"
)

var Instance Hubitat

const (
	appUrlString = "/apps/api/35/devices"
)

// pulse should consantly go through and update the state for the hubitat integration whenever the brain requests it
func pulse() {
	subscriptionBus := make(chan string)
	var subscribeToEvents = func() {
		for {
			msg := <-subscriptionBus
			log.Debugln(msg, " pulse")
			Instance.GetAllDeviceStatus()
		}
	}
	pubsub.Service.Subscribe("pulse", subscriptionBus)
	go subscribeToEvents()
}

func init() {
	subscriptionBus := make(chan string)
	var subscribeToEvents = func() {
		for {
			msg := <-subscriptionBus
			log.Debug("Received message: ", msg)

			// deviceId ;; changetostatus
			message := strings.Split(msg, ";;")
			deviceId, _ := strconv.Atoi(message[0])

			switch message[1] {
			case "energy":
				// Currently hardcode to off rather than being smart
				Instance.updateDeviceStatus(deviceId, "off")
			case "activate":
				// time activated so turn the device on
				Instance.updateDeviceStatus(deviceId, "on")
			default:
				Instance.updateDeviceStatus(deviceId, message[1])

			}
		}
	}
	pubsub.Service.Subscribe("hubitat", subscriptionBus)
	go subscribeToEvents()
	pulse()
}

// GetAllDeviceStatus will reach out to hubitat and get all of the device status'
func (h *Hubitat) GetAllDeviceStatus() (contents []byte, err error) {

	viper.UnmarshalKey("hubitat", &Instance)

	requestUrl := h.HubitatUrl + "all?access_token=" + h.AccessKey

	client := &http.Client{}
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)

	if contents, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Fatal(err)
	}
	var contentResponse HubitatDeviceAllResponse
	if err = json.Unmarshal(contents, &contentResponse); err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	for _, state := range contentResponse {

		capabilities, _ := json.Marshal(state.Capabilities)
		attributes, _ := json.Marshal(state.Attributes)
		commands, _ := json.Marshal(state.Commands)

		msg := fmt.Sprintf("%s;;%s;;%s;;%s;;%s;;%s;;%s;;%s;;%s;;%s",
			state.ID,
			state.Name,
			state.Label,
			state.Type,
			state.Model,
			state.Manufacturer,
			string(capabilities),
			string(attributes),
			string(commands),
			"hubitat",
		)

		pubsub.Service.Publish("learnNewDevice", msg)
	}

	for _, newState := range contentResponse {
		activated := false
		switch newState.Attributes.Switch {
		case "on":
			activated = true
		}

		// hubitat : unique ID : label : attribute type : attribute value
		updateString := fmt.Sprintf("hubitat;;%s;;%s;;%s;;%s;;%t", newState.ID, newState.Label, newState.Name, newState.Attributes.Switch, activated)
		pubsub.Service.Publish("brain", updateString)
	}

	return
}

// GetDatas will take the incoming data, parse, and send to the brain
// example message:
// 	{
//     "content": {
//         "name": "switch",
//         "value": "off",
//         "displayName": "Study",
//         "deviceId": "37",
//         "descriptionText": "Study was turned off",
//         "unit": null,
//         "type": "digital",
//         "data": null
//     }
// }

func (h *Hubitat) GetDatas(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("failed to read, ", err)
	}
	var response HubitatUpdate
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatal("hubitat ::: failed to unmarshal update, ", err)
	}

	activated := false
	switch response.Content.Value {
	case "on":
		activated = true
	}

	// Let's update the previous states with the newly updated entry
	for i := range Instance.PreviousStates {
		if Instance.PreviousStates[i].ID == response.Content.DeviceID {
			// TODO: This only currently supports the switch type, we need to update this in the future
			Instance.PreviousStates[i].Attributes.Switch = response.Content.Value
		}
	}

	// hubitat : unique ID : label : attribute type : attribute value
	pubsub.Service.Publish("brain", fmt.Sprintf("hubitat;;%s;;%s;;%s;;%s;;%t", response.Content.DeviceID, response.Content.DisplayName, response.Content.Name, response.Content.Value, activated))
}

// updateDeviceStatus unused function to update a status
func (h *Hubitat) updateDeviceStatus(deviceId int, command string) (contents []byte, err error) {
	requestUrl := h.HubitatUrl + "/" + strconv.Itoa(deviceId) + "/" + command + "?access_token=" + h.AccessKey

	client := &http.Client{}
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)

	contents, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return
}

// have a go routine constantly running and pinging for updates from hubitat
// create an automation file to use for the time being that will create an automation around a particular device
// have a go routine spin off when an automation is completed that will then check for completion
// will need to check that this automation completion isn't already going (we don't want duplicates)
