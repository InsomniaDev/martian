package brain

import (
	"fmt"
	"strconv"
	"time"

	"github.com/insomniadev/martian/integrations/config"
)

// Init will start up the brain struct
func (b *Brain) Init() {
	url, port := config.LoadRedis()
	b.redisURL = url
	b.redisPort = port

	go b.timeCleanUp()
	go b.timeAutomation()
}

func (b *Brain) timeCleanUp() {
	for {
		b.cleanTempDatabase()
		// Sleep for a week before cleaning up the database again
		time.Sleep(168 * time.Hour)
	}
}

func (b *Brain) timeAutomation() {
	for {
		time.Sleep(15 * time.Minute)
		currentTimeString := assembleTimeString(time.Now())

		// Check if automation is in place for this time
		b.automationCheck(b.checkForTimeAutomations(currentTimeString))
	}
}

// ProcessEvent processes the event
func (b *Brain) ProcessEvent(e Event) {
	b.LastEvent = b.CurrentEvent
	b.CurrentEvent = e

	// If the LastEvent occurred in the last fifteen seconds
	lastFifteenSeconds := time.Now().Add(-15 * time.Second)
	if b.LastEvent.Time.After(lastFifteenSeconds) {
		// These events occurred together
		b.storeEventGraph()
	} else {
		// Just log this as a time event
		b.storeTimeGraph()
	}

	// FIXME: Currently just stop any automations from happening from an automation; more robust logic required here
	if b.CurrentEvent.ID == b.LastEvent.ID && b.CurrentEvent.Type == b.LastEvent.Type {
		return
	}
	// If the automation occurred in the last fifteen seconds
	if b.automationTime.After(lastFifteenSeconds) {
		for _, aEvent := range b.AutomationEvent {
			if aEvent.ID == b.CurrentEvent.ID {
				// TODO: Need to do a delete query for the permanent automation graph
				return
			}
		}
	}

	// Check if automation is in place for this event
	b.automationCheck(b.checkForEventAutomations())
}

func (b *Brain) automationCheck(automationExists bool) {
	if automationExists {
		b.automationTime = time.Now()
		// Fire off the automation since it exists
		for _, aEvent := range b.AutomationEvent {
			if aEvent.ID != b.CurrentEvent.ID {
				if aEvent.Type == "lutron" {
					if s, err := strconv.ParseFloat(aEvent.Value, 64); err == nil {
						fmt.Println(s)
						// graphql.Integrations.LutronData.SetById(aEvent.ID, s)
					}
				}
			}
		}
	}
}
