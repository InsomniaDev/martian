package harmony

import (
	"github.com/gorilla/websocket"
)

// Device type
type Device struct {
	Name            string     `json:"name"`
	ActivityID      string     `json:"activityId"`
	Activities      []Activity `json:"activities"`
	AreaName        string     `json:"areaName"`
	Actions         string
	Connection      *websocket.Conn
	IPAddress       string `json:"ipAddress"`
	CurrentActivity string
	ActiveRemoteID  int
	HostName        string
}

// HTTPMessage is the message returned from the HTTP call
type HTTPMessage struct {
	ID   int    `json:"id"`
	Msg  string `json:"msg"`
	Data struct {
		DiscoveryServerCF string `json:"discoveryServerCF"`
		Email             string `json:"email"`
		Username          string `json:"username"`
		ActiveRemoteID    int    `json:"activeRemoteId"`
		DiscoveryServer   string `json:"discoveryServer"`
		SE                bool   `json:"se"`
		SUSChannel        string `json:"susChannel"`
		Mode              int    `json:"mode"`
		AccountID         string `json:"accountId"`
	} `json:"data"`
	Code string `json:"code"`
}

// Activity is the harmony activity info
type Activity struct {
	ActivityID string            `json:"activityId"`
	Name       string            `json:"name"`
	Actions    []ActivityActions `json:"actions"`
}

// ActivityActions are the actions for the harmony activity
type ActivityActions struct {
	Label  string                `json:"label"`
	Action ActivityActionCommand `json:"action"`
}

// ActivityActionCommand are the commands for the activity actions
type ActivityActionCommand struct {
	Command  string `json:"command"`
	Type     string `json:"type"`
	DeviceID string `json:"deviceId"`
}

// RecvMessage is to retrieve the ID for an incoming websocket broadcast
type RecvMessage struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

// RecvResult is the message received in the websocket
type RecvResult struct {
	Cmd  string `json:"cmd"`
	Code int    `json:"code"`
	ID   int    `json:"id"`
	Msg  string `json:"msg"`
	Data struct {
		Result string `json:"result"`
	} `json:"data"`
}

// RecvConfig is the returned configuration from the harmony websocket
type RecvConfig struct {
	Cmd  string `json:"cmd"`
	Code int64  `json:"code"`
	ID   int64  `json:"id"`
	Msg  string `json:"msg"`
	Data struct {
		Activity []struct {
			ID           string `json:"id"`
			BaseImageURI string `json:"baseImageUri"`
			ControlGroup []struct {
				Name     string `json:"name"`
				Function []struct {
					Action string `json:"action"`
					Name   string `json:"name"`
					Label  string `json:"label"`
				} `json:"function"`
			} `json:"controlGroup"`
			SuggestedDisplay string `json:"suggestedDisplay"`
			ActivityOrder    int64  `json:"activityOrder"`
			Label            string `json:"label"`
		} `json:"activity"`
	} `json:"data"`
}

// Command is used for certain commands such as VolumeUp
type Command struct {
	Command  string `json:"command"`
	Type     string `json:"type"`
	DeviceID string `json:"deviceId"`
}

// StateDigestNotify is the structure for notifications from hub
type StateDigestNotify struct {
	Type string                `json:"type"`
	Data StateDigestNotifyData `json:"data"`
}

// StateDigestNotifyData is the data returned from the StateDigestNotify notification
type StateDigestNotifyData struct {
	SleepTimerID        int           `json:"sleepTimerId"`
	RunningZoneList     []interface{} `json:"runningZoneList"`
	DiscoveryServerCF   string        `json:"discoveryServerCF"`
	ActivityID          string        `json:"activityId"`
	SyncStatus          int           `json:"syncStatus"`
	DiscoveryServer     string        `json:"discoveryServer"`
	StateVersion        int           `json:"stateVersion"`
	TzOffset            string        `json:"tzOffset"`
	Mode                int           `json:"mode"`
	Sequence            bool          `json:"sequence"`
	HubSwVersion        string        `json:"hubSwVersion"`
	DeviceSetupState    []interface{} `json:"deviceSetupState"`
	Tzoffset            string        `json:"tzoffset"`
	IsSetupComplete     bool          `json:"isSetupComplete"`
	RunningActivityList string        `json:"runningActivityList"`
	ContentVersion      int           `json:"contentVersion"`
	ConfigVersion       int           `json:"configVersion"`
	Updates             []interface{} `json:"updates"`
	ActivityStatus      int           `json:"activityStatus"`
	WifiStatus          int           `json:"wifiStatus"`
	Tz                  string        `json:"tz"`
	Time                int           `json:"time"`
	IPIRConversionDate  string        `json:"IPIRConversionDate"`
	HubUpdate           bool          `json:"hubUpdate"`
	ActivitySetupState  bool          `json:"activitySetupState"`
	AccountID           string        `json:"accountId"`
}
