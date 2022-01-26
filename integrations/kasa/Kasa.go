package kasa

// https://github.com/python-kasa/python-kasa

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/insomniadev/martian/database"
	"github.com/insomniadev/martian/modules/pubsub"
)

// Init initializes the instance of kasa for devices on the network
func (d *Devices) Init(configuration string) {
	json.Unmarshal([]byte(configuration), &d)

	d.Discover()

	// Check every second for a change in the connected kasa devices and then update on that change
	for i := range d.Devices {
		go d.Devices[i].WatchForChanges()
	}
}

func (d *Devices) ChangeAreaForKasaDevice(ipAddress, area string) error {
	for i := range d.Devices {
		if d.Devices[i].IPAddress == ipAddress {
			d.Devices[i].AreaName = area
		}
	}
	
	err := database.MartianData.PutIntegrationValue("kasa", d)
	if err != nil {
		return err
	}
	return nil
}

// WatchForChanges will constantly check to assert plug state
func (h *KasaDevice) WatchForChanges() {
	for {
		time.Sleep(1 * time.Second)
		previousState := h.PlugInfo.On
		h.PowerState()
		if previousState != h.PlugInfo.On {
			pubsub.Service.Publish("subscriptions", "kasa")
			eventMessage := "kasa;;" + h.PlugInfo.Alias + ";;" + strconv.FormatBool(h.PlugInfo.On)
			pubsub.Service.Publish("brain", eventMessage)
		}
	}
}

// PowerOff turns the plug off
func (h *KasaDevice) PowerOff() error {
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
	pubsub.Service.Publish("subscriptions", "kasa")
	eventMessage := "kasa;;" + h.PlugInfo.Alias + ";;" + strconv.FormatBool(h.PlugInfo.On)
	pubsub.Service.Publish("brain", eventMessage)
	return nil
}

// PowerOn turns the plug on
func (h *KasaDevice) PowerOn() error {
	_, err := h.do(PowerOnCommand, "set_relay_state")
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
	pubsub.Service.Publish("subscriptions", "kasa")
	eventMessage := "kasa;;" + h.PlugInfo.Alias + ";;" + strconv.FormatBool(h.PlugInfo.On)
	pubsub.Service.Publish("brain", eventMessage)
	return err
}

// PowerState retrieves the current power state of the plug, PowerUnknown when request failed
func (h *KasaDevice) PowerState() (PowerState, error) {
	state, err := h.Info()
	if err != nil {
		return PowerUnknown, fmt.Errorf("could not determine if plug was turned on: %s", err)
	}

	return state.RelayState, nil
}

// Discover will go through and discover all Kasa devices on the network
func (d *Devices) Discover() {
	d.Devices = []KasaDevice{}

	// Get all ips in the cidr
	ip, ipnet, err := net.ParseCIDR(d.IpAddressCidr)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup

	// Check each ip to see if it responds as a Kasa device
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		plug := NewPlug(ip.String())
		wg.Add(1)
		go func() {
			defer wg.Done()
			info, _ := plug.Info()
			if info != nil {
				plug.Name = plug.PlugInfo.Alias
				switch plug.PlugInfo.Model {
				case "HS105(US)":
					plug.Type = "plug"
				case "HS103(US)":
					plug.Type = "plug"
				case "HS200(US)":
					plug.Type = "light"
				}

				// Check if this device has been edited already and update it accordingly
				for i := range d.EditedDevices {
					if d.EditedDevices[i].IPAddress == plug.IPAddress {
						plug.Name = d.EditedDevices[i].Name
						plug.AreaName = d.EditedDevices[i].AreaName
					}
				}
				d.Devices = append(d.Devices, plug)
			}
		}()
	}
	wg.Wait()

	// Insert into the database again with all devices
	// err = database.MartianData.PutIntegrationValue("kasa", d)
	// if err != nil {
	// 	log.Println(err)
	// }
}
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// UpdateSelectedDevices will go through and update the devices as selected or not selected
func (h *Devices) UpdateSelectedDevices(selectedDevices []string, addDevices bool, automationDevice bool) error {

	// Compare either through the automation or interface selections
	if automationDevice {
		h.AutomatedDevices = checkIfDeviceIsInList(h.Devices, h.AutomatedDevices, selectedDevices, addDevices)
	} else {
		h.InterfaceDevices = checkIfDeviceIsInList(h.Devices, h.InterfaceDevices, selectedDevices, addDevices)
	}
	
	err := database.MartianData.PutIntegrationValue("kasa", h)
	if err != nil {
		log.Println(err)
	}

	return nil
}

// checkIfDeviceIsInList is an internal method to see if the value already exists in the list
func checkIfDeviceIsInList(allDevices []KasaDevice, alreadyChosenDevices []string, selectedDevices []string, addDevices bool) []string {
	var newlySelectedDevices []string
	// Cycle through all of the available devices for HomeAssistant
	for _, availableDevice := range allDevices {

		selectedDeviceExists := false
		// Cycle through all of the already selected devices to see if there is a match
		for _, selectedDeviceIpAddress := range alreadyChosenDevices {

			// If this available device is already selected, then set selectedDeviceExists as true
			if availableDevice.IPAddress == selectedDeviceIpAddress {
				selectedDeviceExists = true
				break // break out of the selectedDevice cycle since there is a match
			}
		}

		availableDeviceIsNowSelected := false
		// Cycle through the selectedDevices parameter to see if the available device has a match
		for _, newlySelectedDevice := range selectedDevices {

			// If the available device matches with the newly selected criteria then set as true
			if availableDevice.IPAddress == newlySelectedDevice {
				availableDeviceIsNowSelected = true
			}
		}

		// IF the device is already selected, and is one of the newly selected, and set to be added
		//			IF addDevices is FALSE, then it will be removed from the selectedDevices slice
		if availableDeviceIsNowSelected && addDevices {
			newlySelectedDevices = append(newlySelectedDevices, availableDevice.IPAddress)
		} else if selectedDeviceExists && !availableDeviceIsNowSelected { // IF it was already selected
			newlySelectedDevices = append(newlySelectedDevices, availableDevice.IPAddress)
		}
	}
	return newlySelectedDevices
}

// EditDeviceConfiguration will go through and update information for the passed in device
func (k *Devices) EditDeviceConfiguration(device KasaDevice, removeEdit bool) error {

	// Cycle through all of the devices, interfaceDevices, and automatedDevices and update
	for i := range k.Devices {
		if k.Devices[i].IPAddress == device.IPAddress {
			k.Devices[i] = device
		}
	}

	// Update device in edited section, if not found then add it to the edited devices
	foundDevice := false
	for i := range k.EditedDevices {
		if k.EditedDevices[i].IPAddress == device.IPAddress {
			k.EditedDevices[i] = device
			foundDevice = true
		}
	}
	if !foundDevice {
		k.EditedDevices = append(k.EditedDevices, device)
	}

	// Save in the database
	err := database.MartianData.PutIntegrationValue("kasa", k)
	if err != nil {
		log.Println(err)
	}

	// Let's repopulate with the correct device state
	if removeEdit {
		k.Discover()
	}

	return nil
}
