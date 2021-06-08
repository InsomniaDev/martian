package sleepiq

import (
	config "github.com/insomniadev/martian/integrations/config"
)

// InsertAuth inserts the username and password into bolt
func (si *SleepIQ) InsertAuth() {
	username, password := config.LoadSleepIq()
	si.Username = username
	si.Username = password
}

// RetrieveAuth retrieves auth from the bolt key value store
func (si *SleepIQ) RetrieveAuth() {
	username, password := config.LoadSleepIq()
	si.Username = username
	si.Username = password
}

// Initialize sets up all of the properties for sleep IQ
func (si *SleepIQ) Initialize() {
	si.RetrieveAuth()
}
