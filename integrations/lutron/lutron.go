package lutron

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/insomniadev/martian/database"
)

// Init starts up the lutron instance
func Init(configuration string) (Lutron, error) {
	var lutron Lutron
	err := json.Unmarshal([]byte(configuration), &lutron)
	if err != nil {
		// TODO: Change this away from being a panic
		panic(err)
	}

	if len(lutron.Inventory) == 0 {
		fileContents, err := loadIntegrationFile(lutron.Config)
		if err != nil {
			return Lutron{}, err
		}
		for _, device := range fileContents.LIPIDList.Zones {
			lutron.Inventory = append(lutron.Inventory, &LDevice{
				Name:       device.Name,
				ID:         device.ID,
				AreaName:   device.Area.Name,
				Type:       device.Type,
				LutronName: strings.Replace(device.Area.Name+" "+device.Name, " ", "_", -1),
			})
		}
	}

	lutron.done = make(chan bool)

	lutron.Connect()
	return lutron, nil
}

func loadIntegrationFile(config database.LutronConfig) (CasetaIntegrationFile, error) {
	// TODO: Need to add setup screen in the UI to determine what type the device is, currently it is manually being put in
	jsonfile, err := os.Open("./config/" + config.File)

	fileContents := CasetaIntegrationFile{}
	jsonParser := json.NewDecoder(jsonfile)
	if err = jsonParser.Decode(&fileContents); err != nil {
		return CasetaIntegrationFile{}, err
	}
	return fileContents, nil
}

// UpdateSelectedDevices will go through and update the devices as selected or not selected
func (h *Lutron) UpdateSelectedDevices(selectedDevices []int, addDevices bool, automationDevice bool) error {

	// Compare either through the automation or interface selections
	if automationDevice {
		h.AutomationInventory = checkIfDeviceIsInList(h.Inventory, h.AutomationInventory, selectedDevices, addDevices)
	} else {
		h.InterfaceInventory = checkIfDeviceIsInList(h.Inventory, h.InterfaceInventory, selectedDevices, addDevices)
	}
	var db database.Database
	err := db.PutIntegrationValue("lutron", h)
	if err != nil {
		log.Println(err)
	}

	return nil
}

// checkIfDeviceIsInList is an internal method to see if the value already exists in the list
func checkIfDeviceIsInList(allDevices []*LDevice, alreadyChosenDevices []int, selectedDevices []int, addDevices bool) []int {
	var newlySelectedDevices []int
	// Cycle through all of the available devices for HomeAssistant
	for _, availableDevice := range allDevices {

		selectedDeviceExists := false
		// Cycle through all of the already selected devices to see if there is a match
		for _, selectedDeviceID := range alreadyChosenDevices {

			// If this available device is already selected, then set selectedDeviceExists as true
			if availableDevice.ID == selectedDeviceID {
				selectedDeviceExists = true
				break // break out of the selectedDevice cycle since there is a match
			}
		}

		availableDeviceIsNowSelected := false
		// Cycle through the selectedDevices parameter to see if the available device has a match
		for _, newlySelectedDevice := range selectedDevices {

			// If the available device matches with the newly selected criteria then set as true
			if availableDevice.ID == newlySelectedDevice {
				availableDeviceIsNowSelected = true
			}
		}

		// IF the device is already selected, and is one of the newly selected, and set to be added
		//			IF addDevices is FALSE, then it will be removed from the selectedDevices slice
		if availableDeviceIsNowSelected && addDevices {
			newlySelectedDevices = append(newlySelectedDevices, availableDevice.ID)
		} else if selectedDeviceExists && !availableDeviceIsNowSelected { // IF it was already selected
			newlySelectedDevices = append(newlySelectedDevices, availableDevice.ID)
		}
	}
	return newlySelectedDevices
}

// EditDeviceConfiguration will go through and update information for the passed in device
func (k *Lutron) EditDeviceConfiguration(device LDevice, removeEdit bool) error {

	// Cycle through all of the devices, interfaceDevices, and automatedDevices and update
	for i := range k.Inventory {
		if k.Inventory[i].ID == device.ID {
			k.Inventory[i] = &device
		}
	}

	// Update device in edited section, if not found then add it to the edited devices
	foundDevice := false
	for i := range k.EditedInventory {
		if k.EditedInventory[i].ID == device.ID {
			k.EditedInventory[i] = &device
			foundDevice = true
		}
	}
	if !foundDevice {
		k.EditedInventory = append(k.EditedInventory, &device)
	}

	// Save in the database
	var db database.Database
	err := db.PutIntegrationValue("lutron", k)
	if err != nil {
		log.Println(err)
	}

	return nil
}
