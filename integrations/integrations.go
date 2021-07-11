package integrations

import (
	"encoding/json"
	"fmt"

	"github.com/insomniadev/martian/integrations/config"
	"github.com/insomniadev/martian/integrations/harmony"
	"github.com/insomniadev/martian/integrations/homeassistant"
	"github.com/insomniadev/martian/integrations/kasa"
	"github.com/insomniadev/martian/integrations/life360"
	"github.com/insomniadev/martian/integrations/lutron"
	"github.com/insomniadev/martian/modules/database"
)

type Integrations struct {
	Menu        []config.Menu
	LutronData  lutron.Lutron
	HarmonyData harmony.Device
	Hass        homeassistant.HomeAssistant
	KasaData    kasa.Devices
	Life3       life360.Life360
	Database    database.Database
	// Zwave       zwave.Zwave
}

func (i *Integrations) Init() {

	// Get all the created integrations
	storedIntegrations, err := i.Database.RetrieveAllValuesInBucket(database.IntegrationBucket)
	if err != nil {
		// TODO: Change this away from being a panic
		panic(err)
	}

	for k := range storedIntegrations {
		switch k {
		case "lutron":
			i.LutronData = lutron.Init(storedIntegrations[k])
			jsonValue, _ := json.Marshal(i.LutronData)
			fmt.Println(string(jsonValue))
		case "harmony":
			fmt.Println("Not implemented")
		case "kasa":
			fmt.Println("Not implemented")
		case "life360":
			fmt.Println("Not implemented")
		case "hass":
			fmt.Println("Not implemented")
		default:
			fmt.Println("This integration doesn't exist yet", k)
		}
	}
	// Cycle through the integrations

	// TODO: This needs to load up each based on if it is available, there is no point in loading up all of them
	i.HarmonyData.Init()
	// i.KasaData.Init()
	i.Menu = config.LoadMenu()
	// i.Life3.Authenticate()
	// go i.Life3.SyncMemberStatus()

	go i.Hass.Init()
	// i.Zwave.ConnectToTopic()
}
