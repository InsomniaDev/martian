package graphql

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/graphql-go/graphql"
	"github.com/insomniadev/martian/integrations/area"
	"github.com/insomniadev/martian/integrations/homeassistant"
)

// Return a single lutron element
func lutronOneResolver(params graphql.ResolveParams) (interface{}, error) {
	argID, _ := params.Args["id"].(int)
	for _, d := range Integrations.LutronData.Inventory {
		if d.ID == argID {
			return d, nil
		}
	}
	return nil, nil
}

// Return all of the lutron elements
func lutronAllResolver(params graphql.ResolveParams) (interface{}, error) {
	return Integrations.LutronData.Inventory, nil
}

func getHarmonyActivities(params graphql.ResolveParams) (interface{}, error) {
	return Integrations.HarmonyData.Activities, nil
}

func getCurrentHarmonyActivity(params graphql.ResolveParams) (interface{}, error) {
	for _, data := range Integrations.HarmonyData.Activities {
		if data.ActivityID == Integrations.HarmonyData.CurrentActivity {
			return data, nil
		}
	}
	return nil, nil
}

func life360Members(params graphql.ResolveParams) (interface{}, error) {
	var members []Life360Member
	for _, data := range Integrations.Life3.Members {
		member := Life360Member{
			ID:        data.ID,
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Latitude:  data.Location.Latitude,
			Longitude: data.Location.Longitude,
			Name:      data.Location.Name,
			Address1:  data.Location.Address1,
			Battery:   data.Location.Battery,
			IsDriving: data.Location.IsDriving,
		}
		members = append(members, member)
	}
	return members, nil
}

func homeAssistantDevices(params graphql.ResolveParams) (interface{}, error) {
	var devices []homeassistant.HomeAssistantDevice
	hassType, _ := params.Args["type"].(string)
	hassName, _ := params.Args["name"].(string)
	for _, device := range Integrations.Hass.Devices {
		if hassType != "" && hassName != "" {
			if strings.EqualFold(device.Type, hassType) && strings.Contains(strings.ToLower(device.Name), strings.ToLower(hassName)) {
				devices = append(devices, device)
			}
		} else if hassType != "" {
			if strings.EqualFold(device.Type, hassType) {
				devices = append(devices, device)
			}
		} else if hassName != "" {
			if strings.Contains(strings.ToLower(device.Name), strings.ToLower(hassName)) {
				devices = append(devices, device)
			}
		}
	}
	if len(devices) > 0 {
		return devices, nil
	} else {
		return Integrations.Hass.Devices, nil
	}
}

func menuConfiguration(params graphql.ResolveParams) (interface{}, error) {
	Integrations.Menu = nil
	for _, k := range Integrations.Integrations {
		switch k {
		case "area":
			Integrations.Menu = area.CheckIndexForAreas(Integrations.Menu, Integrations.AreaIndexes)
		case "lutron":
			Integrations.Menu = area.LutronIntegration(Integrations.Menu, Integrations.LutronData.Inventory, Integrations.LutronData.InterfaceInventory)
		case "harmony":
			Integrations.Menu = area.HarmonyIntegration(Integrations.Menu, Integrations.HarmonyData)
		case "kasa":
			Integrations.Menu = area.KasaIntegration(Integrations.Menu, Integrations.KasaData, Integrations.KasaData.InterfaceDevices)
		case "life360":
			log.Info("Not implemented")
		case "hass":
			// TODO: Need to update this and implement it to display on the screen and update the devices accordingly
			log.Info("Not implemented")
		default:
			log.Info("This integration doesn't exist yet", k)
		}
	}

	// Cycle through the integrations and update indexes if there are any
	Integrations.Menu = area.CheckIndexForAreas(Integrations.Menu, Integrations.AreaIndexes)

	return Integrations.Menu, nil
}

func integrationResolver(params graphql.ResolveParams) (interface{}, error) {
	var integration IntegrationQueryType
	integration.Integrations = Integrations.Integrations
	for _, k := range Integrations.Integrations {
		switch k {
		case "lutron":
			integration.Lutron = Integrations.LutronData
		case "kasa":
			integration.Kasa = Integrations.KasaData
		case "hass":
			integration.Hass = Integrations.Hass
		case "harmony":
			integration.Harmony = Integrations.HarmonyData
		default:
			log.Info("This integration doesn't exist yet", k)
		}
	}

	return integration, nil
}

// getAreaNamesResolver will get a string of area names and return them
func getAreaNamesResolver(params graphql.ResolveParams) (interface{}, error) {
	var areaNames []string
	for _, area := range Integrations.Menu {
		areaNames = append(areaNames, area.AreaName)
	}

	return areaNames, nil
}
