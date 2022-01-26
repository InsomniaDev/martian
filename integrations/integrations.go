package integrations

import (
	"log"

	"github.com/insomniadev/martian/database"
	"github.com/insomniadev/martian/integrations/area"
	"github.com/insomniadev/martian/integrations/harmony"
	"github.com/insomniadev/martian/integrations/homeassistant"
	"github.com/insomniadev/martian/integrations/kasa"
	"github.com/insomniadev/martian/integrations/life360"
	"github.com/insomniadev/martian/integrations/lutron"
	"github.com/insomniadev/martian/logger"
	"github.com/sirupsen/logrus"
)

type Integrations struct {
	Integrations []string
	Menu         []area.Area
	AreaIndexes  []area.Area
	LutronData   lutron.Lutron
	HarmonyData  harmony.Device
	Hass         homeassistant.HomeAssistant
	KasaData     kasa.Devices
	Life3        life360.Life360
	Database     database.Database
	// Zwave       zwave.Zwave
}

func (i *Integrations) Init() {

	i.Database = database.MartianData
	i.Integrations = []string{}
	// Get all the created integrations
	storedIntegrations, err := i.Database.RetrieveAllValuesInBucket(database.IntegrationBucket)
	if err != nil {
		// TODO: Change this away from being a panic
		log.Println(err)
		// panic(err)
	}

	for k := range storedIntegrations {
		switch k {
		case "area":
			i.AreaIndexes = area.Init(storedIntegrations[k])
		case "lutron":
			i.LutronData, err = lutron.Init(storedIntegrations[k])
			if err != nil {
				logger.Logger().Log(logrus.ErrorLevel, err)
			}
			i.Menu = area.LutronIntegration(i.Menu, i.LutronData.Inventory, i.LutronData.InterfaceInventory)
			i.Integrations = append(i.Integrations, "lutron")
		case "harmony":
			i.HarmonyData.Init(storedIntegrations[k])
			i.Menu = area.HarmonyIntegration(i.Menu, i.HarmonyData)
			i.Integrations = append(i.Integrations, "harmony")
		case "kasa":
			i.KasaData.Init(storedIntegrations[k])
			i.Menu = area.KasaIntegration(i.Menu, i.KasaData, i.KasaData.InterfaceDevices)
			i.Integrations = append(i.Integrations, "kasa")
		case "life360":
			i.Life3.Authenticate()
			go i.Life3.SyncMemberStatus()
			i.Integrations = append(i.Integrations, "life360")
		case "hass":
			go i.Hass.Init(storedIntegrations[k])
			i.Integrations = append(i.Integrations, "hass")
		default:
			log.Println("This integration doesn't exist yet", k)
		}
	}
	// Cycle through the integrations
	if len(i.AreaIndexes) > 0 {
		log.Println("Devices have been found")
		i.Menu = area.CheckIndexForAreas(i.Menu, i.AreaIndexes)
	}

	// TODO: This needs to load up each based on if it is available, there is no point in loading up all of them
	// i.HarmonyData.Init()

	// i.Zwave.ConnectToTopic()
}

func AddAreas() {

}
