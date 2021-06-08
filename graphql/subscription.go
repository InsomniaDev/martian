package graphql

import (
	"github.com/graphql-go/graphql"
)

var rootSubscription = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootSubscription",
	Fields: graphql.Fields{
		"lutronChanges": &graphql.Field{
			Type:        graphql.NewList(lutronType),
			Description: "Retrieve Lutron state changes",
			Resolve:     lutronAllResolver,
		},
		"harmonyChange": &graphql.Field{
			Type:        harmonyType,
			Description: "Retrieve harmony state changes",
			Resolve:     getCurrentHarmonyActivity,
		},
		"menuChange": &graphql.Field{
			Type:        graphql.NewList(menuType),
			Description: "Retrieve menu configuration changes",
			Resolve:     menuConfiguration,
		},
		"life360Change": &graphql.Field{
			Type:        graphql.NewList(life360Type),
			Description: "Retrieve menu configuration changes",
			Resolve:     life360Members,
		},
	},
})
