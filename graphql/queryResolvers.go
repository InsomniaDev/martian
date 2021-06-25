package graphql

import (
	"encoding/json"
	"strconv"
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
	var menuItems []menu
	harmonyID := 1
	for _, menuValue := range Integrations.Menu {
		menuItem := menu{
			Index:    menuValue.Index,
			AreaName: menuValue.AreaName,
			Active:   false,
		}
		if menuValue.Lutron != nil {
			var lutronDevices []Lutron
			for _, id := range menuValue.Lutron {
				for _, lutronID := range Integrations.LutronData.Inventory {
					if id == lutronID.ID {
						newDev := Lutron{
							AreaName: lutronID.AreaName,
							DeviceID: lutronID.ID,
							Name:     lutronID.Name,
							Type:     lutronID.Type,
							Value:    lutronID.Value,
							State:    lutronID.State,
						}
						if lutronID.State == "on" {
							menuItem.Active = true
						}
						lutronDevices = append(lutronDevices, newDev)
					}
				}
			}
			menuItem.Lutron = lutronDevices
		}
		if menuValue.Kasa != nil {
			var kasaDevices []Kasa
			for _, id := range menuValue.Kasa {
				for _, kasaDev := range Integrations.KasaData.Plugs {
					if id == kasaDev.IPAddress {
						newDev := Kasa{
							ID:        kasaDev.ID,
							AreaName:  kasaDev.AreaName,
							IPAddress: kasaDev.IPAddress,
							IsOn:      kasaDev.PlugInfo.On,
							Name:      kasaDev.Name,
						}
						if newDev.IsOn {
							menuItem.Active = true
						}
						kasaDevices = append(kasaDevices, newDev)
					}
				}
			}
			menuItem.Kasa = kasaDevices
		}
		if menuValue.Harmony != nil {
			var currentActivity []Harmony
			for _, activity := range Integrations.HarmonyData.Activities {
				if activity.ActivityID == Integrations.HarmonyData.CurrentActivity {
					newActivity := Harmony{
						ID:         Integrations.HarmonyData.ActivityID,
						ActivityID: activity.ActivityID,
						Name:       activity.Name,
					}
					harmonyID = harmonyID + 1
					currentActivity = append(currentActivity, newActivity)
				}
			}
			menuItem.Harmony = currentActivity
		}
		if menuValue.Hass != nil {
			var devices []homeassistant.HomeAssistantDevice
			for _, configDevice := range menuValue.Hass {
				for _, device := range Integrations.Hass.Devices {
					if device.EntityId == configDevice {
						device.Name = strings.ToLower(device.Name)
						device.Name = strings.Replace(device.Name, strings.ToLower(menuItem.AreaName), "", -1)
						device.Name = strings.Replace(device.Name, strings.ToLower(device.Type) + "s", "", -1)
						device.Name = strings.Replace(device.Name, strings.ToLower(device.Type), "", -1)
						device.Name = strings.Title(device.Name)
						devices = append(devices, device)
						if strings.ToLower(device.State) == "on" {
							menuItem.Active = true
						}
					}
				}
			}
			menuItem.Hass = devices
		}
		if menuValue.Custom != nil {
			var customActivities []Custom
			for _, data := range menuValue.Custom {
				var newCustom Custom
				jsonData, _ := json.Marshal(data)
				json.Unmarshal(jsonData, &newCustom)
				if strings.ToUpper(newCustom.Type) == "LUTRON" {
					lightOn := false
					for i := range Integrations.LutronData.Inventory {
						for _, val := range newCustom.Devices {
							if strconv.Itoa(Integrations.LutronData.Inventory[i].ID) == val {
								if Integrations.LutronData.Inventory[i].Value > 0 {
									lightOn = true
								}
								break
							}
						}
						if lightOn {
							break
						}
					}
					if lightOn {
						newCustom.State = "on"
						menuItem.Active = true
					} else {
						newCustom.State = "off"
					}
				}
				customActivities = append(customActivities, newCustom)
			}
			menuItem.Custom = customActivities
		}
		menuItems = append(menuItems, menuItem)
	}

	return menuItems, nil
}