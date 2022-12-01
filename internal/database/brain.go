package database

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func hashFunction(hashData ...string) string {
	hashDataString := ""
	for _, ds := range hashData {
		hashDataString += ds
	}

	return fmt.Sprintf("%x", md5.Sum([]byte(hashDataString)))
}

var badDevices = []string{"Fibaro Motion Sensor ZW5", "Generic Zigbee Motion Sensor"}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// GetUniqueHash will return the unique hash created from the deviceid and integration
func GetUniqueHash(deviceId, integration string) string {
	return hashFunction(deviceId, integration)
}

// AssembleDeviceGraph will create the graph sructure required to insert into table, will return a boolean on if the device_graph has the same device in the from->to positions
func AssembleDeviceGraph(fromId, fromIntegration, fromStatus, toId, toIntegration, toStatus string) (deviceGraph DeviceGraph, sameDevice bool) {

	if hashFunction(fromId, fromIntegration) == hashFunction(toId, toIntegration) {
		sameDevice = true
	} else {
		sameDevice = false
	}

	return DeviceGraph{
		FromUniqueHash: hashFunction(fromId, fromIntegration),
		FromStatus:     fromStatus,
		ToUniqueHash:   hashFunction(toId, toIntegration),
		ToStatus:       toStatus,
		Weight:         1,
		Automated:      false,
		TimeAutomated:  false,
	}, sameDevice
}

// DeleteGraphValue will remove the graph provided
func (d *Database) DeleteGraphValue(graphToDelete DeviceGraph) (updated bool, err error) {
	// Grab all of the graphs that match the fromHash
	uniqueGraphs, err := d.GetDeviceGraphValues(graphToDelete.FromUniqueHash)
	if err != nil {
		return false, err
	}

	// Create a new array minus the graph listed above
	newUniqueGraph := UniqueHashGraphs{}
	for _, graph := range uniqueGraphs.Graphs {
		if graph.FromStatus == graphToDelete.FromStatus && graph.ToUniqueHash == graphToDelete.ToUniqueHash && graph.ToStatus == graphToDelete.ToStatus {
			// We don't want to do anything with this one since we are deleting it
			updated = true
		} else {
			newUniqueGraph.Graphs = append(newUniqueGraph.Graphs, graph)
		}
	}

	// Update that  uniquegraph
	if err = d.PutDeviceGraphValues(newUniqueGraph); err != nil {
		return false, err
	}
	return updated, nil
}

// GetDevice will return the device row for the provided uniqueId and integration
func (d *Database) GetDevice(uniqueId, integration string) (device Device) {
	return d.GetDeviceByHash(hashFunction(uniqueId, integration))
}

// GetDeviceByHash will return the device row for the provided uniqueId and integration
func (d *Database) GetDeviceByHash(uniqueHash string) (device Device) {
	device, err := d.GetDeviceValues(uniqueHash)
	if err != nil {
		log.Fatal(err)
	}
	return device
}

// GetGraphs will return all of the graphs for the provided uniqueid and integration
func (d *Database) GetGraphs(uniqueId, integration string) (graphResults []DeviceGraph) {
	graphs, err := d.GetDeviceGraphValues(hashFunction(uniqueId, integration))
	if err != nil {
		log.Fatal(err)
	}
	return graphs.Graphs
}

// GetAutomatedGraphs will return all of the graphs that are automated for the provided uniqueid and integration
func (d *Database) GetAutomatedGraphs(uniqueId, integration, status string) (graphResults []DeviceGraph) {
	graphs, err := d.GetDeviceGraphValues(hashFunction(uniqueId, integration))
	if err != nil {
		log.Fatal(err)
	}

	graphResults = []DeviceGraph{}
	for _, graph := range graphs.Graphs {
		if graph.FromStatus == status && graph.Automated && !graph.TimeAutomated {
			graphResults = append(graphResults, graph)
		}
	}

	return graphResults
}

// GetTimeTableEntry will return the devices that are available for that time
func (d *Database) GetAutomatedTimeTableEntry(timeBlock string) []Device {

	if timeTables, err := d.GetTimeTableValues(timeBlock); err == nil {
		devices := []Device{}
		for _, tblock := range timeTables.Times {
			if tblock.Automated {
				if device, err := d.GetDeviceValues(tblock.UniqueHash); err == nil {
					devices = append(devices, device)
				}
			}
		}

		return devices
	}
	return nil
}

// GetGraphByRelationship gets the graph by the first and last pieces provided
func (d *Database) GetGraphByRelationship(firstId, firstIntegration, firstStatus, lastId, lastIntegration, lastStatus string) (graph DeviceGraph, err error) {
	graphs := d.GetGraphs(firstId, firstIntegration)
	toUniqueHash := GetUniqueHash(lastId, lastIntegration)
	for _, graph := range graphs {
		if graph.ToUniqueHash == toUniqueHash {
			if graph.FromStatus == firstStatus && graph.ToStatus == lastStatus {
				return graph, nil
			}
		}
	}

	return graph, errors.New("graph doesn't exist")
}

// SetDevice will set the device in the database, will return a false if it was unable to insert or device already exists
func (d *Database) SetDevice(data Device) (inserted bool) {

	// Create the unique hash
	data.UniqueHash = hashFunction(data.DeviceId + data.Integration)

	// Check if exists
	deviceExists := d.GetDeviceByHash(data.UniqueHash)
	if deviceExists.DeviceId == "" {
		data.EnergyEfficiencyMinutes = 15
	} else {
		data.EnergyEfficiencyMinutes = deviceExists.EnergyEfficiencyMinutes
	}

	if err := d.PutDeviceValues(data); err != nil {
		return false
	}
	return true
}

// SetGraph Will set the graph for the first time
func (d *Database) SetGraph(graphData DeviceGraph) bool {
	// Get the graphs first
	potentialGraphs, err := d.GetDeviceGraphValues(graphData.FromUniqueHash)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the device is one that should be added
	toDevice, err := d.GetDeviceValues(graphData.ToUniqueHash)
	if err != nil {
		log.Fatal(err)
	}
	if contains(badDevices, toDevice.Type) {
		return false
	}

	potentialGraphs.Graphs = append(potentialGraphs.Graphs, graphData)
	d.PutDeviceGraphValues(potentialGraphs)

	return true
}

// UpdateWeightsForTimeTable will update the weights for all of the devices currently active
func (d *Database) UpdateWeightsForTimeTable(timeTable string, timeDevices []Device) bool {
	// Get all of the current entries in the timetable
	if timeEntries, err := d.GetTimeTableValues(timeTable); err == nil {

		if data, err := json.Marshal(timeEntries); err == nil {
			log.Debugln("checking stored data: ", string(data))
		}

		// cycle through the stored entries and update weights in-place
		for td := range timeDevices {

			// Get the unique hash from the provided device
			uniqueHash := GetUniqueHash(timeDevices[td].DeviceId, timeDevices[td].Integration)

			// Go through and increase the weights if it already exists
			existsInStorage := false
			for i := range timeEntries.Times {
				if timeEntries.Times[i].UniqueHash == uniqueHash {
					// update the weight
					timeEntries.Times[i].Weight = timeEntries.Times[i].Weight + 1
					existsInStorage = true
					break
				}
			}

			// Add to the stored time entries with an initial weight of 1 and not automated
			if !existsInStorage {
				timeEntries.Times = append(timeEntries.Times, TimeBlocks{
					TimeKey:    timeTable,
					UniqueHash: uniqueHash,
					Weight:     1,
					Automated:  false,
				})
			}
		}

		// update the database with the weights
		d.PutTimeTableValues(timeTable, timeEntries)

		if data, err := json.Marshal(timeEntries); err == nil {
			log.Debugln("storing this data: ", string(data))
		}
	}
	return false
}

// UpdateEfficiencyTime will update the efficiency time in the device table for the provided hash by provided amount
func (d *Database) UpdateEfficiencyTime(uniquehash string, updateBy int) bool {

	device := d.GetDeviceByHash(uniquehash)
	device.EnergyEfficiencyMinutes += updateBy
	d.PutDeviceValues(device)
	return true
}

// UpdateGraphWeight will update the weight on the graph entry
func (d *Database) UpdateGraphWeight(graphData DeviceGraph, amount int) bool {
	// attempting to update weight
	boolUpdated := false
	graphs, err := d.GetDeviceGraphValues(graphData.FromUniqueHash)
	if err != nil {
		log.Fatal(err)
		return false
	}

	log.Debugln(graphs)
	for i, graph := range graphs.Graphs {
		if graph.FromStatus == graphData.FromStatus && graph.ToUniqueHash == graphData.ToUniqueHash && graph.ToStatus == graphData.ToStatus {
			graphs.Graphs[i].Weight = graphs.Graphs[i].Weight + amount
			boolUpdated = true
		}
	}

	d.PutDeviceGraphValues(graphs)
	return boolUpdated
}

// ImproveEnergyEfficiencyDaily will decrease the time that devices are on by a minute per day
func (d *Database) ImproveEnergyEfficiency(amount int) bool {

	var contains = func(s []int, e int) bool {
		for _, a := range s {
			if a == e {
				return true
			}
		}
		return false
	}

	devices := d.GetAllDevices()
	energyEfficiencySpecialValues := []int{-1, 0, 1, 2, 3}
	for i := range devices {
		if !contains(energyEfficiencySpecialValues, devices[i].EnergyEfficiencyMinutes) {
			devices[i].EnergyEfficiencyMinutes -= 1
			d.PutDeviceValues(devices[i])
		}
	}
	return true
}

// RetrieveDataFromDeviceGraph grabs data from the device graph above the specified weight amount
func (d *Database) RetrieveDataFromDeviceGraphByWeight(weightAmount int) (graphResults []DeviceGraph) {

	graphs := d.GetAllDeviceGraphs()
	for i := range graphs {
		if graphs[i].Weight >= weightAmount {
			graphResults = append(graphResults, graphs[i])
		}
	}
	return graphResults
}

// UpdateGraphTableWithAutomated will set entries as automated
func (d *Database) UpdateGraphTableWithAutomated(graphData DeviceGraph) bool {

	graphs, err := d.GetDeviceGraphValues(graphData.FromUniqueHash)
	if err != nil {
		log.Fatal(err)
	}
	existing := false
	for i, graph := range graphs.Graphs {
		if graph.FromUniqueHash == graphData.FromUniqueHash && graph.FromStatus == graphData.FromStatus && graph.ToUniqueHash == graphData.ToUniqueHash && graph.ToStatus == graphData.ToStatus {
			graphs.Graphs[i].Automated = true
			graphs.Graphs[i].Weight = 0
			existing = true
			break
		}
	}

	// If we are creating an automation that doesn't already exist
	if !existing {
		graphData.Automated = true
		graphData.Weight = 0

		graphs.Graphs = append(graphs.Graphs, graphData)
	}

	if err := d.PutDeviceGraphValues(graphs); err != nil {
		return false
	}

	return true
}

// GetAllAutomatedGraphs will return all of the automated graphs
func (d *Database) GetAllAutomatedGraphs() (automatedGraphs []DeviceGraph) {
	graphs := d.GetAllGraphs()

	for _, uniqueGraph := range graphs {
		for i := range uniqueGraph.Graphs {
			if uniqueGraph.Graphs[i].Automated && !uniqueGraph.Graphs[i].TimeAutomated {
				automatedGraphs = append(automatedGraphs, uniqueGraph.Graphs[i])
			}
		}
	}
	return
}

// GetAllTimeTables will return all of the timetables
func (d *Database) GetAllAutomatedTimeTables() []TimeTable {
	timeTables := d.getAllTimeBuckets()

	returnTimeTables := []TimeTable{}

	for i := range timeTables {
		timeTable := TimeTable{}
		for ti := range timeTables[i].Times {
			if timeTables[i].Times[ti].Automated {
				timeTable.Times = append(timeTable.Times, timeTables[i].Times[ti])
			}
		}
		if len(timeTable.Times) > 0 {
			returnTimeTables = append(returnTimeTables, timeTable)
		}
	}
	return returnTimeTables
}

// GetAllTimeTables will return all of the timetables
func (d *Database) GetAllTimeTables() []TimeTable {
	return d.getAllTimeBuckets()
}

// ResetDeviceGraph will delete all unautomated entries
func (d *Database) ResetDeviceGraph() bool {
	graphs := d.GetAllGraphs()

	// Delete all of the graphs
	err := d.RecreateGraphBucket()
	if err != nil {
		return false
	}

	for _, uniqueGraph := range graphs {
		newGraph := UniqueHashGraphs{}
		for i := range uniqueGraph.Graphs {
			if (uniqueGraph.Graphs[i].Automated || uniqueGraph.Graphs[i].TimeAutomated) && uniqueGraph.Graphs[i].Weight >= 0 {

				// keep this graph
				// reset the weight to zero
				uniqueGraph.Graphs[i].Weight = 0
				newGraph.Graphs = append(newGraph.Graphs, uniqueGraph.Graphs[i])
			}
		}
		// Only keep the graphs that are automated
		if len(newGraph.Graphs) > 0 {
			if err := d.PutDeviceGraphValues(newGraph); err != nil {
				log.Fatal(err)
				return false
			}
		}
	}
	return true
}
