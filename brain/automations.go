package brain

import "log"

// TODO: Add functionality for dusk to dawn timers
// TODO: Add functionality for one thing to have an action on another (motion sensor turns on light)
// TODO: Add integration for hubitat

type AutomationTypes string

const (
	Time   AutomationTypes = "time"
	Action AutomationTypes = "action"
)

type DeviceTypes string

const (
	Hubitat DeviceTypes = "hubitat"
)

type Automations struct {
	AutomationConfigs []Automation `yaml:"automations"`
}

type Automation struct {
	Type     AutomationTypes `yaml:"yaml"`
	Action   ActionDevice    `yaml:"actionDevice,omitempty"`
	Reaction ReactionDevice  `yaml:"reactionDevice"`
}

type ActionDevice struct {
	Id    string      `yaml:"id,omitempty"`
	Type  DeviceTypes `yaml:"type,omitempty"`
	State string      `yaml:"state,omitempty"`
}

type ReactionDevice struct {
	Id    string      `yaml:"id"`
	Type  DeviceTypes `yaml:"type"`
	State string      `yaml:"state"`
}

// decideAutomation will determine which event was called and how the application should respond
func (a *Automation) decideAutomation(triggering AutomationTypes, deviceState string) {
	switch triggering {
	case Time:
		a.timeEvent()
	case Action:
		a.actionEvent(deviceState)
	}
}

// timeEvent will fire when there is a time based application process
func (a *Automation) timeEvent() {
	log.Println("Need to change based on the time")
	a.callReactionDevice()
}

// actionEvent will fire when there is a change to the state of a connected device that matches with an actionDeviceId
func (a *Automation) actionEvent(deviceState string) {
	if deviceState == "" {
		log.Println("DeviceState is blank, unable to process actionEvent")
		return
	}

	if deviceState == a.Action.State {
		log.Println("DEBUG: Action matches state")
		a.callReactionDevice()
	}
}

func (a *Automation) callReactionDevice() {
	switch a.Reaction.Type {
	case Hubitat:
		// TODO: Have it call the hubitat integration and change it from here
	}
}
