package harmony

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	config "github.com/insomniadev/martian/integrations/config"
	rg "github.com/redislabs/redisgraph-go"
)

const (
	redisGraph = "Devices"
)

// InsertIntoGraph is used to insert a harmony device into the graph
func InsertIntoGraph(name string, activityID string, actions string) Device {
	url, port := config.LoadRedis()
	conn, _ := redis.Dial("tcp", url+":"+port)
	defer conn.Close()
	device := Device{
		Name:       name,
		ActivityID: activityID,
		Actions:    actions,
	}
	graph := rg.GraphNew(redisGraph, conn)
	query := "CREATE (:activity{name:'" + device.Name + "', activityId: '" + device.ActivityID + "', actions: '" + device.Actions + "'})}) "
	query += "MERGE (:devicetype{name:'harmony'}) "
	query += "MATCH (a:activity{deviceId:'" + device.ActivityID + "'}), (c:devicetype{name:'harmony'}) "
	query += "CREATE (a) -[:OF_TYPE]-> (c)"
	graph.Query(query)
	return device
}

// RetrieveAllNodes will retrieve all of the nodes and print them prettily
func RetrieveAllNodes() []*Device {
	url, port := config.LoadRedis()
	conn, _ := redis.Dial("tcp", url+":"+port)
	defer conn.Close()
	var Devices []*Device

	graph := rg.GraphNew(redisGraph, conn)
	query := "MATCH (device:activity) -[:OF_TYPE]-> (:devicetype{name:'harmony'}) RETURN device.name,device.activityId,device.actions"

	result, _ := graph.Query(query)
	for result.Next() {
		r := result.Record()
		name, _ := r.Get("device.name")
		activityID, _ := r.Get("device.activityId")
		// actions, _ := r.Get("device.actions")
		dev := Device{
			Name:       fmt.Sprintf("%v", name),
			ActivityID: fmt.Sprintf("%v", activityID),
		}
		Devices = append(Devices, &dev)
	}

	for _, v := range Devices {
		fmt.Printf("Device name: %s\n", v.Name)
	}
	return Devices
}
