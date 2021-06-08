package life360

var (
	baseURL                          = "https://api.life360.com/v3/"
	tokenURL                         = baseURL + "oauth2/token.json"
	circlesURL                       = baseURL + "circles.json"
	circleURL                        = "circles/"
	membersURL                       = "/members"
	placesURL                        = "/places"
	Life360AuthenticationToken       = "Life360AuthenticationToken"
	Life360AuthenticationBearerToken = "Life360AuthenticationBearerToken"
	Life360Username                  = "Life360Username"
	Life360Password                  = "Life360Password"
)

// Life360 struct for all work on the Life360 API
type Life360 struct {
	Username           string
	Password           string
	AuthorizationToken string
	AccessToken        string
	Circles            []Circle
	Members            []Member
	Places             []Place
}

// AuthenticationPost format for the authentication post
type AuthenticationPost struct {
	GrantType string `json:"grant_type"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

// AuthResponse response from authentication with bearer token
type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

// Circles is the individual response for the CircleResponse
type Circles struct {
	Circles []Circle `json:"circles"`
}

type Circle struct {
	ID                  string         `json:"id"`
	Name                string         `json:"name"`
	Color               string         `json:"color"`
	Type                string         `json:"type"`
	CreatedAt           string         `json:"createdAt"`
	MemberCount         string         `json:"memberCount"`
	UnreadMessages      string         `json:"unreadMessages"`
	UnreadNotifications string         `json:"unreadNotifications"`
	Features            CircleFeatures `json:"features"`
}

type CircleFeatures struct {
	OwnerID             interface{} `json:"ownerId"`
	SkuID               interface{} `json:"skuId"`
	Premium             string      `json:"premium"`
	LocationUpdatesLeft int64       `json:"locationUpdatesLeft"`
	PriceMonth          string      `json:"priceMonth"`
	PriceYear           string      `json:"priceYear"`
	SkuTier             interface{} `json:"skuTier"`
}

type Life360Members struct {
	Members []Member `json:"members"`
}

type Member struct {
	Features       Features        `json:"features"`
	Issues         Issues          `json:"issues"`
	Location       Location        `json:"location"`
	Communications []Communication `json:"communications"`
	Medical        interface{}     `json:"medical"`
	Relation       interface{}     `json:"relation"`
	CreatedAt      string          `json:"createdAt"`
	Activity       interface{}     `json:"activity"`
	ID             string          `json:"id"`
	FirstName      string          `json:"firstName"`
	LastName       string          `json:"lastName"`
	IsAdmin        string          `json:"isAdmin"`
	Avatar         *string         `json:"avatar"`
	PinNumber      interface{}     `json:"pinNumber"`
	LoginEmail     string          `json:"loginEmail"`
	LoginPhone     string          `json:"loginPhone"`
}

type Communication struct {
	Channel string  `json:"channel"`
	Value   string  `json:"value"`
	Type    *string `json:"type"`
}

type Features struct {
	Device                string      `json:"device"`
	Smartphone            string      `json:"smartphone"`
	NonSmartphoneLocating string      `json:"nonSmartphoneLocating"`
	Geofencing            string      `json:"geofencing"`
	ShareLocation         string      `json:"shareLocation"`
	ShareOffTimestamp     interface{} `json:"shareOffTimestamp"`
	Disconnected          string      `json:"disconnected"`
	PendingInvite         string      `json:"pendingInvite"`
	MapDisplay            string      `json:"mapDisplay"`
}

type Issues struct {
	Disconnected    string      `json:"disconnected"`
	Type            interface{} `json:"type"`
	Status          interface{} `json:"status"`
	Title           interface{} `json:"title"`
	Dialog          interface{} `json:"dialog"`
	Action          interface{} `json:"action"`
	Troubleshooting string      `json:"troubleshooting"`
}

type Location struct {
	Latitude       string      `json:"latitude"`
	Longitude      string      `json:"longitude"`
	Accuracy       string      `json:"accuracy"`
	StartTimestamp int64       `json:"startTimestamp"`
	EndTimestamp   string      `json:"endTimestamp"`
	Since          int64       `json:"since"`
	Timestamp      string      `json:"timestamp"`
	Name           string      `json:"name"`
	PlaceType      interface{} `json:"placeType"`
	Source         string      `json:"source"`
	SourceID       string      `json:"sourceId"`
	Address1       string      `json:"address1"`
	Address2       string      `json:"address2"`
	ShortAddress   string      `json:"shortAddress"`
	InTransit      string      `json:"inTransit"`
	TripID         interface{} `json:"tripId"`
	DriveSDKStatus interface{} `json:"driveSDKStatus"`
	Battery        string      `json:"battery"`
	Charge         string      `json:"charge"`
	WifiState      string      `json:"wifiState"`
	Speed          float64     `json:"speed"`
	IsDriving      string      `json:"isDriving"`
	UserActivity   interface{} `json:"userActivity"`
}

type Places struct {
	Places []Place `json:"places"`
}

type Place struct {
	ID        string `json:"id"`
	OwnerID   string `json:"ownerId"`
	CircleID  string `json:"circleId"`
	Name      string `json:"name"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Radius    string `json:"radius"`
	Type      int64  `json:"type"`
	TypeLabel string `json:"typeLabel"`
}

// _PROTOCOL = 'https://'
// _BASE_URL = '{}api.life360.com/v3/'.format(_PROTOCOL)
// _TOKEN_URL = _BASE_URL + 'oauth2/token.json'
// _CIRCLES_URL = _BASE_URL + 'circles.json'
// _CIRCLE_URL = _BASE_URL + 'circles/{}'
// _CIRCLE_MEMBERS_URL = _CIRCLE_URL + '/members'
// _CIRCLE_PLACES_URL = _CIRCLE_URL + '/places'
