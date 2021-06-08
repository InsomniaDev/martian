package lutron

import (
	"fmt"
	"strconv"

	"github.com/gomodule/redigo/redis"
	config "github.com/insomniadev/martian/integrations/config"
	rg "github.com/redislabs/redisgraph-go"
)

const (
	redisGraph = "Devices"
)

// InsertLutronGraph is used to insert a lutron device into the graph
func InsertLutronGraph(name string, deviceID int, areaName string, deviceType string) LDevice {
	url, port := config.LoadRedis()
	conn, _ := redis.Dial("tcp", url+":"+port)
	defer conn.Close()
	device := LDevice{
		Name:     name,
		ID:       deviceID,
		AreaName: areaName,
		Type:     deviceType,
	}
	graph := rg.GraphNew(redisGraph, conn)

	// Create the new connections
	query := "CREATE (:lutron{name:'" + device.Name + "', deviceId:" + strconv.Itoa(device.ID) + "}) "
	query += "MERGE (:area{name:'" + device.AreaName + "'}) "
	query += "MERGE (:devicetype{name:'" + device.Type + "'}) "
	_, err := graph.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	// Create the relationships
	query = "MATCH (a:lutron{deviceId:" + strconv.Itoa(device.ID) + "}), (b:area{name:'" + device.AreaName + "'}), (c:devicetype{name:'" + device.Type + "'}) "
	query += "CREATE (b) <-[:RESIDES_IN]- (a) -[:OF_TYPE]-> (c) "
	_, err = graph.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	return device
}

// RetrieveLutronNodes will retrieve all of the nodes and print them prettily
func RetrieveLutronNodes() []*LDevice {
	url, port := config.LoadRedis()
	conn, _ := redis.Dial("tcp", url+":"+port)
	defer conn.Close()
	var Devices []*LDevice

	graph := rg.GraphNew(redisGraph, conn)
	query := "MATCH (areaName:area) <-[:RESIDES_IN]- (device:lutron) -[:OF_TYPE]-> (type:devicetype) RETURN areaName.name,device.name,device.deviceId,type.name"

	result, _ := graph.Query(query)
	for result.Next() {
		r := result.Record()
		areaName, _ := r.Get("areaName.name")
		device, _ := r.Get("device.name")
		deviceID, _ := r.Get("device.deviceId")
		deviceType, _ := r.Get("type.name")
		numDeviceID, _ := strconv.Atoi(fmt.Sprintf("%v", deviceID))
		Device := LDevice{
			Name:     fmt.Sprintf("%v", device),
			ID:       numDeviceID,
			AreaName: fmt.Sprintf("%v", areaName),
			Type:     fmt.Sprintf("%v", fmt.Sprintf("%v", deviceType)),
		}
		Devices = append(Devices, &Device)
	}

	for _, v := range Devices {
		fmt.Printf("Device name: %s\n", v.Name)
	}
	return Devices
}
