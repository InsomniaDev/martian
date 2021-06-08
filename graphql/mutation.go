package graphql

import (
	"github.com/graphql-go/graphql"
)

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"turnDeviceOn": &graphql.Field{
			Type: lutronType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Description: "Turn on a Lutron Device",
			Resolve:     lutronTurnOnResolver,
		},
		"turnDeviceOff": &graphql.Field{
			Type: lutronType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Description: "Turn off a Lutron Device",
			Resolve:     lutronTurnOffResolver,
		},
		"setLutronDeviceToLevel": &graphql.Field{
			Type: lutronType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"level": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
			},
			Description: "Turn off a Lutron Device",
			Resolve:     lutronChangeDeviceToLevel,
		},
		"turnAllLightsOn": &graphql.Field{
			Type: graphql.Boolean,
			Description: "Turn all Lutron lights on",
			Resolve:     lutronTurnAllLightsOn,
		},
		"turnAllLightsOff": &graphql.Field{
			Type: graphql.Boolean,
			Description: "Turn all Lutron lights off",
			Resolve:     lutronTurnAllLightsOff,
		},
		"startHarmonyActivity": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Description: "Start a Harmony Activity",
			Resolve: harmonyStartActivityResolver,
		},
		"updateKasaArea": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"ipAddress": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"areaName": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Description: "Updates the area for a kasa device",
			Resolve: updateAreaForKasaDevice,
		},
		"turnKasaOn": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"ipAddress": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Description: "Turn on a Kasa plug",
			Resolve:     kasaTurnOnResolver,
		},
		"turnKasaOff": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"ipAddress": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Description: "Turn off a Kasa plug",
			Resolve:     kasaTurnOffResolver,
		},
	},
})
