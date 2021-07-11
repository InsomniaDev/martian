package kasa

import (
	"encoding/json"
	"time"
)

// PowerState represents a state for the plug relay
type PowerState int

type Devices struct {
	Plugs []Plug
}

const (
	// PowerOnCommand is the command sent to turn the relay on
	PowerOnCommand = `{"system":{"set_relay_state":{"state":1}}}`

	// PowerOffCommand is the command sent to turn the relay off
	PowerOffCommand = `{"system":{"set_relay_state":{"state":0}}}`

	// InfoCommand is the command sent to request system into
	InfoCommand = `{"system":{"get_sysinfo":{}}}`

	// EnergyCommand is the command to retrieve energy info
	EnergyCommand = `{"emeter":{"get_realtime":{}}}`

	// RebootCommand is the command to reboot the plug
	RebootCommand = `{"system":{"reboot":{"delay":1}}}`

	// PowerUnknown represents an unknown power state
	PowerUnknown PowerState = -1

	// PowerOff represents the off state off the plug
	PowerOff PowerState = 0

	// PowerOn represents the off state on the plug
	PowerOn PowerState = 1

	port          = 9999
	cryptKey      = byte(0xAB)
	connTimeout   = 400 * time.Millisecond
	writeDeadline = 2
	readDeadline  = 2
)

// Plug represents a management interface for a plug
type Plug struct {
	ID        string
	IPAddress string
	PlugInfo  Info
	Name      string
	AreaName  string
	Type      string
}

// NewPlug creates a new management interface for the TP Link HS1xx plug
func NewPlug(ip string) Plug {
	return Plug{
		ID:        ip,
		IPAddress: ip,
	}
}

// Info is the system
type Info struct {
	SoftwareVersion string     `json:"sw_ver"`
	HardwareVersion string     `json:"hw_ver"`
	Type            string     `json:"type"`
	Model           string     `json:"model"`
	MAC             string     `json:"mac"`
	DeviceName      string     `json:"dev_name"`
	Alias           string     `json:"alias"`
	RelayState      PowerState `json:"relay_state"`
	OnTimeSeconds   int        `json:"on_time"`
	OnTime          string     `json:"on_time_string"`
	ActiveMode      string     `json:"active_mode"`
	Features        string     `json:"feature"`
	Updating        int        `json:"updating"`
	SignalStrength  int        `json:"rssi"`
	LEDOff          int        `json:"led_off"`
	Lon             int        `json:"longitude_i"`
	Lat             int        `json:"latitude_i"`
	HardwareID      string     `json:"hwId"`
	FirmwareID      string     `json:"fwId"`
	DeviceID        string     `json:"deviceId"`
	OEMID           string     `json:"oemId"`
	NTCState        int        `json:"ntc_state"`
	Error           error

	On        bool            `json:"power_on"`
	Off       bool            `json:"power_off"`
	Address   string          `json:"address"`
	RawStatus json.RawMessage `json:"-"`
}
