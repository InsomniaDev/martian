package homeassistant

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/insomniadev/martian/modules/redispub"
)

var (
	subscribeEventsId   int    = 1
	subscribeEventsType string = "subscribe_events"

	getStatesId   int    = 2
	getStatesType string = "get_states"

	getServicesId   int    = 3
	getServicesType string = "get_services"

	getPanelsId   int    = 4
	getPanelsType string = "get_panels"

	getConfigId   int    = 5
	getConfigType string = "get_config"

	getMediaPlayerThumbnailId   int    = 6
	getMediaPlayerThumbnailType string = "media_player_thumbnail"

	getCameraThumbnailId   int    = 7
	getCameraThumbnailType string = "camera_thumbnail"

	callServiceId int = 9
)

func (h *HomeAssistant) Init(configuration string) error {
	err := json.Unmarshal([]byte(configuration), &h)
	if err != nil {
		return err
	}
	h.connect()
	return nil
}

func (h *HomeAssistant) connect() {
	// host := "ws://" + h.Url
	host := "ws://" + h.Url + "/api/websocket"

	conn, _, err := websocket.DefaultDialer.Dial(host, nil)
	h.Connection = conn
	if err != nil {
		log.Fatal("homeassistant dial:", err)
	}
	go h.listen()
}

func (h *HomeAssistant) listen() {
	for {
		_, incoming, err := h.Connection.ReadMessage()
		if err != nil {
			println(err)
		}
		var message Event
		err = json.Unmarshal(incoming, &message)
		if err != nil {
			println(err)
		}
		switch message.Type {
		case "auth_required":
			authMessage := AuthEvent{Type: "auth", AccessToken: h.Token}
			authEvent, err := json.Marshal(authMessage)
			println(string(authEvent))
			if err != nil {
				println(err)
			}
			h.Connection.WriteMessage(1, authEvent)
		case "auth_ok":
			h.subscribeEvents()
			// h.getConfig()
			// h.getServices()
			h.getStates()
		case "event":
			for i := range h.Devices {
				if h.Devices[i].EntityId == message.Event.Data.EntityID && h.Devices[i].State != message.Event.Data.NewState.State {
					h.Devices[i].State = message.Event.Data.NewState.State

					redispub.Service.Publish("subscriptions", h.Devices[i].EntityId)
				}
			}
			// default:
			// 	println(string(incoming))
		}
		switch message.ID {
		case getStatesId:
			for _, result := range message.Result {
				s := strings.Split(result.EntityId, ".")
				deviceType, name := s[0], s[1]
				name = strings.Replace(name, "_", " ", -1)
				areaName := ""
				friendlyName := strings.Split(result.Attributes.FriendlyName, "_")
				if len(friendlyName) > 1 {
					areaName = friendlyName[0]
				}
				newDevice := HomeAssistantDevice{EntityId: result.EntityId, Name: name, Type: deviceType, State: result.State, AreaName: areaName}
				h.Devices = append(h.Devices, newDevice)
			}
			// for _, dev := range h.Devices {
			// 	if dev.Type == "light" {
			// 		fmt.Println(dev.Name)
			// 	}
			// }
		}
	}
}

// CallService will call a service and update the value
func (h *HomeAssistant) CallService(device HomeAssistantDevice, activate bool) {
	setValue := "turn_off"
	if activate {
		setValue = "turn_on"
	}
	serviceJson := `{"id":` + strconv.Itoa(callServiceId) + `,"type":"call_service","domain":"` + device.Type + `","service":"` + setValue + `","service_data":{"entity_id":"` + device.EntityId + `"}}`

	// TODO: We don't need this print statement here once we are done with Hass implementation
	println(serviceJson)

	err := h.Connection.WriteMessage(websocket.TextMessage, []byte(serviceJson))
	if err != nil {
		log.Println("hass write:", err)
	}
	callServiceId = callServiceId + 1
}

// subscribeEvents will subscribe to the events from the home assistant websocket
func (h *HomeAssistant) subscribeEvents() {
	subscription := `{"id":` + strconv.Itoa(subscribeEventsId) + `,"type":"` + subscribeEventsType + `","event_type":"state_changed"}`

	println(subscription)
	err := h.Connection.WriteMessage(websocket.TextMessage, []byte(subscription))
	if err != nil {
		log.Println("hass write:", err)
	}
}

// getStates will get the states of all of the devices used by home assistant
func (h *HomeAssistant) getStates() {
	States := `{"id":` + strconv.Itoa(getStatesId) + `,"type":"` + getStatesType + `"}`

	println(States)
	err := h.Connection.WriteMessage(websocket.TextMessage, []byte(States))
	if err != nil {
		log.Println("hass write:", err)
	}
}

func (h *HomeAssistant) getConfig() {
	config := `{"id":` + strconv.Itoa(getConfigId) + `,"type":"` + getConfigType + `"}`

	println(config)
	err := h.Connection.WriteMessage(websocket.TextMessage, []byte(config))
	if err != nil {
		log.Println("hass write:", err)
	}
}

func (h *HomeAssistant) getServices() {
	Services := `{"id":` + strconv.Itoa(getServicesId) + `,"type":"` + getServicesType + `"}`

	println(Services)
	err := h.Connection.WriteMessage(websocket.TextMessage, []byte(Services))
	if err != nil {
		log.Println("hass write:", err)
	}
}

// UpdateSelectedDevices will go through and update the selected 
func (h *HomeAssistant) UpdateSelectedDevices(selectedDevices []string) error {
	h.SelectedDevices = []HomeAssistantDevice{}
	for _, deviceName := range selectedDevices {
		found := false
		for i := range h.Devices {
			if h.Devices[i].Name == deviceName {
				found = true
				h.SelectedDevices = append(h.SelectedDevices, h.Devices[i])
			}
		}
		if !found {
			return errors.New(deviceName + " was not found in the integrated Hass devices")
		}
	}
	return nil
}
