package brain

import (
	"fmt"
	"time"

	"github.com/insomniadev/martian/internal/database"
	"github.com/insomniadev/martian/pkg/pubsub"
	log "github.com/sirupsen/logrus"
)

// efficiencyTimeCheck is used to check when light is turned off too soon and update efficiency
type efficiencyTimeCheck struct {
	uniqueHash          string
	timestampExpiration time.Time
}

// getTimeDuration returns the set time for the device
func getTimeDuration(deviceId, integration *string) (duration time.Duration, exists bool) {
	dataEvent := database.MartianData.GetDevice(*deviceId, *integration)
	if dataEvent.EnergyEfficiencyMinutes == -1 {
		return time.Duration(0) * time.Minute, false
	}
	return time.Duration(dataEvent.EnergyEfficiencyMinutes) * time.Minute, true
}

// rememberEnergyTime is responsible for assembling the business logic behind devices
func (b *Brain) rememberEnergyTime(event Event) {
	currentTime := time.Now()

	// Get the expiration time from the database
	duration, energySavingsExists := getTimeDuration(&event.DeviceId, &event.Integration)

	for i := range b.Omniscience {
		if event.DeviceId == b.Omniscience[i].DeviceId && event.Integration == b.Omniscience[i].Integration {
			// Ignore this device if it is not timeautomated and has energy savings
			if !b.Omniscience[i].TimeAutomated && energySavingsExists {

				b.Omniscience[i].EnergyTracked = true
				b.Omniscience[i].EnergyExpiration = currentTime.Add(duration)

				log.Debugln("store energy: ", event.Status)

				// If the device was automated, then the newest triggered device matters the most
				if event.TriggeredDeviceId != "" {
					log.Debugln("update triggered devices ", b.Omniscience[i].DeviceId)
					b.Omniscience[i].TriggeredDeviceId = event.TriggeredDeviceId
					b.Omniscience[i].TriggeredDeviceStatus = event.TriggeredDeviceStatus
				}
			}
		}

		// Compare and check the device against the triggered devices
		if event.DeviceId == b.Omniscience[i].TriggeredDeviceId && event.Status != b.Omniscience[i].TriggeredDeviceStatus {

			// If this was an automated event and the triggered device has a new status
			// 		then we should update and set the device to now expire at the required amount of time
			// ie. Motion no longer active and light turns off in ten minutes
			b.Omniscience[i].TriggeredDeviceId = ""
			b.Omniscience[i].TriggeredDeviceStatus = ""

			// Update the expiration time for the device
			duration, _ := getTimeDuration(&b.Omniscience[i].DeviceId, &b.Omniscience[i].Integration)
			// b.Omniscience[i].EnergyTracked = true
			b.Omniscience[i].EnergyExpiration = currentTime.Add(duration)
		}

	}
}

// energyEfficiency is responsible for constantly being aware of the current energy savings
func (b *Brain) energyEfficiency() {
	currentTime := time.Now()

	for i, timeEvent := range b.Omniscience {
		// Check if (1) the triggered device status has changed or (2) if the device status is "inactive" and not likely to change soon
		if (timeEvent.TriggeredDeviceId == "" || timeEvent.TriggeredDeviceStatus == "active") && (timeEvent.EnergyTracked && timeEvent.EnergyExpiration.Before(currentTime)) {
			// Set the hash so that when we get the next event it won't attempt to learn from efficiency savings
			uniquehash := database.GetUniqueHash(timeEvent.DeviceId, timeEvent.Integration)
			lastTimeShutoffUniqueHash = uniquehash

			// Turn off the event
			pubsub.Service.Publish(timeEvent.Integration, fmt.Sprintf("%s;;%s", timeEvent.DeviceId, "energy"))

			// Add uniquehash to array to check against for turning back on, with timestamp
			lastTimeShutOffDevices = append(lastTimeShutOffDevices, efficiencyTimeCheck{
				uniqueHash:          uniquehash,
				timestampExpiration: time.Now().Add(5 * time.Second), // set expiration to five seconds, time is getting set a lot higher much faster than we want
			})

			// Set to no longer tracking the energy since we are going to turn it off
			b.Omniscience[i].EnergyTracked = false
		}
	}
}

func removeIndex(s []Event, index int) []Event {
	return append(s[:index], s[index+1:]...)
}

// checkForDeviceEnergyEfficiency checks if the device was recently disabled by energy efficiency
//
//	if it was recently disabled it will return true
func checkForDeviceEnergyEfficiency(uniquehash string) bool {
	// Check to make sure that the light was turned off less than a minute ago prior to updating

	hashTimeToUpdateBy := 10
	currentTime := time.Now()

	// Remove the expired events from the array
	devicesToCheckAgainst := []efficiencyTimeCheck{}
	for i := range lastTimeShutOffDevices {
		// If we are to expire the entry after the current time, then we should keep this entry
		if lastTimeShutOffDevices[i].timestampExpiration.After(currentTime) {
			devicesToCheckAgainst = append(devicesToCheckAgainst, lastTimeShutOffDevices[i])
		}
	}

	// Update our array with valid timestamps
	lastTimeShutOffDevices = devicesToCheckAgainst

	// Update the efficiencyMinutes in the database
	for i := range lastTimeShutOffDevices {
		// If we are still in the valid time and the hashes match, then we need to update
		if lastTimeShutOffDevices[i].uniqueHash == uniquehash {
			database.MartianData.UpdateEfficiencyTime(uniquehash, hashTimeToUpdateBy)
			return true
		}
	}

	// Thoughts:
	// Possibly update a column that it was changed and if not, then decrease it by one minute per day to slowly get to more efficiency
	// If clicked again after turning off, set the turn off time plus ten minutes, allow this as a configuration in the future

	return false
}
