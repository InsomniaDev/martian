package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/insomniadev/martian/internal/brain"
	"github.com/insomniadev/martian/internal/integrations/hubitat"
)

func StartApi() {

	r := mux.NewRouter()
	r.HandleFunc("/hubitat", hubitat.Instance.GetDatas)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	// TODO: make this cleaner in the future?
	// Get insight into Martian
	r.HandleFunc("/omniscience", brain.Brainiac.GetOmniscience)
	r.HandleFunc("/getTimeExpirations", brain.Brainiac.GetDevicesSetToExpire)
	r.HandleFunc("/setAutomation", brain.Brainiac.SetAutomation)
	r.HandleFunc("/automations", brain.Brainiac.GetAutomatedGraphs).Methods("GET")
	r.HandleFunc("/automations/{id}", brain.Brainiac.GetAutomatedGraphs).Methods("GET")
	r.HandleFunc("/delete/automations", brain.Brainiac.DeleteAutomation)
	r.HandleFunc("/getGraphs", brain.Brainiac.GetGraphs)
	r.HandleFunc("/devices", brain.Brainiac.GetDevices).Methods("GET")
	r.HandleFunc("/devices/{id}", brain.Brainiac.GetDevices).Methods("GET")
	r.HandleFunc("/energy/{id}/{newValue}", brain.Brainiac.UpdateEnergyEfficiency).Methods("GET")

	r.HandleFunc("/times/memory", brain.Brainiac.GetTimeMemory)
	r.HandleFunc("/times/automation", brain.Brainiac.GetTimeAutomation)
	r.HandleFunc("/times/tables", brain.Brainiac.GetTimeTables)
	r.HandleFunc("/times/delete", brain.Brainiac.DeleteTimeTables)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8088", nil))
}
