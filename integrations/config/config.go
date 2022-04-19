package config

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

// LoadFile loads up the yaml file
func (c *Config) LoadFile() {
	yamlFile, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

// LoadLutron loads the lutron configuration from the yaml file
func LoadLutron() LutronConfig {
	c := Config{}
	c.LoadFile()

	return c.LutronData
}

// LoadMyQ loads the myq configuration from the yaml file
func LoadMyQ() (string, string) {
	c := Config{}
	c.LoadFile()

	return c.MyQData.Username, c.MyQData.Password
}

// LoadHarmony loads the Harmony configuration from the yaml file
func LoadHarmony() string {
	c := Config{}
	c.LoadFile()

	return c.HarmonyData.IPAddress
}

// LoadKasa loads the Kasa configuration from the yaml file
func LoadKasa() []string {
	c := Config{}
	c.LoadFile()

	return c.KasaData.Devices
}

// LoadSleepIq loads the sleep number configuration from the yaml file
func LoadSleepIq() (string, string) {
	c := Config{}
	c.LoadFile()

	return c.SleepIqData.Username, c.SleepIqData.Password
}

// LoadLife360 loads the life360 configuration from the yaml file
func LoadLife360() (string, string, string) {
	c := Config{}
	c.LoadFile()

	return c.Life360Data.Username, c.Life360Data.Password, c.Life360Data.AuthenticationToken
}

// LoadZwave loads the zwave configuration from the yaml file
func LoadZwave() string {
	c := Config{}
	c.LoadFile()

	return c.ZwaveData.URL
}

// LoadMenu loads up the set menu configuration
func LoadMenu() []Menu {
	c := Config{}
	c.LoadFile()
	return c.MenuConfig
}

// LoadHomeAssistant loads up the set home assistant configuration
func LoadHomeAssistant() Hass {
	c := Config{}
	c.LoadFile()
	return c.HomeAssistant
}

// LoadHubitat loads up the set hubitat configuration
func LoadHubitat() (accessKey string, URL string) {
	c := Config{}
	c.LoadFile()
	return c.Hubitat.AccessKey, c.Hubitat.URL
}
