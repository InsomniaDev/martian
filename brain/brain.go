package brain

import (
	"encoding/json"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/insomniadev/martian/database"
	"github.com/insomniadev/martian/modules/pubsub"
)

type Brain struct {
	MemoryEvent    []event         `json:"memoryEvent"`
	LongTermMemory []longTermStore `json:"longTermMemory"`
}

type event struct {
	EventType           string    `json:"eventType"`
	EventTime           time.Time `json:"eventTime"`
	EventTimeExpiration time.Time `json:"eventTimeExpiration"`
	EventId             string    `json:"eventId"`
	EventStatus         string    `json:"eventStatus"`
}

type longTermStore struct {
	HourOfOccurrence int       `json:"hourOfOccurrence"`
	EventType        string    `json:"eventType"`
	EventTime        time.Time `json:"eventTime"`
	EventId          string    `json:"eventId"`
	EventStatus      string    `json:"eventStatus"`
	SequentialEvents []event   `json:"sequentialEvents"`
}

var (
	Brainiac       *Brain
	timeDifference = 30 * time.Second // Only remember events for 30 seconds
)

func (b *Brain) SayHello() {
	log.Info("The Brain is ALIVE - Hello!")
}

func init() {
	Brainiac = &Brain{}
	subscriptionBus := make(chan string)
	var subscribeToEvents = func() {
		for {
			msg := <-subscriptionBus
			log.Debug("Received message: ",msg)
			message := strings.Split(msg, ";")
			Brainiac.brainWave(message[0], message[1], message[2])
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
func (b *Brain) brainWave(integrationType, id, status string) {
	// TODO: Add more logic here, just append for now
	b.MemoryEvent = append(b.MemoryEvent, event{EventType: integrationType, EventId: id, EventStatus: status, EventTime: time.Now(), EventTimeExpiration: time.Now().Add(timeDifference)})
	log.Debug("stored event:", integrationType, id, status)
}

// shortTerm checks to see if the current timestamp is greater than the eventTimeExpiration
// IF time is greater than the memoryEvent is removed from the active array
func (b *Brain) shortTerm() {
	for {
		currentTime := time.Now()
		var memoryEvents []event
		for i := range b.MemoryEvent {
			if currentTime.Before(b.MemoryEvent[i].EventTimeExpiration) {
				memoryEvents = append(memoryEvents, b.MemoryEvent[i])
			} else {
				// get the hour timestamp for the event
				hourOfOccurrence, _, _ := b.MemoryEvent[i].EventTime.Clock()

				// go through the rest of the array after the memory event that we are on so that we are only storing the sequential events
				var sequentialEvents []event
				for a := i + 1; a < len(b.MemoryEvent); a++ {
					sequentialEvents = append(sequentialEvents, b.MemoryEvent[a])
				}

				// create an entry to remember
				remember := longTermStore{HourOfOccurrence: hourOfOccurrence, EventTime: b.MemoryEvent[i].EventTime, EventId: b.MemoryEvent[i].EventId, EventStatus: b.MemoryEvent[i].EventStatus, EventType: b.MemoryEvent[i].EventType, SequentialEvents: sequentialEvents}
				resp, err := database.MartianData.GetDayMemoryByHour(strconv.Itoa(hourOfOccurrence))
				if err != nil {
					log.Error("Failure to pull memory from short term collection", err)
				}

				// Grab all of the remembered events from this memory store
				var hourEventsRemembered []longTermStore
				if err := json.Unmarshal(resp, &hourEventsRemembered); err != nil {
					log.WithFields(log.Fields{
						"resp": string(resp),
					}).Error("Error getting short term source from the brain ", err)
				}

				// Store away all of the events for the hour here
				hourEventsRemembered = append(hourEventsRemembered, remember)
				database.MartianData.InsertMemory(strconv.Itoa(hourOfOccurrence), hourEventsRemembered)
			}
		}
		b.MemoryEvent = memoryEvents
		time.Sleep(1 * time.Minute)
	}
}

// processDayMemories processes the memories that occur throughout the day
func (b *Brain) processDayMemories() {
	time.Sleep(24 * time.Hour)

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	recipient := os.Getenv("RECIPIENT")

	var emailBody string
	// pull all messages for the past 24 hours
	memories, err := database.MartianData.RetrieveAllMemories()
	if err != nil {
		log.Error("processDayMemories: Error returning memories:", err)
	}

	// go through each memory hour from the last 24 hours
	for key, memory := range memories {
		log.Println(key, string(memory))

		// convert data into struct
		var hourEventsRemembered []longTermStore
		if err := json.Unmarshal(memory, &hourEventsRemembered); err != nil {
			log.Error("processDayMemories: Error getting short term source from the brain", err)
		}

		// assemble data for the hour
		emailBody += "Hour: " + key + "\n\n"
		for _, memoryRemembered := range hourEventsRemembered {
			emailBody += memoryRemembered.EventId + ":" + memoryRemembered.EventStatus
			for _, memorySequence := range memoryRemembered.SequentialEvents {
				emailBody += " -> " + memorySequence.EventId + ":" + memorySequence.EventStatus
			}
			emailBody += "\n"
		}
		emailBody += "\n\n"

		// delete the hour of data from the database so that it isn't processed again
		if err := database.MartianData.DeleteMemoryHourFromDay(key); err != nil {
			log.Error("processDayMemories: Failed to delete", err)
		}
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
			log.Error(err)
		}
	}
	go b.processDayMemories()
}
