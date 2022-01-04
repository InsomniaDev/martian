package life360

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/insomniadev/martian/modules/pubsub"
)

func (life *Life360) redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+life.AuthorizationToken)
	return nil
}

func (life *Life360) getURL(url string) (contents []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", "Bearer "+life.AccessToken)
	resp, err := client.Do(req)

	contents, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return
}

// GetCircles for Life360
func (life *Life360) GetCircles() (err error) {
	if life.AccessToken == "" {
		life.Authenticate()
	}
	resp, err := life.getURL(circlesURL)
	if err != nil {
		return
	}
	values := Circles{}
	err = json.Unmarshal(resp, &values)
	if err != nil {
		return
	}
	for _, circle := range values.Circles {
		life.Circles = append(life.Circles, circle)
	}
	return
}

// SyncMemberStatus : Constantly sync the life360 member status
func (life *Life360) SyncMemberStatus() {
	for {
		// Update the life360 status every 10 seconds
		time.Sleep(10 * time.Second)
		life.GetMembers()
	}
}

// GetMembers for Life360
func (life *Life360) GetMembers() (err error) {
	for _, circle := range life.Circles {
		if life.AccessToken == "" {
			life.Authenticate()
		}
		var resp []byte
		resp, err = life.getURL(baseURL + circleURL + circle.ID + membersURL)
		if err != nil {
			return
		}
		values := Life360Members{}
		err = json.Unmarshal(resp, &values)
		if err != nil {
			return
		}

		// Save the previous data for comparison
		previousData := life.Members
		// Clear the array to fill with the new data
		life.Members = life.Members[:0]

		for _, member := range values.Members {
			life.Members = append(life.Members, member)
		}
		// Check if locations have changed and if they have then update the subscription
		for _, member := range life.Members {
			for _, oldmember := range previousData {
				if member.FirstName == oldmember.FirstName && (member.Location.Latitude != oldmember.Location.Latitude || member.Location.Longitude != oldmember.Location.Longitude) {
					pubsub.Service.Publish("subscriptions", "life360")
				}
			}
		}
	}
	return
}

// GetPlaces for Life360
func (life *Life360) GetPlaces() (err error) {
	for _, circle := range life.Circles {
		if life.AccessToken == "" {
			life.Authenticate()
		}
		var resp []byte
		resp, err = life.getURL(baseURL + circleURL + circle.ID + placesURL)
		if err != nil {
			return
		}
		values := Places{}
		err = json.Unmarshal(resp, &values)
		if err != nil {
			return
		}
		life.Places = values.Places
	}
	return
}

// Authenticate and return the token for Life360
func (life *Life360) Authenticate() (err error) {
	life.RetrieveAuth()
	client := &http.Client{
		CheckRedirect: life.redirectPolicyFunc,
	}
	data := AuthenticationPost{}
	data.GrantType = "password"
	data.Username = life.Username
	data.Password = life.Password
	postJSON, err := json.Marshal(data)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(postJSON))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+life.AuthorizationToken)
	req.Header.Add("cache-control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	authentication := AuthResponse{}
	err = json.Unmarshal(contents, &authentication)
	life.AccessToken = authentication.AccessToken
	life.GetCircles()
	life.GetPlaces()
	// life.InsertBearerToken("Bearer " + authentication.AccessToken)
	return
}
