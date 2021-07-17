package kasa

// https://github.com/python-kasa/python-kasa

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/insomniadev/martian/database"
	"github.com/insomniadev/martian/modules/redispub"
)

// Init initializes the instance of kasa for devices on the network
func (d *Devices) Init(configuration string) {
	var devices []Plug
	json.Unmarshal([]byte(configuration), &devices)
	d.Plugs = devices
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
	for i := range d.Plugs {
		go d.Plugs[i].WatchForChanges()
	}
}

func (d *Devices) ChangeAreaForKasaDevice(ipAddress, area string) error {
	for i := range d.Plugs {
		if d.Plugs[i].IPAddress == ipAddress {
			d.Plugs[i].AreaName = area
		}
	}
	var db database.Database
	err := db.PutIntegrationValue("kasa", d.Plugs)
	if err != nil {
		return err
	}
	return nil
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
				for i := range d.Plugs {
					if d.Plugs[i].IPAddress == plug.IPAddress {
						d.Plugs[i] = plug
						alreadyUsedPlug = true
					}
				}
				if !alreadyUsedPlug {
					d.Plugs = append(d.Plugs, plug)
				}
				fmt.Println(ip.String(), info)
			}
		}()
	}
	wg.Wait()
	var db database.Database
	db.PutIntegrationValue("kasa", d.Plugs)
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
