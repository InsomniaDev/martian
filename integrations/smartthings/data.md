package smartthings

import (
	"encoding/json"
	"fmt"
	bolt "homesmartie/models"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func (st *SmartThings) SetAuthentication(clientid string, clientsecret string) (err error) {
	bolt.UpdateAccount(SmartThingsClientID, clientid)
	bolt.UpdateAccount(SmartThingsClientSecret, clientsecret)
	// bolt.UpdateAccount(SmartThingsOauthToken, oauth) , oauth string
	return
}

// SaveToken saves the token to the database
func (st *SmartThings) SaveToken(token *oauth2.Token) error {
	insertToken, err := json.Marshal(token)
	if err != nil {
		return err
	}
	bolt.UpdateAccount(SmartThingsOauthToken, string(insertToken))
	return nil
}

func (st *SmartThings) LoadToken() (*oauth2.Token, error) {
	token, err := bolt.ReadAccount(SmartThingsOauthToken)
	if err != nil {
		return nil, err
	}
	oauthToken := &oauth2.Token{}
	if err := json.Unmarshal([]byte(token), &oauthToken); err != nil {
		return nil, err
	}
	return oauthToken, nil
}

func (st *SmartThings) RetrieveAuth() {
	clientId, err := bolt.ReadAccount(SmartThingsClientID)
	if err != nil {
		log.Println(err)
	}
	st.ClientID = clientId
	clientSecret, err := bolt.ReadAccount(SmartThingsClientSecret)
	if err != nil {
		log.Println(err)
	}
	st.ClientSecret = clientSecret
	oauthToken, err := bolt.ReadAccount(SmartThingsOauthToken)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal([]byte(oauthToken), &st.Token)

	config := st.NewOAuthConfig(st.ClientID, st.ClientSecret)
	if !st.Token.Valid() {
		st.Token, err = st.GetToken(config)
	}

	ctx := context.Background()
	st.Client = config.Client(ctx, st.Token)

	st.Endpoint, err = st.GetEndPointsURI()
	if err != nil {
		return
	}
	return
}

// Initialize the smartthings struct
func (st *SmartThings) Initialize() {
	st.RetrieveAuth()
	// id, err := db.RetrieveValue(SmartThingsClientID)
	// if err != nil {
	// 	return
	// }
	// st.ClientID = id.Value

	// secret, err := db.RetrieveValue(SmartThingsClientSecret)
	// if err != nil {
	// 	return
	// }
	// st.ClientSecret = secret.Value

	// oauthToken, err := db.RetrieveValue(SmartThingsOauthToken)
	// if err != nil {
	// 	return
	// }
	// json.Unmarshal([]byte(oauthToken.Value), &st.Token)

	// if st.ClientID == "" || st.ClientSecret == "" {
	// 	return fmt.Errorf("Need to have client id and secret")
	// }

	// config := st.NewOAuthConfig(st.ClientID, st.ClientSecret)

	// if !st.Token.Valid() {
	// 	st.Token, err = st.GetToken(config)
	// }

	// if err != nil {
	// 	return
	// }

}
