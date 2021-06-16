package homeassistant

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/insomniadev/martian/integrations/config"
)

func (h *HomeAssistant) Init() {
	h.connect()
}

func (h *HomeAssistant) connect() {
	// host := "ws://" + h.Url
	h.Config = config.LoadHomeAssistant()
	host := "ws://" + h.Config.URL + "/api/websocket"

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
		if message.Type == "auth_required" {
			authMessage := AuthEvent{Type: "auth", AccessToken: h.Config.Token}
			authEvent, err := json.Marshal(authMessage)
			println(string(authEvent))
			if err != nil {
				println(err)
			}
			h.Connection.WriteMessage(1, authEvent)
		} else {
			println(string(incoming))
			return
		}
	}
}
