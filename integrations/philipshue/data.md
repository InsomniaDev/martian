package philipshue

import (
	"encoding/json"
	"fmt"
	bolt "homesmartie/models"
	"log"
	"strconv"
)

// Refresh gets all lights on system and puts them in the database
func (philips *PhilipsHue) Refresh() {
	lights, err := philips.GetLightsOnSystem()
	if err != nil {
		log.Fatal(err)
	}
	// philips.SaveLightsToDB(lights)
	philips.Lights = lights
	for in := range philips.Lights {
		philips.Lights[in].UUID = "philips" + strconv.Itoa(philips.Lights[in].ID)
		value, err := json.Marshal(philips.Lights[in])
		if err != nil {
			fmt.Println(err)
		}
		bolt.UpdateDevice(philips.Lights[in].UUID, string(value))
	}
}

// RetrieveAuth retrieves the url and secret from bolt
func (philips *PhilipsHue) RetrieveAuth() {

	url, err := bolt.ReadAccount("PHILIPS_HUE_URL")
	if err != nil {
		fmt.Println(err)
	}
	philips.URL = url
	secret, err := bolt.ReadAccount("PHILIPS_HUE_SECRET")
	if err != nil {
		fmt.Println(err)
	}
	philips.Secret = secret
}

// RetrieveDevices retrieves all of the devices from the bolt key value store
func (philips *PhilipsHue) RetrieveDevices() {
	values := bolt.ScanForDevicePrefix("philips")
	var lights []Light
	for _, val := range values {
		byteValue := []byte(val)
		var newLight Light
		err := json.Unmarshal(byteValue, &newLight)
		if err != nil {
			fmt.Println(err)
		}
		lights = append(lights, newLight)
	}
	philips.Lights = lights
}

// bolt.UpdateAccount("PHILIPS_HUE_URL", url)
func (philips *PhilipsHue) CreateUrl(url string) {
	bolt.UpdateAccount("PHILIPS_HUE_URL", url)
	philips.URL = url
}

// CreateAuth inserts the url and secret into bolt
func (philips *PhilipsHue) CreateAuth(secret string) {
	bolt.UpdateAccount("PHILIPS_HUE_SECRET", secret)
	philips.RetrieveAuth()
	philips.Refresh()
}

// Initialize sets up the struct to begin taking orders
func (philips *PhilipsHue) Initialize() {
	philips.RetrieveAuth()
	philips.RetrieveDevices()
}
