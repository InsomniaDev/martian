package brain

import (
	"encoding/json"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/insomniadev/martian/internal/database"
	"github.com/insomniadev/martian/pkg/cron"
	"github.com/insomniadev/martian/pkg/pubsub"
)

type Device struct {
	DeviceId                string
	UniqueHash              string
	Name                    string
	Label                   string
	DeviceType              string
	Model                   string
	Integration             string
	EnergyEfficiencyMinutes string
}

type Brain struct {
	Omniscience      []Event           `json:"omniscience"`
	LastEvent        Event             `json:"lastEvent"`
	AutomationMemory []automationEvent `json:"automationMemory"`
}

type Event struct {
	Integration           string    `json:"eventType"`
	EventTime             time.Time `json:"eventTime"`
	DeviceId              string    `json:"DeviceId"`
	Status                string    `json:"status"`
	FullEvent             string    `json:"fullEvent"`
	TriggeredDeviceId     string    `json:"triggeredDeviceId"`
	TriggeredDeviceStatus string    `json:"triggeredDeviceStatus"`
	TimeAutomated         bool      `json:"timeAutomated"`
	TimeTracked           bool      `json:"timeTracked"`
	TimeExpiration        time.Time `json:"timeExpiration"`
	EnergyTracked         bool      `json:"energyTracked"`
	EnergyExpiration      time.Time `json:"energyExpiration"`
	MemoryTracked         bool      `json:"memoryTracked"`
	MemoryExpiration      time.Time `json:"memoryExpiration"`

	Consecutive ConsecutiveEvent `json:"consecutive"`
}

type ConsecutiveEvent struct {
	DeviceId    string `json:"DeviceId"`
	Status      string `json:"status"`
	Integration string `json:"eventType"`
}

type LongTermStore struct {
	HourOfOccurrence int       `json:"hourOfOccurrence"`
	EventType        string    `json:"eventType"`
	EventTime        time.Time `json:"eventTime"`
	DeviceId         string    `json:"DeviceId"`
	Status           string    `json:"status"`
	SequentialEvents []Event   `json:"sequentialEvents"`
	FullEvent        string    `json:"fullEvent"`
}

var (
	Brainiac       *Brain
	timeDifference = 30 * time.Second // Only remember events for 30 seconds

	// Values for remembering and updating energy efficiency
	lastTimeShutoffUniqueHash = ""
	lastTimeShutOffDevices    = []efficiencyTimeCheck{}
)

func (b *Brain) SayHello() {
	log.Info("The Brain is ALIVE - Hello!")
}

// learnNewDevices will put new devices into the database
func learnNewDevices() {
	subscriptionBus := make(chan string)
	var subscribeToEvents = func() {
		for {
			msg := <-subscriptionBus
			log.Debug("Received message: ", msg)

			// deviceId ;; name ;; label ;; type ;; model ;; manufacturer ;; capabilities ;; attributes ;; commands ;; integration
			message := strings.Split(msg, ";;")

			newDevice := database.Device{
				DeviceId:    message[0],
				Name:        message[1],
				Label:       message[2],
				Type:        message[3],
				Integration: message[9],
			}

			inserted := database.MartianData.SetDevice(newDevice)
			if !inserted {
				log.Debug(message[2] + " was not inserted")
			}
		}
	}
	pubsub.Service.Subscribe("learnNewDevice", subscriptionBus)
	go subscribeToEvents()
}

func init() {
	Brainiac = &Brain{}

	// attempt to pull from the brain if it exists
	if exists, value := database.MartianData.GetMemoryData(); exists {
		if err := json.Unmarshal(value, &Brainiac); err != nil && value != nil {
			log.Fatal("Error when pulling in the memory data")
		}
	}

	learnNewDevices()

	subscriptionBus := make(chan string)
	var subscribeToEvents = func() {
		for {
			msg := <-subscriptionBus
			log.Debug("Received message: ", msg)

			// integration ;; unique ID ;; label ;; attribute type ;; attribute value ;; activated
			message := strings.Split(msg, ";;")
			log.Debugln(message)
			uniqueHash := database.GetUniqueHash(message[1], message[0])

			// Check if device was turned off by the energy efficiency learning
			if uniqueHash != lastTimeShutoffUniqueHash {
				Brainiac.brainWave(message[0], message[1], message[4], msg)
			} else {
				// Reset the shutoff hash so that we can continue to learn
				lastTimeShutoffUniqueHash = ""
			}
		}
	}
	pubsub.Service.Subscribe("brain", subscriptionBus)
	go subscribeToEvents()

	Brainiac.setTimeAutomations()
	c := cron.New()
	c.AddFunc("@every 10s", func() {
		log.Debugln("energyEfficiency")
		Brainiac.energyEfficiency()
	})
	c.AddFunc("*/1 * * * *", func() { // at every minute
		log.Debugln("shortTerm")
		Brainiac.shortTerm()
	})
	c.AddFunc("*/3 * * * *", func() { // every five minutes
		// Request all integrations to update the status for their connected devices
		pubsub.Service.Publish("pulse", "true")
	})
	c.AddFunc("0 0 * * *", func() { // at midnight every day
		log.Debugln("processDayMemories")
		Brainiac.processDayMemories()
	})
	c.AddFunc("0 0 * * 6", func() { // 12am on Saturday
		log.Debugln("setTimeAutomations")
		Brainiac.setTimeAutomations()
	})
	c.AddFunc("*/10 * * * *", func() { // every ten minutes
		log.Println("checkTimeAutomations")
		Brainiac.checkTimeAutomations()
	}) // @every 10m
	c.Start()
}

// StoreMemoryData will store memory data into the database
func (b *Brain) StoreMemoryData() bool {
	return database.MartianData.StoreMemoryData(b)
}

// brainWave will add events into the brain with a populated eventTimeExpiration
func (b *Brain) brainWave(integrationType, id, status, fullEvent string) {

	// Turn this into the event type
	newEvent := Event{Integration: integrationType, DeviceId: id, Status: status, EventTime: time.Now(), MemoryExpiration: time.Now().Add(timeDifference), FullEvent: fullEvent}

	// Update the memory of the device here...
	deviceKnown := false
	omniscienceIndex := 0
	for i := range b.Omniscience {
		if b.Omniscience[i].DeviceId == id && b.Omniscience[i].Integration == integrationType {

			deviceKnown = true

			// constantly update if time tracked
			b.checkForTimeActivation(newEvent, fullEvent)

			// update if this is a different status
			if b.Omniscience[i].Status != status {
				// This is an updated status
				b.Omniscience[i].EventTime = newEvent.EventTime
				b.Omniscience[i].Status = newEvent.Status
				b.Omniscience[i].FullEvent = newEvent.FullEvent
				b.Omniscience[i].MemoryExpiration = newEvent.MemoryExpiration
				omniscienceIndex = i

				break
			}
			return
		}
	}

	for i := range b.Omniscience {
		// Update this as the consecutive device since it is not related to energy and is
		if b.Omniscience[i].DeviceId == b.LastEvent.DeviceId && b.Omniscience[i].Integration == b.LastEvent.Integration {
			b.Omniscience[i].Consecutive.DeviceId = newEvent.DeviceId
			b.Omniscience[i].Consecutive.Integration = newEvent.Integration
			b.Omniscience[i].Consecutive.Status = newEvent.Status
		}
	}

	if !deviceKnown {
		b.Omniscience = append(b.Omniscience, newEvent)
	}
	b.LastEvent = newEvent

	b.checkForTimeActivation(newEvent, fullEvent)

	// Check if the device was recently turned off and update the period of time the device should be on if that is the case, returns a boolean on if the device was just turned back on quickly after going off
	deviceUpdateWasEnergyRelated := checkForDeviceEnergyEfficiency(database.GetUniqueHash(newEvent.DeviceId, newEvent.Integration))

	// add the event to the short term memory
	// ignore the event if someone was just turning the light back on after it turned off
	if !deviceUpdateWasEnergyRelated {

		// Remember this as tracked by the active memory
		b.Omniscience[omniscienceIndex].MemoryTracked = true
		log.Debugln("stored event:", integrationType, id, status)
	}

	// Remember the event in relation to the energy savings
	b.rememberEnergyTime(newEvent)

	// determine if there is a learned behavior
	go b.learnedBrainAction(newEvent)
}

// shortTerm checks to see if the current timestamp is greater than the eventTimeExpiration
// IF time is greater than the memoryEvent is removed from the active array
func (b *Brain) shortTerm() {

	for i := range b.Omniscience {
		if b.Omniscience[i].MemoryTracked && time.Now().After(b.Omniscience[i].MemoryExpiration) {
			b.Omniscience[i].MemoryTracked = false

			// If there isn't a consecutive event then just return
			if b.Omniscience[i].Consecutive.DeviceId == "" {
				return
			}

			// We only want to store the current event and the consequtive event, all other events don't matter
			// Attempt to retrieve existing graph
			existingGraph, err := database.MartianData.GetGraphByRelationship(b.Omniscience[i].DeviceId, b.Omniscience[i].Integration, b.Omniscience[i].Status, b.Omniscience[i].Consecutive.DeviceId, b.Omniscience[i].Consecutive.Integration, b.Omniscience[i].Consecutive.Status)

			//  If there is a returned error then the graph doesn't exist and we need to create it
			if err != nil {
				// Create a device_graph from that existing information
				existingGraph, sameDevice := database.AssembleDeviceGraph(
					b.Omniscience[i].DeviceId, b.Omniscience[i].Integration, b.Omniscience[i].Status,
					b.Omniscience[i].Consecutive.DeviceId, b.Omniscience[i].Consecutive.Integration, b.Omniscience[i].Consecutive.Status,
				)
				// for now only insert
				if !sameDevice {
					database.MartianData.SetGraph(existingGraph)
				}
			} else {
				database.MartianData.UpdateGraphWeight(existingGraph, 1)
			}
		}
	}
}

// processDayMemories processes the memories that occur every morning at 1am
func (b *Brain) processDayMemories() {

	// Since we are resetting every 24hrs then lets set the weight lower
	// TODO: This should be a configuration in the future
	// checking 5 for making events require more occurrence before becoming automated
	// previous (3) was more aggressive and there were a lot of false automations getting put into place
	setWeight := 5

	// Get the records that we are going to automate that have a weight above setWeight
	recordsToAutomate := database.MartianData.RetrieveDataFromDeviceGraphByWeight(setWeight)

	// Mark these items as automated and set weight equal to zero
	for _, record := range recordsToAutomate {
		// if record is already automated, then don't reset the weight
		if !record.Automated {
			database.MartianData.UpdateGraphTableWithAutomated(record)
		}
	}

	// Reset all records that did not occur more than setWeight times
	database.MartianData.ResetDeviceGraph()
	// }
}
