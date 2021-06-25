package life360

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/insomniadev/martian/modules/redispub"
)

func (life *Life360) redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+life.AuthorizationToken)
	return nil
}

func (life *Life360) getURL(url string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+life.AccessToken)
	resp, err := client.Do(req)

	contents, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return contents
}

// GetCircles for Life360
func (life *Life360) GetCircles() {
	if life.AccessToken == "" {
		life.Authenticate()
	}
	resp := life.getURL(circlesURL)
	values := Circles{}
	err := json.Unmarshal(resp, &values)
	if err != nil {
		fmt.Println(err)
	}
	for _, circle := range values.Circles {
		life.Circles = append(life.Circles, circle)
	}
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
func (life *Life360) GetMembers() {
	for _, circle := range life.Circles {
		if life.AccessToken == "" {
			life.Authenticate()
		}
		resp := life.getURL(baseURL + circleURL + circle.ID + membersURL)
		values := Life360Members{}
		err := json.Unmarshal(resp, &values)
		if err != nil {
			fmt.Println(err)
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
					redispub.Service.Publish("subscriptions", member)
				}
			}
		}
	}
}

// GetPlaces for Life360
func (life *Life360) GetPlaces() {
	for _, circle := range life.Circles {
		if life.AccessToken == "" {
			life.Authenticate()
		}
		resp := life.getURL(baseURL + circleURL + circle.ID + placesURL)
		values := Places{}
		err := json.Unmarshal(resp, &values)
		if err != nil {
			fmt.Println(err)
		}
		life.Places = values.Places
	}
}

// Authenticate and return the token for Life360
func (life *Life360) Authenticate() {
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
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+life.AuthorizationToken)
	req.Header.Add("cache-control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
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
}