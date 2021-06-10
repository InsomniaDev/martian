package homeassistant

import "github.com/gorilla/websocket"

type HomeAssistant struct {
	Url        string
	Token      string
	Connection *websocket.Conn
}

type Event struct {
	ID    int64      `json:"id"`
	Type  string     `json:"type"`
	Event EventClass `json:"event"`
}

type EventClass struct {
	Data      Data    `json:"data"`
	EventType string  `json:"event_type"`
	TimeFired string  `json:"time_fired"`
	Origin    string  `json:"origin"`
	Context   Context `json:"context"`
}

type Context struct {
	ID       string      `json:"id"`
	ParentID interface{} `json:"parent_id"`
	UserID   string      `json:"user_id"`
}

type Data struct {
	EntityID string   `json:"entity_id"`
	NewState NewState `json:"new_state"`
	OldState OldState `json:"old_state"`
}

type NewState struct {
	EntityID    string             `json:"entity_id"`
	LastChanged string             `json:"last_changed"`
	State       string             `json:"state"`
	Attributes  NewStateAttributes `json:"attributes"`
	LastUpdated string             `json:"last_updated"`
	Context     Context            `json:"context"`
}

type NewStateAttributes struct {
	RGBColor          []int64   `json:"rgb_color"`
	ColorTemp         int64     `json:"color_temp"`
	SupportedFeatures int64     `json:"supported_features"`
	XyColor           []float64 `json:"xy_color"`
	Brightness        int64     `json:"brightness"`
	WhiteValue        int64     `json:"white_value"`
	FriendlyName      string    `json:"friendly_name"`
}

type OldState struct {
	EntityID    string             `json:"entity_id"`
	LastChanged string             `json:"last_changed"`
	State       string             `json:"state"`
	Attributes  OldStateAttributes `json:"attributes"`
	LastUpdated string             `json:"last_updated"`
	Context     Context            `json:"context"`
}

type OldStateAttributes struct {
	SupportedFeatures int64  `json:"supported_features"`
	FriendlyName      string `json:"friendly_name"`
}
