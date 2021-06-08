package philipshue

var (
	// PhilipsHueURL const name for the url
	PhilipsHueURL = "PHILIPS_HUE_URL"
	// PhilipsHueSecret const name for the secret
	PhilipsHueSecret = "PHILIPS_HUE_SECRET"
	// PhilipsHueLights const name for all of the lights
	PhilipsHueLights = "PHILIPS_HUE_LIGHTS"
)

type errorsInResponse struct {
	TypeError   int    `json:"type"`
	Description string `json:"description"`
}

type errorResponse struct {
	Err errorsInResponse `json:"error"`
}

// Light is struct for philips light information
type Light struct {
	UUID             string
	State            *State `json:"state,omitempty"`
	Type             string `json:"type,omitempty"`
	Name             string `json:"name,omitempty"`
	ModelID          string `json:"modelid,omitempty"`
	ManufacturerName string `json:"manufacturername,omitempty"`
	UniqueID         string `json:"uniqueid,omitempty"`
	SwVersion        string `json:"swversion,omitempty"`
	SwConfigID       string `json:"swconfigid,omitempty"`
	ProductID        string `json:"productid,omitempty"`
	ID               int    `json:"ID"`
}

// State defines the attributes and properties of a light
type State struct {
	On             bool      `json:"on"`
	Bri            uint8     `json:"bri,omitempty"`
	Hue            uint16    `json:"hue,omitempty"`
	Sat            uint8     `json:"sat,omitempty"`
	Xy             []float32 `json:"xy,omitempty"`
	Ct             uint16    `json:"ct,omitempty"`
	Alert          string    `json:"alert,omitempty"`
	Effect         string    `json:"effect,omitempty"`
	TransitionTime uint16    `json:"transitiontime,omitempty"`
	BriInc         int       `json:"bri_inc,omitempty"`
	SatInc         int       `json:"sat_inc,omitempty"`
	HueInc         int       `json:"hue_inc,omitempty"`
	CtInc          int       `json:"ct_inc,omitempty"`
	XyInc          int       `json:"xy_inc,omitempty"`
	ColorMode      string    `json:"colormode,omitempty"`
	Reachable      bool      `json:"reachable,omitempty"`
	Scene          string    `json:"scene,omitempty"`
}

type deviceType struct {
	DeviceType string `json:"devicetype"`
}

type username struct {
	Username string `json:"username"`
}

type successResponse struct {
	Success username `json:"success"`
}

type state struct {
	On bool `json:"on"`
}

// PhilipsHue is the struct for information of the lights
type PhilipsHue struct {
	URL    string
	Secret string
	Lights []Light
}
