package area

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/insomniadev/martian/database"
	"github.com/insomniadev/martian/integrations/kasa"
	"github.com/insomniadev/martian/integrations/lutron"
)

func Init(configuration string) []Area {
	var areas []Area
	json.Unmarshal([]byte(configuration), &areas)
	return areas
}

// InsertAreaIndex will insert an index or priority number for the area
func InsertAreaIndex(areas []Area, areaIndex Area) ([]Area, error) {
	found := false
	var storedAreas []Area
	for i := range areas {
		if areas[i].AreaName == areaIndex.AreaName {
			newArea := Area{
				AreaName: areaIndex.AreaName,
				Index:    areaIndex.Index,
			}
			areas[i].Index = areaIndex.Index
			storedAreas = append(storedAreas, newArea)
			found = true
		} else {
			newArea := Area{
				AreaName: areas[i].AreaName,
				Index:    areas[i].Index,
			}
			storedAreas = append(storedAreas, newArea)
		}
	}
	if !found {
		areas = append(areas, areaIndex)
	}
	var db database.Database
	err := db.PutIntegrationValue("area", storedAreas)
	if err != nil {
		return nil, err
	}
	return areas, nil
}

// CheckIndexForAreas will go through the available areas and apply indexes to those that have them
func CheckIndexForAreas(areas []Area, areaWithIndexes []Area) []Area {
	for _, areaIndex := range areaWithIndexes {
		for i := range areas {
			if areas[i].AreaName == areaIndex.AreaName {
				areas[i].Index = areaIndex.Index
				break
			}
		}
	}
	for i := range areas {
		if areas[i].Index == 0 {
			areas[i].Index = 999
		}
	}
	// Sort by index here and then return the sorted indexes
	// TODO: determine which way this is sorting
	sort.Slice(areas[:], func(i, j int) bool {
		return areas[i].Index < areas[j].Index
	})
	return areas
}

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
