package smartthings

import (
	"net/http"

	"golang.org/x/oauth2"
)

const (
	authDone  = "<html><body>Authentication Completed.</body></html>"
	authError = "<html><body>AUthentication error. Please see terminal output for details.</body></html>"

	// Endpoints URL
	endPointsURI = "https://graph.api.smartthings.com/api/smartapps/endpoints"

	// URL paths used for Oauth authentication on localhost
	callbackPath = "/OAuthCallback"
	donePath     = "/OauthDone"
	rootPath     = "/"

	// default local HTTP server port
	defaultPort = 4567
)

var (
	SmartThingsOauthToken = "SMART_THINGS_OAUTH_TOKEN"
	// SmartThingsClientID client ID that is used
	SmartThingsClientID = "SMART_THINGS_CLIENT_ID"
	// SmartThingsClientSecret client secret that is used
	SmartThingsClientSecret = "SMART_THINGS_CLIENT_SECRET"
)

// DeviceList holds the list of devices returned by /devices
type DeviceList struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

// DeviceCommand holds one command a device can accept.
type DeviceCommand struct {
	Command string                 `json:"command"`
	Params  map[string]interface{} `json:"params"`
}

// DeviceInfo holds information about a specific device.
type DeviceInfo struct {
	DeviceList
	Attributes map[string]interface{} `json:"attributes"`
}

// SmartThings is the struct for the smartthings integration
type SmartThings struct {
	ClientID     string
	ClientSecret string
	Token        *oauth2.Token
	Client       *http.Client
	Endpoint     string
}

// Auth contains the SmartThings authentication related data.
type Auth struct {
	port             int
	config           *oauth2.Config
	rchan            chan oauthReturn
	oauthStateString string
}

// oauthReturn contains the values returned by the OAuth callback handler.
type oauthReturn struct {
	token *oauth2.Token
	err   error
}

// endpoints holds the values returned by the SmartThings endpoints URI.
type endpoints struct {
	OauthClient struct {
		ClientID string `json:"clientId"`
	}
	Location struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	URI     string `json:"uri"`
	BaseURL string `json:"base_url"`
	URL     string `json:"url"`
}
