package graphql

import (
	"github.com/graphql-go/graphql"
)

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
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
		"changeDeviceStatus": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"status": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"level": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"integration": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Description: "Change the device status independent of it's integration type",
			Resolve:     changeDeviceStatus,
		},
		"changeAreaForKasaDevice": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"ipAddress": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"area": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Description: "Create or update an integration with the Martian API",
			Resolve:     changeKasaDeviceArea,
		},
		"updateIntegration": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"type": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"value": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Description: "Create or update an integration with the Martian API",
			Resolve:     updateIntegration,
		},
	},
})
