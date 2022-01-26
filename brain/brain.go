package brain

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/insomniadev/martian/database"
	"github.com/insomniadev/martian/modules/pubsub"
)

type Brain struct {
	memoryEvent    []event
	longTermMemory []longTermStore
}

type event struct {
	eventTime           time.Time
	eventTimeExpiration time.Time
	eventId             string
	eventStatus         string
}

type longTermStore struct {
	hourOfOccurrence int
	eventTime        time.Time
	eventId          string
	eventStatus      string
	sequentialEvents []event
}

var (
	Brainiac       *Brain
	timeDifference = 30 * time.Second // Only remember events for 30 seconds
)

func (b *Brain) SayHello() {
	log.Println("Hello")
}

func init() {
	Brainiac = &Brain{}
	subscriptionBus := make(chan string)
	var subscribeToEvents = func() {
		for {
			msg := <-subscriptionBus
			message := strings.Split(msg, ";")
			Brainiac.brainWave(message[0], message[1])
		}
	}
	pubsub.Service.Subscribe("brain", subscriptionBus)
	go subscribeToEvents()

	// Start the short term memory which will delete events
	go Brainiac.shortTerm()

	// Process the daily memories
	go Brainiac.processDayMemories()
}

// brainWave will add events into the brain with a populated eventTimeExpiration
func (b *Brain) brainWave(id, status string) {
	// TODO: Add more logic here, just append for now
	b.memoryEvent = append(b.memoryEvent, event{eventId: id, eventStatus: status, eventTime: time.Now(), eventTimeExpiration: time.Now().Add(timeDifference)})
	fmt.Println("I remember that", id, status)
}

// shortTerm checks to see if the current timestamp is greater than the eventTimeExpiration
// IF time is greater than the memoryEvent is removed from the active array
func (b *Brain) shortTerm() {
	for {
		currentTime := time.Now()
		var memoryEvents []event
		for i := range b.memoryEvent {
			if currentTime.Before(b.memoryEvent[i].eventTimeExpiration) {
				memoryEvents = append(memoryEvents, b.memoryEvent[i])
			} else {
				// get the hour timestamp for the event
				hourOfOccurrence, _, _ := b.memoryEvent[i].eventTime.Clock()

				// go through the rest of the array after the memory event that we are on so that we are only storing the sequential events
				var sequentialEvents []event
				for a := i; a < len(b.memoryEvent); a++ {
					sequentialEvents = append(sequentialEvents, b.memoryEvent[a])
				}

				// create an entry to remember
				remember := longTermStore{hourOfOccurrence: hourOfOccurrence, eventTime: b.memoryEvent[i].eventTime, eventId: b.memoryEvent[i].eventId, eventStatus: b.memoryEvent[i].eventStatus, sequentialEvents: sequentialEvents}
				resp, err := database.MartianData.GetDayMemoryByHour(strconv.Itoa(hourOfOccurrence))
				if err != nil {
					log.Println("Failure to pull memory from short term collection", err)
				}

				// Grab all of the remembered events from this memory store
				var hourEventsRemembered []longTermStore
				if err := json.Unmarshal(resp, &hourEventsRemembered); err != nil {
					log.Println("Error getting short term source from the brain", err)
				}

				// Store away all of the events for the hour here
				hourEventsRemembered = append(hourEventsRemembered, remember)
				database.MartianData.InsertMemory(strconv.Itoa(hourOfOccurrence), hourEventsRemembered)
			}
		}
		b.memoryEvent = memoryEvents
		time.Sleep(1 * time.Minute)
	}
}

// processDayMemories processes the memories that occur throughout the day
func (b *Brain) processDayMemories() {
	time.Sleep(24 * time.Hour)

	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	recipient := os.Getenv("EMAIL_RECIPIENT")

	var emailBody string
	// need to pull the messages for every hour of the day
	for i := 0; i < 24; i++ {
		// Retrieve the data for the hour
		var hourEventsRemembered []longTermStore
		resp, err := database.MartianData.GetDayMemoryByHour(strconv.Itoa(i + 1))
		if err != nil {
			log.Println("processDayMemories: Failed to return hour key for the nightly memory process:", err)
		}
		if err := json.Unmarshal(resp, &hourEventsRemembered); err != nil {
			log.Println("processDayMemories: Error getting short term source from the brain", err)
		}

		// assemble data for the hour
		emailBody += "Hour: " + strconv.Itoa(i+1) + "\n\n"
		for _, memoryRemembered := range hourEventsRemembered {
			emailBody += memoryRemembered.eventId + ":" + memoryRemembered.eventStatus
			for _, memorySequence := range memoryRemembered.sequentialEvents {
				emailBody += " -> " + memorySequence.eventId + ":" + memorySequence.eventStatus
			}
			emailBody += "\n"
		}
		emailBody += "\n\n"
	}

	// email the report to me so that I can see what it looks like and can appropriately develop
	if len(emailBody) > 0 {
		host := "smtp.gmail.com:587"
		auth := smtp.PlainAuth("", username, password, "smtp.gmail.com")

		body := []byte(
			"To: " + recipient + "\r\n" +
				"Subject: Daily Home Status Report\r\n\r\n" +
				emailBody)

		header := make(map[string]string)
		header["From"] = username
		header["To"] = recipient
		header["Subject"] = "Daily Home Status Report"
		// toList is list of email address that email is to be sent.
		toList := []string{recipient}
		// Need to fix how it looks when sent and also add in a subject line
		err := smtp.SendMail(host, auth, "username", toList, body)

		if err != nil {
			log.Println(err)
		}
	}
	go b.processDayMemories()
}
