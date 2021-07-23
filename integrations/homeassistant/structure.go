package homeassistant

import (
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
)

type HomeAssistant struct {
	Url             string `json:"url"`
	Token           string `json:"token"`
	Connection      *websocket.Conn
	Devices         []HomeAssistantDevice `json:"devices"`
	SelectedDevices []HomeAssistantDevice `json:"selectedDevices"`
}

// GraphqlHomeAssistantType is the graphql object for the HomeAssistant integration
var GraphqlHomeAssistantType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "HomeAssistantType",
		Fields: graphql.Fields{
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"token": &graphql.Field{
				Type: graphql.String,
			},
			"devices": &graphql.Field{
				Type: graphql.NewList(GraphqlHomeAssistantDeviceType),
			},
			"selectedDevices": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

type HomeAssistantDevice struct {
	EntityId string `json:"entityId"`
	Name     string `json:"name"`
	Group    string `json:"group"`
	Type     string `json:"type"`
	State    string `json:"state"`
	AreaName string `json:"areaName"`
}

// GraphqlHomeAssistantDeviceType is the graphql object for the HomeAssistantDevice type
var GraphqlHomeAssistantDeviceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "HomeAssistantDeviceType",
		Fields: graphql.Fields{
			"entityId": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"group": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"state": &graphql.Field{
				Type: graphql.String,
			},
			"areaName": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

type AuthEvent struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

type Event struct {
	ID      int        `json:"id"`
	Type    string     `json:"type"`
	Success bool       `json:"success"`
	Event   EventClass `json:"event"`
	Result  []Results  `json:"result"`
}

type Results struct {
	EntityId    string     `json:"entity_id"`
	State       string     `json:"state"`
	Attributes  Attributes `json:"attributes"`
	LastChanged string     `json:"last_changed"`
	LastUpdated string     `json:"last_updated"`
}

type Attributes struct {
	FriendlyName string `json:"friendly_name"`
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
	RGBColor          []int     `json:"rgb_color"`
	ColorTemp         int       `json:"color_temp"`
	SupportedFeatures int       `json:"supported_features"`
	XyColor           []float64 `json:"xy_color"`
	Brightness        int       `json:"brightness"`
	WhiteValue        int       `json:"white_value"`
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
	SupportedFeatures int    `json:"supported_features"`
	FriendlyName      string `json:"friendly_name"`
}
