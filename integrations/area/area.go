package area

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/insomniadev/martian/database"
	"github.com/insomniadev/martian/integrations/harmony"
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

func LutronIntegration(areas []Area, devices []*lutron.LDevice, interfaceDevices []int) []Area {
	for _, lutronDev := range devices {
		for _, iDevice := range interfaceDevices {
			// Match the provided interface device id with the inventory device id
			if lutronDev.ID == iDevice {
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

func KasaIntegration(areas []Area, devices kasa.Devices, interfaceDevices []string) []Area {
	for _, kasaDev := range devices.Devices {
		for _, iDevice := range interfaceDevices {
			// Match the provided interface device IPAddress with the inventory device IPAddress
			if kasaDev.IPAddress == iDevice {
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
		}
	}
	return areas
}

func (a *Area) addKasa(device kasa.KasaDevice) {
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

// HarmonyIntegration will grab the harmony integration and return it to the calling application
func HarmonyIntegration(areas []Area, device harmony.Device) []Area {
	foundArea := false
	areaName := strings.TrimSpace(device.AreaName)
	if areaName == "" {
		areaName = "UKN"
	}
	for area := range areas {
		if strings.EqualFold(areas[area].AreaName, areaName) {
			areas[area].addHarmony(device)
			foundArea = true
		}
	}
	if !foundArea {
		newArea := Area{
			AreaName: areaName,
		}
		newArea.addHarmony(device)
		areas = append(areas, newArea)
	}
	return areas
}

// addHarmony will attach the harmony instance to the room
func (a *Area) addHarmony(device harmony.Device) {
	// TODO: Will need to replace this once we have something more concrete in the future, or we do actions against harmony
	type Activities struct {
		ActivityID string `json:"activityID"`
		Name       string `json:"name"`
	}
	var activities []Activities
	for i := range device.Activities {
		newDev := Activities{
			ActivityID: device.Activities[i].ActivityID,
			Name:       device.Activities[i].Name,
		}
		activities = append(activities, newDev)
	}
	if device.CurrentActivity != "-1" {
		a.Active = true
	}
	jsonActivities, _ := json.Marshal(activities)
	newDev := AreaDevice{
		AreaName:    device.AreaName,
		Id:          device.IPAddress,
		Name:        string(jsonActivities),
		Type:        "tv",
		State:       device.CurrentActivity,
		Integration: "harmony",
	}
	a.Devices = append(a.Devices, newDev)
}
