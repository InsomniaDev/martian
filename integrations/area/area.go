package area

import (
	"strconv"
	"strings"

	"github.com/insomniadev/martian/integrations/kasa"
	"github.com/insomniadev/martian/integrations/lutron"
)

func LutronIntegration(areas []Area, devices []*lutron.LDevice) []Area {
	for _, lutronDev := range devices {
		foundArea := false
		areaName := strings.TrimSpace(lutronDev.AreaName)
		if areaName == "" {
			areaName = "UKN"
		}
		for area := range areas {
			if strings.EqualFold(areas[area].AreaName, areaName) {
				areas[area].addLutron(lutronDev)
				foundArea = true
			}
		}
		if !foundArea {
			newArea := Area{
				AreaName: areaName,
			}
			newArea.addLutron(lutronDev)
			areas = append(areas, newArea)
		}
	}
	return areas
}

func (a *Area) addLutron(device *lutron.LDevice) {
	newDev := AreaDevice{
		AreaName:    device.AreaName,
		Id:          strconv.Itoa(device.ID),
		Name:        device.Name,
		Type:        device.Type,
		Value:       strconv.FormatFloat(device.Value, 'E', -1, 64),
		State:       device.State,
		Integration: "lutron",
	}
	if strings.ToLower(device.State) != "off" {
		a.Active = true
	}
	a.Devices = append(a.Devices, newDev)
}

func KasaIntegration(areas []Area, devices kasa.Devices) []Area {
	for _, kasaDev := range devices.Plugs {
		foundArea := false
		areaName := strings.TrimSpace(kasaDev.AreaName)
		if areaName == "" {
			areaName = "UKN"
		}
		for area := range areas {
			if strings.EqualFold(areas[area].AreaName, areaName) {
				areas[area].addKasa(kasaDev)
				foundArea = true
			}
		}
		if !foundArea {
			newArea := Area{
				AreaName: areaName,
			}
			newArea.addKasa(kasaDev)
			areas = append(areas, newArea)
		}
	}
	return areas
}

func (a *Area) addKasa(device kasa.Plug) {
	state := "off"
	if device.PlugInfo.On {
		state = "on"
		a.Active = true
	}
	newDev := AreaDevice{
		AreaName:    device.AreaName,
		Id:          device.IPAddress,
		Name:        device.Name,
		Type:        device.Type,
		State:       state,
		Integration: "kasa",
	}
	a.Devices = append(a.Devices, newDev)
}
