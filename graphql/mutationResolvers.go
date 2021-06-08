package graphql

import (
	"strconv"
	"strings"

	"github.com/graphql-go/graphql"
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
	for i, dev := range Integrations.KasaData.Plugs {
		if dev.IPAddress == argID {
			Integrations.KasaData.Plugs[i].UpdateArea(areaName)
		}
	}
	return true, nil
}

// kasaTurnOffResolver turns off a kasa device by setting value to zero
func kasaTurnOffResolver(params graphql.ResolveParams) (interface{}, error) {
	argString := params.Args["ipAddress"].(string)
	for i, d := range Integrations.KasaData.Plugs {
		if d.IPAddress == argString {
			Integrations.KasaData.Plugs[i].PowerOff()
		}
	}
	return true, nil
}

// kasaTurnOnResolver turns on a kasa device by setting value to one hundred
func kasaTurnOnResolver(params graphql.ResolveParams) (interface{}, error) {
	argString := params.Args["ipAddress"].(string)
	for i, d := range Integrations.KasaData.Plugs {
		if d.IPAddress == argString {
			Integrations.KasaData.Plugs[i].PowerOn()
		}
	}
	return true, nil
}
