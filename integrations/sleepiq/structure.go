package sleepiq

var (
	// SleepIqURL const name for the url
	SleepIqURL = "https://api.sleepiq.sleepnumber.com/"
	// SleepIqUsername const name for the username
	SleepIqUsername = "SLEEP_IQ_USERNAME"
	// SleepIqPassword const name for the password
	SleepIqPassword = "SLEEP_IQ_PASSWORD"
)

type BedInfo struct {
	IsInBed              bool   `json:"isInBed"`
	AlertDetailedMessage string `json:"alertDetailedMessage"`
	SleepNumber          int    `json:"sleepNumber"`
	AlertID              int    `json:"alertId"`
	LastLink             string `json:"lastLink"`
	Pressure             int    `json:"pressure"`
}

type Status struct {
	Status    string  `json:"status"`
	BedID     string  `json:"bedId"`
	LeftSide  BedInfo `json:"leftSide"`
	RightSide BedInfo `json:"rightSide"`
}

type FamilyStatus struct {
	Beds []Status `json:"beds"`
}

type login struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	UserID            string `json:"userId"`
	Key               string `json:"key"`
	RegistrationState string `json:"registrationState"`
	EdpLoginStatus    string `json:"edpLoginStatus"`
	EdpLoginMessage   string `json:"edpLoginMessage"`
}

// SleepIQ is the struct for the SleepIQ API
type SleepIQ struct {
	Username string
	Password string
	Cookies  string
	Key      string
}
