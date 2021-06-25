package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/insomniadev/martian/integrations/homeassistant"
)

//Lutron struct
type Lutron struct {
	Name     string  `json:"name"`
	DeviceID int     `json:"id"`
	AreaName string  `json:"areaName"`
	Type     string  `json:"type"`
	Value    float64 `json:"value"`
	State    string  `json:"state"`
}

var lutronType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Lutron",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"id": &graphql.Field{
				Type: graphql.ID,
			},
			"areaName": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"value": &graphql.Field{
				Type: graphql.Float,
			},
			"state": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// Harmony struct for changing into the graphql response
type Harmony struct {
	ID         string `json:"id"`
	ActivityID string `json:"activityId"`
	Name       string `json:"name"`
	Actions    []struct {
		Label  string `json:"label"`
		Action struct {
			Command  string `json:"command"`
			Type     string `json:"type"`
			DeviceID string `json:"deviceId"`
		} `json:"action"`
	} `json:"actions"`
}

var hassType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Hass",
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
		},
	},
)

var harmonyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Harmony",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"activityId": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"actions": &graphql.Field{
				Type: graphql.NewList(harmonyActionType),
			},
		},
	},
)
var harmonyActionType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "HarmonyAction",
		Fields: graphql.Fields{
			"label": &graphql.Field{
				Type: graphql.String,
			},
			"action": &graphql.Field{
				Type: harmonyCommandType,
			},
		},
	},
)
var harmonyCommandType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "HarmonyCommand",
		Fields: graphql.Fields{
			"command": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"deviceId": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// Kasa is the struct used for GraphQL
type Kasa struct {
	ID        string `json:"id"`
	IPAddress string `json:"ipAddress"`
	Name      string `json:"name"`
	AreaName  string `json:"areaName"`
	IsOn      bool   `json:"on"`
}

var kasaType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Kasa",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"ipAddress": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"areaName": &graphql.Field{
				Type: graphql.String,
			},
			"on": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)

type menu struct {
	Index    int                                 `yaml:"index"`
	AreaName string                              `yaml:"areaName"`
	Active   bool                                `yaml:"active"`
	Lutron   []Lutron                            `yaml:"lutron"`
	Kasa     []Kasa                              `yaml:"kasa"`
	Harmony  []Harmony                           `yaml:"harmony"`
	Hass     []homeassistant.HomeAssistantDevice `yaml:"hass"`
	Custom   []Custom                            `yaml:"custom"`
}

type Custom struct {
	Type    string   `yaml:"type"`
	State   string   `yaml:"state"`
	Name    string   `yaml:"name"`
	Devices []string `yaml:"devices"`
}

var menuType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Menu",
		Fields: graphql.Fields{
			"index": &graphql.Field{
				Type: graphql.Int,
			},
			"areaName": &graphql.Field{
				Type: graphql.String,
			},
			"active": &graphql.Field{
				Type: graphql.Boolean,
			},
			"lutron": &graphql.Field{
				Type: graphql.NewList(lutronType),
			},
			"kasa": &graphql.Field{
				Type: graphql.NewList(kasaType),
			},
			"harmony": &graphql.Field{
				Type: graphql.NewList(harmonyType),
			},
			"hass": &graphql.Field{
				Type: graphql.NewList(hassType),
			},
			"custom": &graphql.Field{
				Type: graphql.NewList(customType),
			},
		},
	},
)

var customType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Custom",
		Fields: graphql.Fields{
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"devices": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"state": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// Life360Member : The data for a life360 member
type Life360Member struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Name      string `json:"name"`
	Address1  string `json:"address1"`
	Battery   string `json:"battery"`
	IsDriving string `json:"isDriving"`
}

var life360Type = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Life360",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"firstName": &graphql.Field{
				Type: graphql.String,
			},
			"lastName": &graphql.Field{
				Type: graphql.String,
			},
			"latitude": &graphql.Field{
				Type: graphql.String,
			},
			"longitude": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"address1": &graphql.Field{
				Type: graphql.String,
			},
			"battery": &graphql.Field{
				Type: graphql.String,
			},
			"isDriving": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var homeAssistantType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "HomeAssistantDevice",
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
		},
	},
)
