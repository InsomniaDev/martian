package life360

import (
	log "github.com/sirupsen/logrus"

	"github.com/insomniadev/martian/integrations/config"
)

// InsertAuth inserts the username and password into the database
func (life *Life360) InsertAuth(username string, password string) {
	log.Info("Not implemented")

	// bolt.UpdateAccount(Life360AuthenticationToken, "cFJFcXVnYWJSZXRyZTRFc3RldGhlcnVmcmVQdW1hbUV4dWNyRUh1YzptM2ZydXBSZXRSZXN3ZXJFQ2hBUHJFOTZxYWtFZHI0Vg==")
}

// InsertBearerToken inserts the authenticated token to the database
func (life *Life360) InsertBearerToken(token string) {
	log.Info("Not implemented")
}

// RetrieveAuth from the bolt key value store
func (life *Life360) RetrieveAuth() {

	life.Username, life.Password, life.AuthorizationToken = config.LoadLife360()
}

// Initialize the Life360 application
func (life *Life360) Initialize() {
	life.RetrieveAuth()
}
