package brain

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/insomniadev/martian/internal/database"
	"github.com/insomniadev/martian/pkg/pubsub"
)

/*
*

- energy efficiency
  - Set to look at every ten minute block and turn on the device if it is required to be activated
  - allow for automation to override device
  - should energy efficiency cut down be more aggressive?

- automating time efficiency
  - only set weight if the light turning on isn't through an automation?
  - After 5 or greater weight set as automated
  - Go through and check times once per week

- activating automation
  - check every ten minutes

- data design:
key - timeBlock (ie., "17.3")
value -	{
		timeBlocks: [
			{
				uniqueHash: "",
				weight: 0,
				automated: false,
				timeKey: "17.3"
			}
		]
	}

- integration will dictate if a device is "activated"
- there is a permanent list of "activated" devices that are tracked
*/

// getCurrentTimeStampByTen will return the timestamp, ie: for 5:31PM it will return 17.3
func getCurrentTimeStampByTen() string {
	if s, err := strconv.ParseFloat(strings.Replace(time.Now().Format("15:4"), ":", ".", -1), 32); err == nil {
		ratio := math.Pow(10, float64(1))
		return strings.Replace(fmt.Sprintf("%f", (math.Round(s*ratio)/ratio)), "0", "", -1)
	}
	return ""
}

// organizeActivatedDevices will keep track of devices in the time memory
func (b *Brain) organizeActivatedDevices(event Event, activated bool) {
	// Set the device to being tracked by time on if it is activated or not
	// 		currently the time tracking is only targeted on when switches are activated
	for i := range b.Omniscience {
		if b.Omniscience[i].DeviceId == event.DeviceId && b.Omniscience[i].Integration == event.Integration {
			b.Omniscience[i].TimeTracked = activated
		}
	}
}

// checkTimeAutomations will constantly peruse what is active
func (b *Brain) checkTimeAutomations() {
	// get the current timestamp for the automation
	if currentTimeBlock := getCurrentTimeStampByTen(); currentTimeBlock != "" {

		// Update weights first for the ones currently being time tracked
		devicesToUpdateWeights := []database.Device{}
		for i := range b.Omniscience {
			if b.Omniscience[i].TimeTracked {
				devicesToUpdateWeights = append(devicesToUpdateWeights, database.Device{
					DeviceId:    b.Omniscience[i].DeviceId,
					Integration: b.Omniscience[i].Integration,
				})
			}
		}

		// Logging lines
		type logger struct {
			Data []database.Device `json:"updatingDevices"`
		}
		daterz := logger{Data: devicesToUpdateWeights}
		if data, err := json.Marshal(daterz); err == nil {
			log.Println("data to store and update weights: ", string(data))
		}

		database.MartianData.UpdateWeightsForTimeTable(currentTimeBlock, devicesToUpdateWeights)

		// grab all of the devices that need to be activated during this time block
		devices := database.MartianData.GetAutomatedTimeTableEntry(currentTimeBlock)

		// turn on all of the devices that need to be activated for this time block
		// 		set the device in omniscience to be currently time tracked
		for i := range devices {
			for bi := range b.Omniscience {
				if devices[i].DeviceId == b.Omniscience[bi].DeviceId && devices[i].Integration == b.Omniscience[bi].Integration && !b.Omniscience[bi].TimeTracked {
					// update device to be time tracked
					b.Omniscience[bi].TimeTracked = true
					// update device as being time automated since that is how it was triggered
					b.Omniscience[bi].TimeAutomated = true

					// activate the device
					pubsub.Service.Publish(devices[i].Integration, fmt.Sprintf("%s;;%s", devices[i].DeviceId, "activate"))
					break
				}
			}
		}

		// check if any devices are currently time automated and are no longer valid, shut off if they are found
		for bi := range b.Omniscience {
			if b.Omniscience[bi].TimeAutomated {
				needsToBeTerminated := true
				for i := range devices {
					if devices[i].DeviceId == b.Omniscience[bi].DeviceId && devices[i].Integration == b.Omniscience[bi].Integration {
						needsToBeTerminated = false
						break
					}
				}
				// If the device needs to be terminated and is currently automated through the time learning
				if needsToBeTerminated && b.Omniscience[bi].TimeAutomated {

					// Set the device to no longer be time automated
					b.Omniscience[bi].TimeAutomated = false
					// set the device to no longer be time tracked
					b.Omniscience[bi].TimeTracked = false

					// Add to the energy efficiency time
					// TODO: if energy efficiency is not configured then we need to just turn off the device
					b.rememberEnergyTime(Event{
						Integration: b.Omniscience[bi].Integration,
						DeviceId:    b.Omniscience[bi].DeviceId,
						Status:      "time",
					})
				}
			}
		}
	}
}

func (b *Brain) checkForTimeActivation(newEvent Event, fullEvent string) {
	// store the event in the timetable memory for activated devices
	if boolValue, err := strconv.ParseBool(strings.Split(fullEvent, ";;")[5]); err == nil {
		// add the device to the activated devices if not activated, or remove if false now
		b.organizeActivatedDevices(newEvent, boolValue)
	}
}

func (b *Brain) setTimeAutomations() {
	// we are checking this weekly, so the automation needs to occur at least 6 times out of 7
	weightLimit := 6

	// set the automated time entries to zero
	// set remaining time entries as automated that are already not automated
	allTimeEntries := database.MartianData.GetAllTimeTables()

	// Decrease time that devices are on by a minute, this is now a weekly process
	// TODO: Time to decrease by should be a configuration
	database.MartianData.ImproveEnergyEfficiency(1)

	if err := database.MartianData.RecreateTimeBucket(); err == nil {
		for i := range allTimeEntries {
			newTimeBlock := database.TimeTable{}

			for ti := range allTimeEntries[i].Times {
				// Automate the entry if it has a weight greater than the set weight limit
				// 		OR
				// The entry is already automated and the weight is greater than zero
				if (allTimeEntries[i].Times[ti].Weight >= weightLimit && !allTimeEntries[i].Times[ti].Automated) || (allTimeEntries[i].Times[ti].Weight >= 0 && allTimeEntries[i].Times[ti].Automated) {

					newTimeBlock.Times = append(newTimeBlock.Times, database.TimeBlocks{
						TimeKey:    allTimeEntries[i].Times[ti].TimeKey,
						UniqueHash: allTimeEntries[i].Times[ti].UniqueHash,
						Automated:  true,
						Weight:     0,
					})
				}
			}

			if len(newTimeBlock.Times) > 0 {
				database.MartianData.PutTimeTableValues(newTimeBlock.Times[0].TimeKey, newTimeBlock)
			}
		}
	}

	// TODO:
	// energy efficiency - take off one timeblock at the end of the chain every week
	// cycle through the time database per device to get all of the entries that are consecutive and then remove the last one of each consecutive block
	// this should only happen if we aren't automating anything for that device
}
