package brain

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/insomniadev/martian/internal/database"
)

type DeviceExpiration struct {
	DeviceLabel  string `json:"device"`
	DeactivateAt string `json:"deactivate"`
}

// Return all of the devices that are set to expire
func (b *Brain) GetDevicesSetToExpire(w http.ResponseWriter, r *http.Request) {
	devices := []DeviceExpiration{}
	for _, dev := range b.Omniscience {
		if dev.EnergyTracked {
			devices = append(devices, DeviceExpiration{
				DeviceLabel:  dev.DeviceId,
				DeactivateAt: dev.MemoryExpiration.Format("2006-01-02 15:04:05"),
			})
		}
	}
	formattedJson, err := json.Marshal(devices)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(formattedJson)
}

// SetAutomation will set an automation in the database
func (b *Brain) SetAutomation(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var newAutomationGraph database.DeviceGraph
	err = json.Unmarshal(body, &newAutomationGraph)
	if err != nil {
		log.Println("Incorrect input: ", err)
		return
	}

	updated := database.MartianData.UpdateGraphTableWithAutomated(newAutomationGraph)

	if updated {
		w.Write([]byte("Updated Successfully"))
	} else {
		w.Write([]byte("Failed to Automate"))
	}
}

// GetAutomatedGraphs will return all of the stored automations
func (b *Brain) GetAutomatedGraphs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if id == "" {
		response := database.MartianData.GetAllAutomatedGraphs()
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(jsonResponse)
	} else {
		response, err := database.MartianData.GetDeviceGraphValues(id)
		if err != nil {
			log.Fatal(err)
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(jsonResponse)
	}
}

// GetDevices will return all of the stored devices
func (b *Brain) GetDevices(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if id == "" {
		response := database.MartianData.GetAllDevices()
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(jsonResponse)
	} else {
		response := database.MartianData.GetDeviceByHash(id)
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(jsonResponse)
	}
}

// GetGraphs will return all of the stored graphs
func (b *Brain) GetGraphs(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	if len(body) == 0 {
		response := database.MartianData.GetAllGraphs()
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(jsonResponse)
	} else {
		// TODO: In the future set it so that we can pull a unique devicehash
	}
}

// UpdateEnergyEfficiency will update the energy efficiency
func (b *Brain) UpdateEnergyEfficiency(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	newTime := params["newValue"]
	intTime, err := strconv.Atoi(newTime)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
	}

	updated := database.MartianData.UpdateEfficiencyTime(id, intTime)
	if !updated {
		w.WriteHeader(400)
	}
	w.WriteHeader(200)
}

// DeleteAutomation is an API to remove the automation provided
func (b *Brain) DeleteAutomation(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var graphToDelete database.DeviceGraph
	if err = json.Unmarshal(body, &graphToDelete); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}

	deleted, err := database.MartianData.DeleteGraphValue(graphToDelete)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	} else if !deleted {
		w.Write([]byte("Failed to delete"))
		w.WriteHeader(400)
	} else {
		w.WriteHeader(200)
	}
}

// GetOmniscience will return the time memory
func (b *Brain) GetOmniscience(w http.ResponseWriter, r *http.Request) {
	if len(b.Omniscience) > 0 {
		switch responseData, err := json.Marshal(b.Omniscience); {
		case err == nil:
			w.Write(responseData)
		default:
			w.Write([]byte(err.Error()))
			w.WriteHeader(400)
		}
	}
}

// GetTimeMemory will return the time memory
func (b *Brain) GetTimeMemory(w http.ResponseWriter, r *http.Request) {
	if len(b.Omniscience) > 0 {
		timeMemories := []Event{}
		for i := range b.Omniscience {
			if b.Omniscience[i].TimeTracked {
				timeMemories = append(timeMemories, b.Omniscience[i])
			}
		}
		switch responseData, err := json.Marshal(timeMemories); {
		case err == nil:
			w.Write(responseData)
		default:
			w.Write([]byte(err.Error()))
			w.WriteHeader(400)
		}
	}
}

// GetTimeAutomation will return the time automation
func (b *Brain) GetTimeAutomation(w http.ResponseWriter, r *http.Request) {
	timeTables := database.MartianData.GetAllAutomatedTimeTables()
	if timeTables != nil {
		switch responseData, err := json.Marshal(timeTables); {
		case err == nil:
			w.Write(responseData)
		default:
			w.Write([]byte(err.Error()))
			w.WriteHeader(400)
		}
	}
}

// GetTimeTables will return the time tables
func (b *Brain) GetTimeTables(w http.ResponseWriter, r *http.Request) {
	timeTables := database.MartianData.GetAllTimeTables()
	if timeTables != nil {
		switch responseData, err := json.Marshal(timeTables); {
		case err == nil:
			w.Write(responseData)
		default:
			w.Write([]byte(err.Error()))
			w.WriteHeader(400)
		}
	}
}

// DeleteTimeTables will delete the time tables
func (b *Brain) DeleteTimeTables(w http.ResponseWriter, r *http.Request) {
	if err := database.MartianData.RecreateTimeBucket(); err != nil {
		w.WriteHeader(400)
	}
}
