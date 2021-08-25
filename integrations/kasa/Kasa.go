package kasa

// https://github.com/python-kasa/python-kasa

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/insomniadev/martian/database"
	"github.com/insomniadev/martian/modules/redispub"
)

// Init initializes the instance of kasa for devices on the network
func (d *Devices) Init(configuration string) {
	json.Unmarshal([]byte(configuration), &d)
	// d.Plugs = RetrieveKasaNodes()
	// devices := config.LoadKasa()

	// for _, ipAdd := range devices {
	// 	found := false
	// 	for _, dev := range d.Plugs {
	// 		if ipAdd == dev.IPAddress {
	// 			found = true
	// 			dev.PowerState()
	// 		}
	// 	}
	// 	if !found {
	// 		plug := NewPlug(ipAdd)
	// 		plug.PowerOff()
	// 		plug.Name = plug.PlugInfo.Alias
	// 		InsertKasaGraph(ipAdd, plug.Name)
	// 		d.Plugs = append(d.Plugs, plug)
	// 	}
	// }

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
	var db database.Database
	err := db.PutIntegrationValue("kasa", d)
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
			redispub.Service.Publish("subscriptions", h)
		}
	}
}

// UpdateArea assigns the plug to an area
func (h *KasaDevice) UpdateArea(areaName string) {
	UpdateAreaForKasaDevice(h.IPAddress, areaName)
	h.AreaName = areaName
}

// PowerOff turns the plug off
func (h *KasaDevice) PowerOff() error {
	data, err := h.do(PowerOffCommand, "set_relay_state")
	fmt.Println(data)
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
func (h *KasaDevice) PowerOn() error {
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
func (h *KasaDevice) PowerState() (PowerState, error) {
	state, err := h.Info()
	if err != nil {
		return PowerUnknown, fmt.Errorf("could not determine if plug was turned on: %s", err)
	}

	return state.RelayState, nil
}

func (d *Devices) Discover() {
	// ifaces, err := net.Interfaces()
	// if err != nil {
	// 	fmt.Print(fmt.Errorf("localAddresses: %+v", err.Error()))
	// 	return
	// }
	// found := false
	// for _, i := range ifaces {
	// 	addrs, err := i.Addrs()
	// 	if err != nil {
	// 		fmt.Print(fmt.Errorf("localAddresses: %+v", err.Error()))
	// 		continue
	// 	}
	// 	for _, a := range addrs {
	// 		switch v := a.(type) {
	// 		case *net.IPNet:
	addr := "192.168.1.1/24"
	// if strings.Contains(addr, "192.168") || strings.Contains(addr, "10.10") {
	ip, ipnet, err := net.ParseCIDR(addr)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		plug := NewPlug(ip.String())
		wg.Add(1)
		go func() {
			defer wg.Done()
			info, _ := plug.Info()
			// d.Plugs = append(d.Plugs, plug)
			if info != nil {
				plug.Name = plug.PlugInfo.Alias
				switch plug.PlugInfo.Model {
				case "HS105(US)":
					plug.Type = "plug"
				case "HS200(US)":
					plug.Type = "light"
				}
				alreadyUsedPlug := false
				for i := range d.Devices {
					if d.Devices[i].IPAddress == plug.IPAddress {
						d.Devices[i] = plug
						alreadyUsedPlug = true
					}
				}
				if !alreadyUsedPlug {
					d.Devices = append(d.Devices, plug)
				}
			}
		}()
	}
	wg.Wait()
	var db database.Database
	db.PutIntegrationValue("kasa", d)
	// 	}
	// case *net.IPAddr:
	// 	fmt.Printf("%v : %s (%s)\n", i.Name, v, v.IP.DefaultMask())
	// }
	// 	if found {
	// 		break
	// 	}
	// }
	// if found {
	// 	break
	// }
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
	var db database.Database
	err := db.PutIntegrationValue("kasa", h)
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

	// Save in the database
	var db database.Database
	err := db.PutIntegrationValue("hass", k)
	if err != nil {
		log.Println(err)
	}

	// Let's repopulate with the correct device state
	if removeEdit {
		k.Discover()
	}

	return nil
}
