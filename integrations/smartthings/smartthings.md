package smartthings

import (
	"encoding/json"
	"log"
	"io/ioutil"
)

// issueCommand sends a given command to an URI and returns the contents
func (st *SmartThings) issueCommand(cmd string) ([]byte, error) {
	uri := st.Endpoint + cmd
	resp, err := st.Client.Get(uri)
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return contents, nil
}

// GetDevices will return all devices
func (st *SmartThings) GetDevices() (deviceList []DeviceList, err error) {

	contents, err := st.issueCommand("/devices")
	if err != nil {
		return nil, err
	}
	log.Println(string(contents))
	if err := json.Unmarshal(contents, &deviceList); err != nil {
		return nil, err
	}
	return deviceList, nil
}

// ExecuteCommand Executes command to change the state of the IoT device
func (st *SmartThings) ExecuteCommand(id string, command string) (err error) {
	_, err = st.issueCommand("/devices/" + id + "/" + command)
	if err != nil {
		return
	}
	return
}

// GetDeviceCommands gets a list of commands
func (st *SmartThings) GetDeviceCommands(id string) (deviceCommands []DeviceCommand, err error) {
	contents, err := st.issueCommand("/devices/" + id + "/commands")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(contents, &deviceCommands); err != nil {
		return nil, err
	}
	return deviceCommands, nil
}

// GetDeviceInfo Retrieves information about the device
func (st *SmartThings) GetDeviceInfo(id string) (deviceInfo DeviceInfo, err error) {
	contents, err := st.issueCommand("/devices/" + id)
	if err != nil {
		return
	}
	log.Println(string(contents))

	err = json.Unmarshal(contents, &deviceInfo)
	if err != nil {
		return
	}
	return
}

// GetDeviceEvents Retrieve events about the device
func (st *SmartThings) GetDeviceEvents(id string) (deviceInfo DeviceInfo, err error) {
	contents, err := st.issueCommand("/devices/" + id + "/events")
	if err != nil {
		return
	}

	err = json.Unmarshal(contents, &deviceInfo)
	if err != nil {
		return
	}
	return
}

// https://github.com/marcopaganini/gosmart/blob/master/gosmart.go
