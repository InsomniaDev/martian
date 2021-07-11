package graphql

import (
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
			for _, id := range menuValue.Lutron {
				for _, lutronID := range Integrations.LutronData.Inventory {
					if id == lutronID.ID {
						newDev := menuDevice{
							AreaName:    lutronID.AreaName,
							Id:          strconv.Itoa(lutronID.ID),
							Name:        lutronID.Name,
							Type:        lutronID.Type,
							Value:       strconv.FormatFloat(lutronID.Value, 'E', -1, 64),
							State:       lutronID.State,
							Integration: "lutron",
						}
						if strings.ToLower(lutronID.State) != "off" {
							menuItem.Active = true
						}
						menuItem.Devices = append(menuItem.Devices, newDev)
					}
				}
			}
		}
		if menuValue.Kasa != nil {
			for _, id := range menuValue.Kasa {
				for _, kasaDev := range Integrations.KasaData.Plugs {
					if id == kasaDev.IPAddress {
						state := "off"
						if kasaDev.PlugInfo.On {
							state = "on"
							menuItem.Active = true
						}
						newDev := menuDevice{
							AreaName:    kasaDev.AreaName,
							Id:          kasaDev.IPAddress,
							Name:        kasaDev.Name,
							Type:        "UKN", //TODO: Need to work on discovering what type of device the kasa device is for the UI
							State:       state,
							Integration: "kasa",
						}
						menuItem.Devices = append(menuItem.Devices, newDev)
					}
				}
			}
		}
		if menuValue.Harmony != nil {
			for _, activity := range Integrations.HarmonyData.Activities {
				if activity.ActivityID == Integrations.HarmonyData.CurrentActivity {
					newDev := menuDevice{
						Id:          Integrations.HarmonyData.ActivityID,
						Name:        activity.Name,
						Integration: "harmony",
					}
					harmonyID = harmonyID + 1
					menuItem.Devices = append(menuItem.Devices, newDev)
				}
			}
		}
		if menuValue.Hass != nil {
			for _, configDevice := range menuValue.Hass {
				for _, device := range Integrations.Hass.Devices {
					if device.EntityId == configDevice {
						if strings.ToLower(device.State) == "on" {
							menuItem.Active = true
						}
						device.Name = strings.ToLower(device.Name)
						device.Name = strings.Replace(device.Name, strings.ToLower(menuItem.AreaName), "", -1)
						device.Name = strings.Replace(device.Name, strings.ToLower(device.Type)+"s", "", -1)
						device.Name = strings.Replace(device.Name, strings.ToLower(device.Type), "", -1)
						device.Name = strings.Title(device.Name)

						newDev := menuDevice{
							AreaName:    device.Group,
							Id:          device.EntityId,
							Name:        device.Name,
							Type:        device.Type,
							State:       device.State,
							Integration: "hass",
						}
						menuItem.Devices = append(menuItem.Devices, newDev)
					}
				}
			}
		}
		// if menuValue.Custom != nil {
		// 	var customActivities []Custom
		// 	for _, data := range menuValue.Custom {
		// 		var newCustom Custom
		// 		jsonData, _ := json.Marshal(data)
		// 		json.Unmarshal(jsonData, &newCustom)
		// 		if strings.ToUpper(newCustom.Type) == "LUTRON" {
		// 			lightOn := false
		// 			for i := range Integrations.LutronData.Inventory {
		// 				for _, val := range newCustom.Devices {
		// 					if strconv.Itoa(Integrations.LutronData.Inventory[i].ID) == val {
		// 						if Integrations.LutronData.Inventory[i].Value > 0 {
		// 							lightOn = true
		// 						}
		// 						break
		// 					}
		// 				}
		// 				if lightOn {
		// 					break
		// 				}
		// 			}
		// 			if lightOn {
		// 				newCustom.State = "on"
		// 				menuItem.Active = true
		// 			} else {
		// 				newCustom.State = "off"
		// 			}
		// 		}
		// 		customActivities = append(customActivities, newCustom)
		// 	}
		// 	menuItem.Custom = customActivities
		// }
		menuItems = append(menuItems, menuItem)
	}

	return menuItems, nil
}
