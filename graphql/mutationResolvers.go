package graphql

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/insomniadev/martian/integrations/area"
	"github.com/insomniadev/martian/integrations/harmony"
	"github.com/insomniadev/martian/integrations/homeassistant"
	"github.com/insomniadev/martian/integrations/kasa"
	"github.com/insomniadev/martian/integrations/lutron"
)

// lutronTurnOffResolver turns off a lutron device by setting value to zero
func lutronTurnOffResolver(params graphql.ResolveParams) (interface{}, error) {
	argString := params.Args["id"].(string)
	argID, _ := strconv.Atoi(argString)
	for _, d := range Integrations.LutronData.Inventory {
		if d.ID == argID {
			Integrations.LutronData.SetById(argID, 0)
			return d, nil
		}
	}
	return false, nil
}

// lutronTurnOnResolver turns on a lutron device by setting value to one hundred
func lutronTurnOnResolver(params graphql.ResolveParams) (interface{}, error) {
	argString := params.Args["id"].(string)
	argID, _ := strconv.Atoi(argString)
	for _, d := range Integrations.LutronData.Inventory {
		if d.ID == argID {
			Integrations.LutronData.SetById(argID, 100)
			return d, nil
		}
	}
	return false, nil
}

// lutronTurnOnResolver turns on a lutron device by setting value to one hundred
func lutronChangeDeviceToLevel(params graphql.ResolveParams) (interface{}, error) {
	argString := params.Args["id"].(string)
	argID, _ := strconv.Atoi(argString)
	argLevel := params.Args["level"].(float64)
	for _, d := range Integrations.LutronData.Inventory {
		if d.ID == argID {
			Integrations.LutronData.SetById(argID, argLevel)
			return d, nil
		}
	}
	return false, nil
}

// lutronTurnOnResolver turns on a lutron device by setting value to one hundred
func lutronTurnAllLightsOn(params graphql.ResolveParams) (interface{}, error) {
	for _, d := range Integrations.LutronData.Inventory {
		if strings.ToUpper(d.Type) == "LIGHT" && d.Value == 0 {
			Integrations.LutronData.SetById(d.ID, 100)
		}
	}
	return true, nil
}

// lutronTurnOnResolver turns on a lutron device by setting value to one hundred
func lutronTurnAllLightsOff(params graphql.ResolveParams) (interface{}, error) {
	for _, d := range Integrations.LutronData.Inventory {
		if strings.ToUpper(d.Type) == "LIGHT" && d.Value > 0 {
			Integrations.LutronData.SetById(d.ID, 0)
		}
	}
	return true, nil
}

// harmonyStartActivityResolver will change the activity to the one specified
func harmonyStartActivityResolver(params graphql.ResolveParams) (interface{}, error) {
	argID, _ := params.Args["id"].(string)
	Integrations.HarmonyData.StartActivity(argID)
	return true, nil
}

// updateAreaForKasaDevice will update the area for the kasa device to match
func updateAreaForKasaDevice(params graphql.ResolveParams) (interface{}, error) {
	argID, _ := params.Args["ipAddress"].(string)
	areaName, _ := params.Args["areaName"].(string)
	for i, dev := range Integrations.KasaData.Devices {
		if dev.IPAddress == argID {
			Integrations.KasaData.Devices[i].UpdateArea(areaName)
		}
	}
	return true, nil
}

// kasaTurnOffResolver turns off a kasa device by setting value to zero
func kasaTurnOffResolver(params graphql.ResolveParams) (interface{}, error) {
	argString := params.Args["ipAddress"].(string)
	for i, d := range Integrations.KasaData.Devices {
		if d.IPAddress == argString {
			Integrations.KasaData.Devices[i].PowerOff()
		}
	}
	return true, nil
}

// kasaTurnOnResolver turns on a kasa device by setting value to one hundred
func kasaTurnOnResolver(params graphql.ResolveParams) (interface{}, error) {
	argString := params.Args["ipAddress"].(string)
	for i, d := range Integrations.KasaData.Devices {
		if d.IPAddress == argString {
			Integrations.KasaData.Devices[i].PowerOn()
		}
	}
	return true, nil
}

// changeHassDeviceStatusResolver changes the status of the Hass device
func changeHassDeviceStatusResolver(params graphql.ResolveParams) (interface{}, error) {
	entityId := params.Args["entityId"].(string)
	activated := params.Args["activated"].(bool)
	fmt.Println("got something here", entityId, activated)
	for _, d := range Integrations.Hass.Devices {
		if d.EntityId == entityId {
			Integrations.Hass.CallService(d, activated)
		}
	}
	return true, nil
}

// changeDeviceStatus is the one mutation to rule them all
func changeDeviceStatus(params graphql.ResolveParams) (interface{}, error) {
	id := params.Args["id"].(string)
	status := params.Args["status"].(string)
	level := params.Args["level"].(string)

	switch params.Args["integration"].(string) {
	case "lutron":
		lutronId, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		switch status {
		case "off":
			changeLutronDevice(lutronId, 0)
		case "on":
			changeLutronDevice(lutronId, 100)
		case "dim":
			level, err := strconv.ParseFloat(level, 64)
			if err != nil {
				return nil, err
			}
			changeLutronDevice(lutronId, level)
		}
	case "kasa":
		changeKasaDevice(id, status)
	case "hass":
		hassDevice(id, status)
	case "harmony":
		_, err := harmonyStartActivityResolver(params)
		if err != nil {
			return false, nil
		}
	}
	return true, nil
}

func changeLutronDevice(id int, level float64) {
	for _, d := range Integrations.LutronData.Inventory {
		if d.ID == id {
			Integrations.LutronData.SetById(id, level)
		}
	}
}

func changeKasaDevice(id string, status string) {
	for i, d := range Integrations.KasaData.Devices {
		if d.IPAddress == id {
			if status == "on" {
				Integrations.KasaData.Devices[i].PowerOn()
			} else {
				Integrations.KasaData.Devices[i].PowerOff()
			}
		}
	}
}

func hassDevice(id string, status string) {
	activated := false
	if status == "on" {
		activated = true
	}
	for _, d := range Integrations.Hass.Devices {
		if d.EntityId == id {
			Integrations.Hass.CallService(d, activated)
		}
	}
}

// updateIntegration will create new integrations that will be stored in the local database
func updateIntegration(params graphql.ResolveParams) (interface{}, error) {
	integrationType := params.Args["type"].(string)
	integrationValue := params.Args["value"].(string)
	newIntegration := false
	switch integrationType {
	case "lutron":
		var lutron lutron.Lutron
		err := json.Unmarshal([]byte(integrationValue), &lutron)
		if err != nil {
			return false, err
		}
		Integrations.Database.PutIntegrationValue(integrationType, lutron)
		newIntegration = true
	case "harmony":
		// var harmony harmony.Device
		// err := json.Unmarshal([]byte(integrationValue), &harmony)
		// if err != nil {
		// 	return false, err
		// }
		Integrations.Database.PutIntegrationValue(integrationType, "")
		newIntegration = true
	case "kasa":
		currentDevices := len(Integrations.KasaData.Devices)
		var ipAddress kasa.Devices
		err := json.Unmarshal([]byte(integrationValue), &ipAddress)
		if err != nil {
			return false, err
		}
		if ipAddress.IpAddressCidr != Integrations.KasaData.IpAddressCidr {
			Integrations.KasaData.IpAddressCidr = ipAddress.IpAddressCidr
			Integrations.Database.PutIntegrationValue(integrationType, Integrations.KasaData)
		}
		Integrations.KasaData.Discover()
		if len(Integrations.KasaData.Devices) > currentDevices {
			Integrations.Database.PutIntegrationValue(integrationType, Integrations.KasaData)
		} else {
			// TODO: Need to fix this piece, it is constantly assigning it as "", even though there are new ones
			Integrations.Database.PutIntegrationValue(integrationType, "")
		}
		newIntegration = true
	case "life360":
		fmt.Println("Not implemented")
	case "hass":
		// Will pass in {"url":"","token":""}
		var hass homeassistant.HomeAssistant
		err := json.Unmarshal([]byte(integrationValue), &hass)
		if err != nil {
			return false, err
		}
		Integrations.Database.PutIntegrationValue(integrationType, hass)
		newIntegration = true
	default:
		fmt.Println("This integration doesn't exist yet", integrationType)
	}

	if newIntegration {
		Integrations.Init()
	}
	return true, nil
}

func changeKasaDeviceArea(params graphql.ResolveParams) (interface{}, error) {
	ipAddress := params.Args["ipAddress"].(string)
	areaName := params.Args["area"].(string)

	err := Integrations.KasaData.ChangeAreaForKasaDevice(ipAddress, areaName)
	Integrations.Init()
	return true, err
}

// updateIndexForArea will update the index for the provided areaname
func updateIndexForArea(params graphql.ResolveParams) (interface{}, error) {
	areaName := params.Args["areaName"].(string)
	index := params.Args["index"].(int)

	udpatedAreaIndex := area.Area{AreaName: areaName, Index: index}
	menuValues, err := area.InsertAreaIndex(Integrations.Menu, udpatedAreaIndex)
	if err != nil {
		return false, err
	}
	Integrations.Menu = menuValues
	Integrations.Init()
	return true, err
}

type dataDevices struct {
	Data string
}

// selectDevicesForIntegration will add in the selected devices for the provided integration type
func selectDevicesForIntegration(params graphql.ResolveParams) (interface{}, error) {
	integration := params.Args["integration"].(string)
	devices := params.Args["devices"].([]interface{})
	addDevices := params.Args["addDevices"].(bool)
	automationDevice := params.Args["automationDevice"].(bool)
	var daters []string
	for a := range devices {
		daters = append(daters, devices[a].(string))
	}

	switch integration {
	case "hass":
		err := Integrations.Hass.UpdateSelectedDevices(daters, addDevices, automationDevice)
		if err != nil {
			return false, err
		}
	case "kasa":
		err := Integrations.KasaData.UpdateSelectedDevices(daters, addDevices, automationDevice)
		if err != nil {
			return false, err
		}
	case "lutron":
		var intDaters []int
		for a := range daters {
			number, err := strconv.Atoi(daters[a])
			if err != nil {
				return false, err
			}
			intDaters = append(intDaters, number)
		}
		err := Integrations.LutronData.UpdateSelectedDevices(intDaters, addDevices, automationDevice)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// editDeviceConfiguration will edit the passed in configuration device
func editDeviceConfiguration(params graphql.ResolveParams) (interface{}, error) {
	integration := params.Args["integration"].(string)
	device := params.Args["device"].(string)
	removeEdit := params.Args["removeEdit"].(bool)

	switch integration {
	case "hass":
		var hass homeassistant.HomeAssistantDevice
		err := json.Unmarshal([]byte(device), &hass)
		Integrations.Hass.EditDeviceConfiguration(hass, removeEdit)
		if err != nil {
			return false, err
		}
	case "kasa":
		var kasaDevice kasa.KasaDevice
		err := json.Unmarshal([]byte(device), &kasaDevice)
		Integrations.KasaData.EditDeviceConfiguration(kasaDevice, removeEdit)
		if err != nil {
			return false, err
		}
	case "lutron":
		var lutronDevice lutron.LDevice
		err := json.Unmarshal([]byte(device), &lutronDevice)
		Integrations.LutronData.EditDeviceConfiguration(lutronDevice, removeEdit)
		if err != nil {
			return false, err
		}
	case "harmony":
		var harmony harmony.Device
		err := json.Unmarshal([]byte(device), &harmony)
		Integrations.HarmonyData.EditDeviceConfiguration(harmony, removeEdit)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
