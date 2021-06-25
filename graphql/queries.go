package graphql

import (
	"github.com/graphql-go/graphql"
)

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"lutronDevice": &graphql.Field{
			Type: lutronType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Description: "All of the Lutron devices",
			Resolve:     lutronOneResolver,
		},
		"lutronDevices": &graphql.Field{
			Type:        graphql.NewList(lutronType),
			Description: "All of the Lutron devices",
			Resolve:     lutronAllResolver,
		},
		"getHarmonyActivities": &graphql.Field{
			Type:        graphql.NewList(harmonyType),
			Description: "All of the Harmony Activities",
			Resolve:     getHarmonyActivities,
		},
		"currentHarmonyActivity": &graphql.Field{
			Type:        harmonyType,
			Description: "Current harmony activity",
			Resolve:     getCurrentHarmonyActivity,
		},
		"kasaDevices": &graphql.Field{
			Type:        graphql.NewList(kasaType),
			Description: "All of the kasa devices",
			Resolve:     getKasaDevices,
		},
		"menuConfiguration": &graphql.Field{
			Type:        graphql.NewList(menuType),
			Description: "The configuration that is returned for the UI to display",
			Resolve:     menuConfiguration,
		},
		"life360": &graphql.Field{
			Type:        graphql.NewList(life360Type),
			Description: "The life360 members and their current locations",
			Resolve:     life360Members,
		},
		"homeAssistantDevices": &graphql.Field{
			Type:        graphql.NewList(homeAssistantType),
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
					DefaultValue: "",
				},
				"type": &graphql.ArgumentConfig{
					Type: graphql.String,
					DefaultValue: "",
				},
			},
			Description: "The Home Assistant Device by Name or Type",
			Resolve:     homeAssistantDevices,
		},
	},
})