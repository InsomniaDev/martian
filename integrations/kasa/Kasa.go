package kasa

import (
	"fmt"
	"time"

	"github.com/insomniadev/martian/integrations/config"
	"github.com/insomniadev/martian/modules/redispub"
)

// Init initializes the instance of kasa for devices on the network
func (d *Devices) Init() {
	d.Plugs = RetrieveKasaNodes()
	devices := config.LoadKasa()

	for _, ipAdd := range devices {
		found := false
		for _, dev := range d.Plugs {
			if ipAdd == dev.IPAddress {
				found = true
				dev.PowerState()
			}
		}
		if !found {
			plug := NewPlug(ipAdd)
			plug.PowerOff()
			plug.Name = plug.PlugInfo.Alias
			InsertKasaGraph(ipAdd, plug.Name)
			d.Plugs = append(d.Plugs, plug)
		}
	}

	for i := range d.Plugs {
		go d.Plugs[i].WatchForChanges()
	}
}

// WatchForChanges will constantly check to assert plug state
func (h *Plug) WatchForChanges() {
	for {
		time.Sleep(1 * time.Second)
		previousState := h.PlugInfo.On
		h.PowerState()
		if previousState != h.PlugInfo.On {
			redispub.Service.Publish("subscriptions", h)
		}
	}
}

// UpdateArea assigns the plug to an area
func (h *Plug) UpdateArea(areaName string) {
	UpdateAreaForKasaDevice(h.IPAddress, areaName)
	h.AreaName = areaName
}

// PowerOff turns the plug off
func (h *Plug) PowerOff() error {
	_, err := h.do(PowerOffCommand, "set_relay_state")
	if err != nil {
		return err
	}

	state, err := h.PowerState()
	if err != nil {
		return err
	}

	if state != PowerOff {
		return fmt.Errorf("power off was requested but device stayed on")
	}
	redispub.Service.Publish("subscriptions", h)
	return nil
}

// PowerOn turns the plug on
func (h *Plug) PowerOn() error {
	data, err := h.do(PowerOnCommand, "set_relay_state")
	fmt.Println(string(data))
	if err != nil {
		return err
	}

	state, err := h.PowerState()
	if err != nil {
		return err
	}

	if state != PowerOn {
		return fmt.Errorf("power on was requested but device stayed off")
	}
	redispub.Service.Publish("subscriptions", h)
	return err
}

// PowerState retrieves the current power state of the plug, PowerUnknown when request failed
func (h *Plug) PowerState() (PowerState, error) {
	state, err := h.Info()
	if err != nil {
		return PowerUnknown, fmt.Errorf("could not determine if plug was turned on: %s", err)
	}

	return state.RelayState, nil
}
