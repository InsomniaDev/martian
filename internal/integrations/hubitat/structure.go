package hubitat

type Hubitat struct {
	AccessKey      string `mapstructure:"accessKey"`
	HubitatUrl     string `mapstructure:"url"`
	PreviousStates HubitatDeviceAllResponse
	CurrentStates  HubitatDeviceAllResponse
}

// HubitatDeviceAllResponse is the response returned for getting the current status of all devices
type HubitatDeviceAllResponse []HubitatDeviceAllResponseElement

type HubitatDeviceAllResponseElement struct {
	Name         string      `json:"name"`
	Label        string      `json:"label"`
	Type         string      `json:"type"`
	ID           string      `json:"id"`
	Date         *string     `json:"date"`
	Model        interface{} `json:"model"`
	Manufacturer interface{} `json:"manufacturer"`
	Capabilities []string    `json:"capabilities"`
	Attributes   Attributes  `json:"attributes"`
	Commands     []Command   `json:"commands"`
}

type Attributes struct {
	Presence         interface{} `json:"presence"`
	DataType         DataType    `json:"dataType"`
	Values           []string    `json:"values"`
	NotificationText interface{} `json:"notificationText"`
	Battery          string     `json:"battery,omitempty"`
	ThreeAxis        string     `json:"threeAxis,omitempty"`
	Acceleration     string     `json:"acceleration"`
	Tamper           interface{} `json:"tamper"`
	Temperature      string     `json:"temperature,omitempty"`
	Illuminance      string     `json:"illuminance,omitempty"`
	Motion           string     `json:"motion,omitempty"`
	PendingChanges   string     `json:"pendingChanges,omitempty"`
	Held             interface{} `json:"held"`
	Pushed           interface{} `json:"pushed"`
	Released         interface{} `json:"released"`
	DoubleTapped     interface{} `json:"doubleTapped"`
	NumberOfButtons  string     `json:"numberOfButtons,omitempty"`
	Switch           string     `json:"switch,omitempty"`
	Level            string     `json:"level,omitempty"`
	Status           string     `json:"status,omitempty"`
	TrackData        string     `json:"trackData,omitempty"`
	Volume           string     `json:"volume,omitempty"`
	Mute             string     `json:"mute,omitempty"`
	TrackDescription string     `json:"trackDescription,omitempty"`
	PowerSource      string     `json:"powerSource,omitempty"`
}

type Command struct {
	Command string `json:"command"`
}

type DataType string

const (
	Enum   DataType = "ENUM"
	String DataType = "STRING"
)

type Switch string

const (
	Off Switch = "off"
	On  Switch = "on"
)

// HubitatUpdate is what gets sent from the Hubitat device to this service
type HubitatUpdate struct {
    Content Content `json:"content"`
}

type Content struct {
    Name            string      `json:"name"`           
    Value           string      `json:"value"`          
    DisplayName     string      `json:"displayName"`    
    DeviceID        string      `json:"deviceId"`       
    DescriptionText string      `json:"descriptionText"`
    Unit            interface{} `json:"unit"`           
    Type            string      `json:"type"`           
    Data            interface{} `json:"data"`           
}