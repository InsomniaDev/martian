package lutron

import (
	"bufio"
	"net"

	config "github.com/insomniadev/martian/integrations/config"
	"github.com/insomniadev/martian/modules/pubsub"
)

// LDevice type
type LDevice struct {
	Name     string  `json:"name"`
	ID       int     `json:"id"`
	AreaName string  `json:"areaName"`
	Type     string  `json:"type"`
	Value    float64 `json:"value"`
	State    string  `json:"state"`
}

type MsgType int
type Command string

const (
	Get MsgType = iota
	Set
	Watch
	Response
)

const (
	Output  Command = "OUTPUT"
	Device  Command = "DEVICE"
	Group   Command = "GROUP"
	Unknown Command = "UNKNOWN"
)

type Lutron struct {
	Config    config.LutronConfig
	conn      net.Conn
	reader    *bufio.Reader
	done      chan bool
	Inventory []*LDevice
	broker    *pubsub.PubSub
	Changed   bool
}

type LutronMsg struct {
	// the lutron component number
	Id    int
	Name  string
	Value float64
	// duration in seconds for a set action
	// TODO parse > 60 seconds into string "M:SS"
	Fade float64
	// the action to take with the command, Get, Set, Watch, Default: Get
	Type MsgType
	// the integration command type - Output, Device
	Cmd Command
	// usually the button press
	Action int
	// in Unix nanos format
	Timestamp int64
	// TODO
	// Action Number - default to 1 for now
}

// CasetaIntegrationFile is the structure for the integration file for Lutron
type CasetaIntegrationFile struct {
	LIPIDList struct {
		Devices []struct {
			ID      int    `json:"ID"`
			Name    string `json:"Name"`
			Buttons []struct {
				Name   string `json:"Name"`
				Number int    `json:"Number"`
			} `json:"Buttons"`
		} `json:"Devices"`
		Zones []struct {
			ID   int    `json:"ID"`
			Name string `json:"Name"`
			Area struct {
				Name string `json:"Name"`
			} `json:"Area"`
			Type string `json:"Type"`
		} `json:"Zones"`
	} `json:"LIPIdList"`
}

type ResponseWatcher struct {
	matchMsg  *LutronMsg
	incomming chan interface{}
	Responses chan *LutronMsg
	stop      chan bool
}
