package integrations

import (
	"fmt"

	"github.com/insomniadev/martian/database"
	"github.com/insomniadev/martian/integrations/area"
	"github.com/insomniadev/martian/integrations/harmony"
	"github.com/insomniadev/martian/integrations/homeassistant"
	"github.com/insomniadev/martian/integrations/kasa"
	"github.com/insomniadev/martian/integrations/life360"
	"github.com/insomniadev/martian/integrations/lutron"
)

type Integrations struct {
	Integrations []string
	Menu         []area.Area
	LutronData   lutron.Lutron
	HarmonyData  harmony.Device
	Hass         homeassistant.HomeAssistant
	KasaData     kasa.Devices
	Life3        life360.Life360
	Database     database.Database
	// Zwave       zwave.Zwave
}

func (i *Integrations) Init() {

	i.Integrations = []string{}
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
			i.Menu = area.LutronIntegration(i.Menu, i.LutronData.Inventory)
			i.Integrations = append(i.Integrations, "lutron")
		case "harmony":
			fmt.Println("Not implemented")
			i.Integrations = append(i.Integrations, "harmony")
		case "kasa":
			i.KasaData.Init(storedIntegrations[k])
			if len(i.KasaData.Plugs) == 0 {
				i.KasaData.Discover()
			}
			i.Menu = area.KasaIntegration(i.Menu, i.KasaData)
			i.Integrations = append(i.Integrations, "kasa")
		case "life360":
			i.Life3.Authenticate()
			go i.Life3.SyncMemberStatus()
			i.Integrations = append(i.Integrations, "life360")
		case "hass":
			go i.Hass.Init()
			i.Integrations = append(i.Integrations, "hass")
		default:
			fmt.Println("This integration doesn't exist yet", k)
		}
	}
	// Cycle through the integrations

	// TODO: This needs to load up each based on if it is available, there is no point in loading up all of them
	i.HarmonyData.Init()

	// i.Zwave.ConnectToTopic()
}

func AddAreas() {

}
