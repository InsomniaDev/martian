package philipshue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// AddIntegration adds in the creation of the API key
func (philips *PhilipsHue) AddIntegration() (err error) {
	var response []successResponse
	content := deviceType{"InsomniaDev"}

	data, err := json.Marshal(content)
	if err != nil {
		return
	}

	resp, err := http.Post(philips.URL+"/api", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	err = json.Unmarshal(contents, &response)
	if err != nil {
		return
	}

	if resp.StatusCode == 200 {
		philips.CreateAuth(response[0].Success.Username)
	}
	return
}

// TestURL tests to see if the response received is the one from philips hue
func (philips *PhilipsHue) TestURL(URL string) (valid bool, err error) {
	var errorbody []errorResponse

	resp, err := http.Get(URL + "/api/newdeveloper")
	if err != nil {
		return false, err
	}
	contents, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	err = json.Unmarshal(contents, &errorbody)
	if err != nil {
		return false, err
	}

	if errorbody[0].Err.Description == "unauthorized user" {
		philips.CreateUrl(URL)
		return true, nil
	}

	return false, nil
}

// GetLightsOnSystem Retrieve all of the lights that can be found on the system
func (philips *PhilipsHue) GetLightsOnSystem() ([]Light, error) {
	mLights := map[string]Light{}

	resp, err := http.Get(philips.URL + "/api/" + philips.Secret + "/lights")
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, &mLights)
	if err != nil {
		return nil, err
	}
	lights := make([]Light, 0)

	for i, l := range mLights {
		id, err := strconv.Atoi(i)
		l.ID = id
		if err != nil {
			return nil, err
		}
		lights = append(lights, l)
	}

	return lights, nil
}

// ChangeLightStatus turns the light on and off
func (philips *PhilipsHue) ChangeLightStatus(lightNumber int, lightOn bool) (err error) {
	client := &http.Client{}
	url := philips.URL + "/api/" + philips.Secret + "/lights/" + strconv.Itoa(lightNumber) + "/state"
	log.Println(url)
	changeState := state{lightOn}
	data, err := json.Marshal(changeState)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	log.Println(resp.StatusCode)
	contents, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	log.Println(string(contents))
	if err != nil {
		return
	}
	return
}
