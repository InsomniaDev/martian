package sleepiq

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// CheckStatus to check the status of the API
func (si *SleepIQ) CheckStatus() (familyStatus FamilyStatus, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, SleepIqURL+"/rest/bed/familyStatus?_k="+si.Key, nil)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("DNT", "1")
	req.Header.Set("Cookie", si.Cookies)

	resp, err := client.Do(req)
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	json.Unmarshal(contents, &familyStatus)

	return
}

// Login is to login the user
func (si *SleepIQ) Login() (err error) {

	loginData := login{si.Username, si.Password}
	body, err := json.Marshal(loginData)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, SleepIqURL+"rest/login", bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("DNT", "1")
	res, err := client.Do(req)
	if err != nil {
		return
	}

	contents, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var responseBody loginResponse
	json.Unmarshal(contents, &responseBody)
	si.Key = responseBody.Key
	for _, line := range res.Header["Set-Cookie"] {
		parts := strings.Split(strings.TrimSpace(line), ";")
		parts[0] = strings.TrimSpace(parts[0])
		si.Cookies = si.Cookies + parts[0] + ";"
	}
	return
}

// Family Status
// {
// 	"beds": [
// 	  {
// 		"status": 1,
// 		"bedId": "-9223372019938824180",
// 		"leftSide": {
// 		  "isInBed": false,
// 		  "alertDetailedMessage": "No Alert",
// 		  "sleepNumber": 40,
// 		  "alertId": 0,
// 		  "lastLink": "00:00:00",
// 		  "pressure": 1197
// 		},
// 		"rightSide": {
// 		  "isInBed": false,
// 		  "alertDetailedMessage": "No Alert",
// 		  "sleepNumber": 35,
// 		  "alertId": 0,
// 		  "lastLink": "00:00:00",
// 		  "pressure": 990
// 		}
// 	  }
// 	]
//   }
