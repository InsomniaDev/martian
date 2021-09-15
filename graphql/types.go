package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/insomniadev/martian/integrations/harmony"
	"github.com/insomniadev/martian/integrations/homeassistant"
	"github.com/insomniadev/martian/integrations/kasa"
	"github.com/insomniadev/martian/integrations/lutron"
)

type IntegrationQueryType struct {
	Lutron       lutron.Lutron               `json:"lutron"`
	Hass         homeassistant.HomeAssistant `json:"hass"`
	Harmony      harmony.Device              `json:"harmony"`
	Kasa         kasa.Devices                `json:"kasa"`
	Integrations []string                    `json:"integrations"`
}

var integrationsType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Integrations",
		Fields: graphql.Fields{
			"integrations": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"lutron": &graphql.Field{
				Type: lutron.GraphqlLutronType,
			},
			"hass": &graphql.Field{
				Type: homeassistant.GraphqlHomeAssistantType,
			},
			"harmony": &graphql.Field{
				Type: harmony.GraphqlType,
			},
			"kasa": &graphql.Field{
				Type: kasa.GraphqlKasaType,
			},
		},
	},
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

type Custom struct {
	Type    string   `yaml:"type"`
	State   string   `yaml:"state"`
	Name    string   `yaml:"name"`
	Devices []string `yaml:"devices"`
}

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
			"areaName": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
