package harmony

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/insomniadev/martian/database"
	"github.com/insomniadev/martian/modules/pubsub"
)

var (
	defaultHubPort = 8088
)

// Init the harmony hub communication
func (d *Device) Init(configuration string) error {
	// Load up the configuration as the device itself
	json.Unmarshal([]byte(configuration), d)

	// IF the IPAddress doesn't exist, then we need to find the endpoint
	if d.IPAddress == "" {
		d.discover()
	}

	d.connect()
	err := database.MartianData.PutIntegrationValue("harmony", d)
	if err != nil {
		return err
	}
	return nil
}

// discover will find the harmony device on the network
func (d *Device) discover() {
	addr := "192.168.1.1/24"
	ip, ipnet, err := net.ParseCIDR(addr)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		wg.Add(1)
		go func(ipdad string) {
			defer wg.Done()
			message, err := checkIpAddress(ipdad)
			if err != nil {
				return
			}
			d.IPAddress = ipdad + ":" + strconv.Itoa(defaultHubPort)
			d.ActiveRemoteID = message.Data.ActiveRemoteID
			u, err := url.Parse(message.Data.DiscoveryServer)
			if err != nil {
				log.Fatal(err)
			}
			d.HostName = u.Hostname()

		}(ip.String())
	}
	wg.Wait()
}
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
func checkIpAddress(ipAddress string) (HTTPMessage, error) {
	host := "http://" + ipAddress + ":" + strconv.Itoa(defaultHubPort)
	jsonBody := []byte(`{"id":1,"cmd":"setup.account?getProvisionInfo","params":{}}`)
	client := &http.Client{Timeout: 500 * time.Millisecond}
	req, err := http.NewRequest("Post", host, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Origin", "http://sl.dhg.myharmony.com")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Charset", "utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return HTTPMessage{}, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	message := HTTPMessage{}
	if err = json.Unmarshal(body, &message); err != nil {
		log.Fatal(err)
		return HTTPMessage{}, err
	}
	return message, nil
}

func (d *Device) connect() {
	host := "ws://" + d.IPAddress + "/?domain=" + d.HostName + "&hubId=" + strconv.Itoa(d.ActiveRemoteID)

	conn, _, err := websocket.DefaultDialer.Dial(host, nil)
	d.Connection = conn
	if err != nil {
		log.Fatal("harmony dial:", err)
	}
	go d.listen()
	d.GetConfig()
	d.GetCurrentActivity()
}

// heartbeat sends a websocket request every 30 seconds so that the websocket connection stays open
func (d *Device) heartbeat() {
	for {
		time.Sleep(30 * time.Second)
		d.WriteMessage("")
	}
}

func (d *Device) listen() {
	go d.heartbeat()
	for {
		time.Sleep(1 * time.Second)
		_, message, err := d.Connection.ReadMessage()
		if err != nil {
			log.Println("harmony read:", err.Error())
			time.Sleep(30 * time.Second) // Wait for thirty seconds
			go d.connect()               // Start a new process to connect
			return                       // exit out of this current listening loop
		}
		if len(string(message)) > 0 {
			receivedMessage := RecvMessage{}
			err = json.Unmarshal(message, &receivedMessage)
			if receivedMessage.ID == 1 {
				receivedResult := RecvResult{}
				err = json.Unmarshal(message, &receivedResult)
				d.CurrentActivity = receivedResult.Data.Result
				d.ActivityID = receivedResult.Data.Result
			} else if receivedMessage.ID == 2 {
				receivedResult := RecvConfig{}
				err = json.Unmarshal(message, &receivedResult)
				activitiesRecorded := make([]Activity, 0)
				for _, data := range receivedResult.Data.Activity {
					deviceActions := make([]ActivityActions, 0)
					for _, device := range data.ControlGroup {
						for _, funct := range device.Function {
							actCommand := ActivityActionCommand{
								Command:  funct.Action,
								DeviceID: funct.Name,
								Type:     funct.Label,
							}
							actAction := ActivityActions{
								Label:  device.Name,
								Action: actCommand,
							}
							deviceActions = append(deviceActions, actAction)
						}
					}
					nextActivity := Activity{
						ActivityID: data.ID,
						Name:       data.Label,
						Actions:    deviceActions,
					}
					activitiesRecorded = append(activitiesRecorded, nextActivity)
				}
				d.Activities = activitiesRecorded

				database.MartianData.PutIntegrationValue("harmony", d)
			} else if receivedMessage.Type == "connect.stateDigest?notify" {
				receivedResult := StateDigestNotify{}
				err = json.Unmarshal(message, &receivedResult)
				d.CurrentActivity = receivedResult.Data.ActivityID
				// FIXME: Fix the brain notification
				formattedString := "harmony;;harmony;;" + d.CurrentActivity
				pubsub.Service.Publish("brain", formattedString)
				pubsub.Service.Publish("subscriptions", d.CurrentActivity)
			}
		}
	}
}

// GetCurrentActivity calls the websocket to return the current activity
func (d *Device) GetCurrentActivity() {
	message := `{"hubId":"` + strconv.Itoa(d.ActiveRemoteID) + `","timeout":30, "hbus":{"cmd":"vnd.logitech.harmony/vnd.logitech.harmony.engine?getCurrentActivity","id":1,"params":{"verb":"get","format":"json"}}}`
	d.WriteMessage(message)
}

// GetConfig calls the websocket to return the hub configuration
func (d *Device) GetConfig() {
	message := `{"hubId":"` + strconv.Itoa(d.ActiveRemoteID) + `","timeout":30, "hbus":{"cmd":"vnd.logitech.harmony/vnd.logitech.harmony.engine?config","id":2,"params":{"verb":"get","format":"json"}}}`
	d.WriteMessage(message)
}

// StartActivity starts an activity
func (d *Device) StartActivity(activityID string) {
	message := `{"hubId":"` + strconv.Itoa(d.ActiveRemoteID) + `","timeout":30, "hbus":{"cmd":"harmony.activityengine?runactivity","id":3,"params":{"async":"true","timestamp":"0","args":{"rule":"start"},"activityId":"` + activityID + `"}}}`
	d.WriteMessage(message)
}

// SendCommand sends a command to the harmony hub
func (d *Device) SendCommand(command string, commandType string, deviceID string) {
	nextCommand := &Command{
		Command:  command,
		Type:     commandType,
		DeviceID: deviceID,
	}
	d.sendCommandToHub(*nextCommand)
}

// sendCommandToHub sends a command to the harmony hub
func (d *Device) sendCommandToHub(action Command) {
	actionJSON, err := json.Marshal(action)
	if err != nil {
		log.Println("Harmony Command Error:", err)
	}
	message := `{"hubId":` + strconv.Itoa(d.ActiveRemoteID) + `,"timeout":30, "hbus":{"cmd":"vnd.logitech.harmony/vnd.logitech.harmony.engine?holdAction","id":4,"params":{"status":"pressrelease","timestamp":"0","verb":"render","action":"` + string(actionJSON) + `"}}}`
	d.WriteMessage(message)
}

// WriteMessage will write a message to the harmony websocket
func (d *Device) WriteMessage(message string) {
	err := d.Connection.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("harmony write:", err)
	}
}

// EditDeviceConfiguration will update the name and areaname for the harmony configuration
func (d *Device) EditDeviceConfiguration(device Device, removeEdit bool) {
	d.Name = device.Name
	d.AreaName = device.AreaName

	// Save in the database
	err := database.MartianData.PutIntegrationValue("harmony", d)
	if err != nil {
		log.Println(err)
	}
}
