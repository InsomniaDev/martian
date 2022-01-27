package zwave

import (
	log "github.com/sirupsen/logrus"
)

func placeholder() {

}

// // ConnectToTopic will connect to the topic
// func (z *Zwave) ConnectToTopic() {
// 	// zwaveMqttURL := config.LoadZwave()

// 	const TOPIC = "zw/_EVENTS/ZWAVE_GATEWAY-server/value_changed"

// 	opts := mqtt.NewClientOptions().AddBroker("mqtt://192.168.1.19:30862")

// 	client := mqtt.NewClient(opts)
// 	if token := client.Connect(); token.Wait() && token.Error() != nil {
// 		log.Println(token.Error())
// 	}

// 	var wg sync.WaitGroup
// 	wg.Add(1)

// 	if token := client.Subscribe(TOPIC, 0, func(client mqtt.Client, msg mqtt.Message) {
// 		log.Println(string(msg.Payload()))
// 		// r := bytes.NewReader(msg.Payload())
// 		// req, _ := http.NewRequest("POST", "http://localhost:4000/zwave", r)
// 		// hclient := &http.Client{}
// 		// res, e := hclient.Do(req)
// 		// if e != nil {
// 		// 	log.Println(e)
// 		// }
// 		// defer res.Body.Close()
// 		// log.Println("response status:", res.Status)
// 		// wg.Done()
// 	}); token.Wait() && token.Error() != nil {
// 		log.Println(token.Error())
// 	}

// 	//if token := client.Publish(TOPIC, 0, false, "mymessage"); token.Wait() && token.Error() != nil {
// 	//	log.Println(token.Error())
// 	//}
// 	wg.Wait()
// }

func doSomethingWithMessage(payload []byte) {
	log.Info(string(payload))
}

// Format required to be able to publish to a mqtt topic
//if token := client.Publish(TOPIC, 0, false, "mymessage"); token.Wait() && token.Error() != nil {
//	log.Println(token.Error())
//}
