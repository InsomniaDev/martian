package area

import (
	"strings"

	"github.com/insomniadev/martian/integrations/kasa"
)

func (a *Area) AddLutron() {

}

func KasaIntegration(areas []Area, devices kasa.Devices) []Area {
	for _, kasaDev := range devices.Plugs {
		foundArea := false
		if kasaDev.AreaName == "" {
			kasaDev.AreaName = "UKN"
		}
		for area := range areas {
			if strings.EqualFold(areas[area].AreaName, kasaDev.AreaName) {
				areas[area].AddKasa(kasaDev)
				foundArea = true
			}
		}
		if !foundArea {
			newArea := Area{
				AreaName: kasaDev.AreaName,
			}
			newArea.AddKasa(kasaDev)
			areas = append(areas, newArea)
		}
	}
	return areas
}

func (a *Area) AddKasa(device kasa.Plug) {
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
