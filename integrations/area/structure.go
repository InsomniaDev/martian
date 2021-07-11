package area

import "github.com/graphql-go/graphql"

type Area struct {
	Index    int          `yaml:"index"`
	AreaName string       `yaml:"areaName"`
	Active   bool         `yaml:"active"`
	Devices  []AreaDevice `yaml:"devices"`
}

var AreaType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Area",
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
			"devices": &graphql.Field{
				Type: graphql.NewList(AreaDeviceType),
			},
		},
	},
)

type AreaDevice struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Integration string `json:"integration"`
	Name        string `json:"name"`
	State       string `json:"state"`
	AreaName    string `json:"areaName"`
	Value       string `json:"value"`
}

var AreaDeviceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Device",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"integration": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"state": &graphql.Field{
				Type: graphql.String,
			},
			"areaName": &graphql.Field{
				Type: graphql.String,
			},
			"value": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)
