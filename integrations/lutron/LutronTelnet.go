package lutron

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/insomniadev/martian/modules/pubsub"
)

// custom io scanner splitter
// splits on either '>' or '\n' as depending on whether
// the session is at a prompt - or just sent a change event
func lutronSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	delim := strings.IndexAny(string(data), ">\n")
	if delim == -1 {
		// keep reading
		return 0, nil, nil
	}
	// else split the token
	return delim + 1, data[:delim], nil
}

// Initialize the lutron telnet communication

// func Initialize(hostName, inventoryPath string) *Lutron {
// 	inv := NewCasetaInventory(inventoryPath)

// 	l.broker = pubsub.New(10)
// 	return l
// }

// Connect to the lutron instance and start receiving messages
func (l *Lutron) Connect() error {
	conn, err := net.Dial("tcp", l.Config.URL+":"+l.Config.Port)
	if err != nil {
		return err
	}
	l.conn = conn
	loginReader := bufio.NewReader(l.conn)
	l.reader = loginReader
	// TODO turn to logging
	log.Printf("Connection established between %s and localhost.\n", l.Config.URL)
	log.Printf("Local Address : %s \n", l.conn.LocalAddr().String())
	log.Printf("Remote Address : %s \n", l.conn.RemoteAddr().String())
	message, _ := loginReader.ReadString(':')
	log.Info("Message from server: " + message + "\n")
	// send to socket
	fmt.Fprintf(conn, l.Config.Username+"\n")
	// listen for reply
	message, _ = loginReader.ReadString(':')
	log.Info("Message from server: " + message + "\n")
	fmt.Fprintf(l.conn, l.Config.Password+"\n")
	message, _ = loginReader.ReadString('>')
	log.Info("prompt ready: " + message + "\n")
	// TODO set up scanner on l.conn
	scanner := bufio.NewScanner(l.conn)
	scanner.Split(lutronSplitter)
	go func() {
		re := regexp.MustCompile(
			// ^~(?P<command>[^,]+),(?P<id>\d+),(?P<action>\d+)(?:,(?P<value>\d+(?:\.\d+)?))?$
			`^~(?P<command>[^,]+),` + // the the commmand
				`(?P<id>\d+),` +
				`(?P<action>\d+)` +
				`(?:,(?P<value>\d+` + //values are optional
				`(?:\.\d+)?` + // not all values are floats
				`))?$`) // end optional value capture
		for scanner.Scan() {
			scannedMsg := strings.TrimSpace(scanner.Text())
			// fmt.Printf("scannedMsg: %v\n", scannedMsg)
			select {
			case <-l.done:
				return
			// case l.Responses <- scannedMsg:
			default:
			}
			response := &LutronMsg{}
			groups := re.FindStringSubmatch(scannedMsg)
			if len(groups) == 0 {
				continue
			}
			lutronItems := make(map[string]string)

			// fmt.Printf("%v\n", groups)
			for i, name := range re.SubexpNames() {
				if i > 0 && i <= len(groups) {
					lutronItems[name] = groups[i]
				}
			}
			switch lutronItems["command"] {
			case "OUTPUT":
				response.Cmd = Output
			case "DEVICE":
				response.Cmd = Device
			default:
				response.Cmd = Unknown
			}
			// response.Cmd = lutronItems["command"]
			// response.Cmd = "OUTPUT".(Command)
			response.Id, err = strconv.Atoi(lutronItems["id"])
			response.Action, err = strconv.Atoi(lutronItems["action"])
			if err != nil {
				log.Warn(err.Error())
			}
			response.Type = Response
			response.Value, _ = strconv.ParseFloat(lutronItems["value"], 64)
			if err != nil {
				log.Warn(err.Error())
			}

			if response.Cmd == Output {
				for index, data := range l.Inventory {
					if data.ID == response.Id {
						oldValue := l.Inventory[index].Value
						l.Inventory[index].Value = response.Value

						// Set the value for the field
						switch value := response.Value; value {
						case 100:
							if strings.ToUpper(l.Inventory[index].Type) == "FAN" {
								l.Inventory[index].State = "HIGH"
							} else {
								l.Inventory[index].State = "ON"
							}
						case 75:
							if strings.ToUpper(l.Inventory[index].Type) == "FAN" {
								l.Inventory[index].State = "MEDIUM_HIGH"
							} else {
								l.Inventory[index].State = "DIMMED"
							}
						case 50:
							if strings.ToUpper(l.Inventory[index].Type) == "FAN" {
								l.Inventory[index].State = "MEDIUM"
							} else {
								l.Inventory[index].State = "DIMMED"
							}
						case 25:
							if strings.ToUpper(l.Inventory[index].Type) == "FAN" {
								l.Inventory[index].State = "LOW"
							} else {
								l.Inventory[index].State = "DIMMED"
							}
						default:
							if strings.ToUpper(l.Inventory[index].Type) == "FAN" {
								l.Inventory[index].Value = 0
								l.Inventory[index].State = "OFF"
							} else if response.Value > 0 {
								l.Inventory[index].State = "DIMMED"
							} else {
								l.Inventory[index].Value = 0
								l.Inventory[index].State = "OFF"
							}
						}
						pubsub.Service.Publish("subscriptions", "lutron")

						// Check to see if the values have been changed
						if oldValue != response.Value {

							// Verify that the device is one that should be automated
							for x := range l.AutomationInventory {
								if l.AutomationInventory[x] == response.Id {
									eventData := fmt.Sprintf("{\"id\":%d,\"type\":\"lutron\",\"value\":\"%s\",\"time\":\"0001-01-01T00:00:00Z\"}", response.Id, fmt.Sprintf("%f", l.Inventory[index].Value))
									eventMessage := "lutron;;" + strconv.Itoa(response.Id) + ";;" + eventData
									pubsub.Service.Publish("brain", eventMessage)
								}
							}
						}
					}
				}
			}
			// fmt.Printf("publishing %+v\n", response)
		}
	}()

	// Get all of the device current states
	go l.getAllDeviceStates()
	return nil
}

// getAllDeviceStates will loop and update fresh every minute, this is to catch if some of the devices didn't actually turn off
func (l *Lutron) getAllDeviceStates() {
	for {
		for _, device := range l.Inventory {
			// device.State = "OFF" // TODO: Commenting this out for now, if it gives an invalid response, then we will need to change this
			l.SendCommand(&LutronMsg{
				Id:   device.ID,
				Type: Get,
			})
		}
		time.Sleep(1 * time.Minute)
	}
}

func (l *Lutron) Disconnect() error {
	l.done <- true
	return l.conn.Close()
}

// TODO - how many API variations to support - need to have one
// with Fade
func (l *Lutron) SetById(id int, level float64) error {
	return l.Send(fmt.Sprintf("#OUTPUT,%d,1,%f", id, level))
}

func (l *Lutron) Send(msg string) error {
	fmt.Fprintf(l.conn, msg+"\n")
	// TODO return meaningful error
	return nil
}

func (l *Lutron) Watch(c *LutronMsg) (responses chan *LutronMsg, stop chan bool) {
	watcher := &ResponseWatcher{
		matchMsg: c,
	}
	watcher.incomming = make(chan interface{}, 5)
	watcher.Responses = make(chan *LutronMsg, 5)
	watcher.stop = make(chan bool)
	go func() {
		for {
			select {
			case msg := <-watcher.incomming:
				// match msg
				watcher.Responses <- msg.(*LutronMsg)
			case <-watcher.stop:
				close(watcher.Responses)
				return
			}
		}

	}()
	return watcher.Responses, watcher.stop
}

func (l *Lutron) SendCommand(c *LutronMsg) (resp string, err error) {
	var cmd string
	if c.Cmd == "" {
		c.Cmd = Output
	}

	switch c.Type {
	case Get:
		cmd = fmt.Sprintf("?%s,%d,1", c.Cmd, c.Id)
		// TODO confirm level and fade are 0
	case Set:
		cmd = fmt.Sprintf("#%s,%d,1,%.2f", c.Cmd, c.Id, c.Value)
	case Watch:
		// TODO
		// create mechanism to add a fmt.scanner on responses in a goroutine
		// with a dedicated channel for matches
		log.Fatal("Watch not implemented")
	}

	if c.Fade > 0.0 {
		// TODO - longer fades don't expose themselves well in the integration
		// the final value is reported for the item immediately on the sending
		// of the command. So if you set a light to dim from 100 to 10 over 20
		// seconds, the light reports out at 10 immediately. The way to perhaps
		// to approximate (as an option) is to manage the fade here, with a ticker
		cmd = fmt.Sprintf("%s,%.2f", cmd, c.Fade)
	}
	// TODO need to decide how to capture and bubble up either
	// transport/connection errors, or semantic lighting errors
	l.Send(cmd)
	return "", nil
}
