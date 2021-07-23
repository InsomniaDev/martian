package database

import "github.com/graphql-go/graphql"

// LutronConfig is the configuration for the Lutron Hub
type LutronConfig struct {
	URL      string `json:"url"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	File     string `json:"file"`
}

// GrapqhlLutronConfigType is the graphql object for the lutron devices
var GrapqhlLutronConfigType = graphql.NewObject(graphql.ObjectConfig{
	Name: "LutronConfigType",
	Fields: graphql.Fields{
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"port": &graphql.Field{
			Type: graphql.Int,
		},
		"username": &graphql.Field{
			Type: graphql.String,
		},
		"password": &graphql.Field{
			Type: graphql.String,
		},
		"file": &graphql.Field{
			Type: graphql.Float,
		},
	},
})
