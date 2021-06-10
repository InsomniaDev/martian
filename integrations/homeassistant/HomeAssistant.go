package homeassistant

import (
	"log"

	"github.com/gorilla/websocket"
)

func (h *HomeAssistant) Init() {
	h.connect()
}

func (h *HomeAssistant) connect() {
	// host := "ws://" + h.Url
	host := "ws://" + "192.168.1.19:8123"

	conn, _, err := websocket.DefaultDialer.Dial(host, nil)
	h.Connection = conn
	if err != nil {
		log.Fatal("homeassistant dial:", err)
	}
	go h.listen()
}

func (h *HomeAssistant) listen() {
	for {
		_, message, err := h.Connection.ReadMessage()
		if err != nil {
			println(err)
		}
		println(string(message))
	}
}
