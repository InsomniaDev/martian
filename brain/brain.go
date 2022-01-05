package brain

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/insomniadev/martian/modules/pubsub"
)

type Brain struct {
	memoryEvent []event
}

type event struct {
	eventTime           time.Time
	eventTimeExpiration time.Time
	eventId             string
	eventStatus         string
}

var (
	Brainiac       *Brain
	timeDifference = 2 * time.Minute
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
}

func (b *Brain) brainWave(id, status string) {
	// TODO: Add more logic here, just append for now
	b.memoryEvent = append(b.memoryEvent, event{eventId: id, eventStatus: status, eventTime: time.Now(), eventTimeExpiration: time.Now().Add(timeDifference)})
	fmt.Println("I remember that", id, status)
}
