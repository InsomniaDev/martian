package graphql

import (
	"strings"

	"github.com/graphql-go/graphql"
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

func getKasaDevices(params graphql.ResolveParams) (interface{}, error) {
	var devices []Kasa
	for _, data := range Integrations.KasaData.Plugs {
		dev := Kasa{
			AreaName:  data.AreaName,
			IPAddress: data.IPAddress,
			IsOn:      data.PlugInfo.On,
			Name:      data.Name,
			Type:      data.Type,
		}
		devices = append(devices, dev)
	}
	return devices, nil
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
	
	return Integrations.Menu, nil
}
