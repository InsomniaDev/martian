package brain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func assembleTimeString(eventTime time.Time) (timeString string) {
	newTime := fmt.Sprintf("%s", eventTime.Format("15:04"))
	newSplit := strings.Split(newTime, ":")
	
	minutes, err := strconv.Atoi(newSplit[1])
	minuteString := ""
	if err != nil {
		fmt.Println(err)
	}
	switch {
	case minutes > 45:
		minuteString = ":45"
	case minutes > 30:
		minuteString = ":30"
	case minutes > 15:
		minuteString = ":15"
	default:
		minuteString = ":00"
	}
	timeString = newSplit[0] + minuteString
	return
}
