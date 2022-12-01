package brain

import (
	"fmt"
	"time"

	"github.com/insomniadev/martian/internal/database"
	"github.com/insomniadev/martian/pkg/pubsub"
	log "github.com/sirupsen/logrus"
)

type automationEvent struct {
	graph          database.DeviceGraph
	expirationTime time.Time
}

const checkBySeconds = 15

// TODO: Add functionality for dusk to dawn timers

// learnedBrainAction will check to see if there is an automation in place and then update accordingly if there is
func (b *Brain) learnedBrainAction(event Event) {
	// Get any automated graphs for the function
	automatedGraphs := database.MartianData.GetAutomatedGraphs(event.DeviceId, event.Integration, event.Status)

	for _, automatedGraph := range automatedGraphs {
		// Get the device from the database for the integration type
		device := database.MartianData.GetDeviceByHash(automatedGraph.ToUniqueHash)
		fromDevice := database.MartianData.GetDeviceByHash(automatedGraph.FromUniqueHash)

		// Publish the required action
		pubsub.Service.Publish(device.Integration, fmt.Sprintf("%s;;%s", device.DeviceId, automatedGraph.ToStatus))

		// print log statement for tracking to start
		log.Debugln("automated the start of: ", fmt.Sprintf("%s;;%s;;%s", device.Integration, device.Label, automatedGraph.ToStatus))

		b.AutomationMemory = append(b.AutomationMemory, automationEvent{
			expirationTime: time.Now().Add(checkBySeconds * time.Second),
			graph:          automatedGraph,
		})

		// reset the energy savings time if device is already set to that state
		b.rememberEnergyTime(Event{Integration: device.Integration, DeviceId: device.DeviceId, TriggeredDeviceId: fromDevice.DeviceId, TriggeredDeviceStatus: automatedGraph.FromStatus})
	}

	// Check to see if we should adjust the weight for the automation if we are immediatly changing status of automated device
	b.checkUnsetAutomation(event)
}

// checkUnsetAutomation will determine if the event called is conflicting with the automation that just occurred
func (b *Brain) checkUnsetAutomation(event Event) {
	// things are becoming unautomated too quickly, changing this to -1 to see if that helps with automations
	unsetWeightBy := -1

	currentTime := time.Now()

	// Remove the already expired events
	// TODO: Pull this logic into a common function, exists in energyefficiency.go as well
	eventsToCheckAgainst := []automationEvent{}
	for i := range b.AutomationMemory {
		if b.AutomationMemory[i].expirationTime.After(currentTime) {
			eventsToCheckAgainst = append(eventsToCheckAgainst, b.AutomationMemory[i])
		}
	}

	// Update array with valid events
	b.AutomationMemory = eventsToCheckAgainst

	partialGraph, _ := database.AssembleDeviceGraph("", "", "", event.DeviceId, event.Integration, event.Status)

	for i := range b.AutomationMemory {
		// If it is the same device with a different status,
		if b.AutomationMemory[i].graph.ToUniqueHash == partialGraph.ToUniqueHash && b.AutomationMemory[i].graph.ToStatus != partialGraph.ToStatus {
			// then decrease the weight by the unsetWeightBy

			database.MartianData.UpdateGraphWeight(b.AutomationMemory[i].graph, unsetWeightBy)
		}
	}
}
