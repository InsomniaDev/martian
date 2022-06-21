package hubitat

// URLs

// get data
// http://192.168.1.131/apps/api/35/devices/all?access_token=bfd56b33-a58b-4fcd-afc8-0d72e41525ac
// http://192.168.1.131/apps/api/35/devices/37?access_token=bfd56b33-a58b-4fcd-afc8-0d72e41525ac
// http://192.168.1.131/apps/api/35/devices?access_token=bfd56b33-a58b-4fcd-afc8-0d72e41525ac

// update data
// curl http://192.168.1.131/apps/api/35/devices/37/off?access_token=bfd56b33-a58b-4fcd-afc8-0d72e41525ac

// mutation createHubitat {
//   updateIntegration(type: "hubitat", value: "{\"url\": \"http://192.168.1.131\",\"accessKey\": \"bfd56b33-a58b-4fcd-afc8-0d72e41525ac\"}")
// }

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/insomniadev/martian/integrations/config"
)

var Hubitat HubitatData

const (
	appUrlString = "/apps/api/35/devices"
)

func init() {
	Hubitat.Config.AccessKey, Hubitat.Config.HubitatUrl = config.LoadHubitat()
}

func (h *HubitatData) Init(configuration string) {
	err := json.Unmarshal([]byte(configuration), &h.Config)
	if err != nil {
		log.Println(err)
	}

	// cycle through and constantly check for any updates to the devices
	go func() {
		for {
			// check every second for an update
			time.Sleep(1 * time.Second)
			h.getAllDeviceStatus()
		}
	}()
	// run the configuration piece here
}

func (h *HubitatData) getAllDeviceStatus() (contents []byte, err error) {

	requestUrl := h.Config.HubitatUrl + "all?access_token=" + h.Config.AccessKey

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

func (h *HubitatData) updateDeviceStatus(deviceId int, command string) (contents []byte, err error) {
	requestUrl := h.Config.HubitatUrl + "/" + strconv.Itoa(deviceId) + "/" + command + "?access_token=" + h.Config.AccessKey

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
