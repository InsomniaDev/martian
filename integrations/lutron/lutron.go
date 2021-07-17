package lutron

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/insomniadev/martian/database"
)

// Init starts up the lutron instance
func Init(configuration string) Lutron {
	var lutron Lutron
	err := json.Unmarshal([]byte(configuration), &lutron)
	if err != nil {
		// TODO: Change this away from being a panic
		panic(err)
	}

	if len(lutron.Inventory) == 0 {
		fileContents := loadIntegrationFile(lutron.Config)
		for _, device := range fileContents.LIPIDList.Zones {
			lutron.Inventory = append(lutron.Inventory, &LDevice{
				Name: device.Name,
				ID: device.ID,
				AreaName: device.Area.Name,
				Type: device.Type,
			})
		}
	}

	lutron.done = make(chan bool)

	lutron.Connect()
	return lutron
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
