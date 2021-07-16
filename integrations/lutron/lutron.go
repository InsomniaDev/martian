package lutron

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/insomniadev/martian/database"
)

// Init starts up the lutron instance
func Init(configuration string) Lutron {
	var lutron database.LutronConfig
	err := json.Unmarshal([]byte(configuration), &lutron)
	if err != nil {
		// TODO: Change this away from being a panic
		panic(err)
	}
	inventory := RetrieveLutronNodes()
	if len(inventory) == 0 {
		// Load the lutron configuration file
		fileContents := loadIntegrationFile(lutron)
		// insert into the database
		for _, name := range fileContents.LIPIDList.Zones {
			InsertLutronGraph(name.Name, name.ID, name.Area.Name, name.Type)
		}
		inventory = RetrieveLutronNodes()
	}
	l := &Lutron{
		Config: lutron,
		// Responses: make(chan string, 5),
		done:      make(chan bool),
		Inventory: inventory,
	}

	l.Connect()
	return *l
}

func loadIntegrationFile(config database.LutronConfig) CasetaIntegrationFile {
	jsonfile, err := os.Open("./config/" + config.File)

	fileContents := CasetaIntegrationFile{}
	jsonParser := json.NewDecoder(jsonfile)
	if err = jsonParser.Decode(&fileContents); err != nil {
		fmt.Println("parsing lutron config file", err.Error())
	}
	return fileContents
}
