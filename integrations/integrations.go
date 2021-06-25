package integrations

import (
	"github.com/insomniadev/martian/integrations/config"
	"github.com/insomniadev/martian/integrations/harmony"
	"github.com/insomniadev/martian/integrations/homeassistant"
	"github.com/insomniadev/martian/integrations/kasa"
	"github.com/insomniadev/martian/integrations/life360"
	"github.com/insomniadev/martian/integrations/lutron"
)

type Integrations struct {
	Menu        []config.Menu
	LutronData  lutron.Lutron
	HarmonyData harmony.Device
	Hass        homeassistant.HomeAssistant
	KasaData    kasa.Devices
	Life3       life360.Life360
	// Zwave       zwave.Zwave
}

func (i *Integrations) Init() {
	// TODO: This needs to load up each based on if it is available, there is no point in loading up all of them
	// i.LutronData = lutron.Init()
	i.HarmonyData.Init()
	// i.KasaData.Init()
	i.Menu = config.LoadMenu()
	// i.Life3.Authenticate()
	// go i.Life3.SyncMemberStatus()

	go i.Hass.Init()
	// i.Zwave.ConnectToTopic()
}
