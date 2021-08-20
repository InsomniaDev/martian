package kasa

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/insomniadev/martian/integrations/config"
	rg "github.com/redislabs/redisgraph-go"
)

const (
	kasaGraph = "Devices"
)

// InsertKasaGraph inserts device known into graph database
func InsertKasaGraph(ipAddress string, name string) {
	url, port := config.LoadRedis()

	conn, _ := redis.Dial("tcp", url+":"+port)
	defer conn.Close()
	graph := rg.GraphNew(kasaGraph, conn)

	query := "CREATE (:kasa{name:'" + name + "', ipAddress:'" + ipAddress + "'}) "
	query += "MERGE (:devicetype{name:'plug'}) "
	_, err := graph.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	query = "MATCH (a:kasa{name:'" + name + "', ipAddress:'" + ipAddress + "'}), (b:devicetype{name:'plug'}) "
	query += "CREATE (a) -[:OF_TYPE]-> (b) "
	_, err = graph.Query(query)
	if err != nil {
		fmt.Println(err)
	}
}

// UpdateAreaForKasaDevice will update an area for the
func UpdateAreaForKasaDevice(ipAddress string, areaName string) {
	url, port := config.LoadRedis()

	conn, _ := redis.Dial("tcp", url+":"+port)
	defer conn.Close()
	graph := rg.GraphNew(kasaGraph, conn)

	query := "MATCH (:area) <-[r:RESIDES_IN]- (:kasa{ipAddress:'" + ipAddress + "'}) "
	query += "DELETE r "
	_, err := graph.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	query = "MATCH (a:kasa{ipAddress:'" + ipAddress + "'}) "
	query += "MERGE (b:area{name:'" + areaName + "'}) "
	query += "MERGE (b) <-[:RESIDES_IN]- (a) "
	_, err = graph.Query(query)
	if err != nil {
		fmt.Println(err)
	}
}

// RetrieveKasaNodes will retrieve all of the nodes in the graph
func RetrieveKasaNodes() []KasaDevice {
	url, port := config.LoadRedis()

	conn, _ := redis.Dial("tcp", url+":"+port)
	defer conn.Close()
	graph := rg.GraphNew(kasaGraph, conn)
	var plugs []KasaDevice

	query := "MATCH (device:kasa) -[:OF_TYPE]-> (devicetype{name:'plug'}) "
	query += "RETURN device.ipAddress, device.name "
	result, _ := graph.Query(query)
	for result.Next() {
		r := result.Record()
		ipAddress, _ := r.Get("device.ipAddress")
		name, _ := r.Get("device.name")
		plug := NewPlug(fmt.Sprintf("%v", ipAddress))
		plug.Name = fmt.Sprintf("%v", name)
		plugs = append(plugs, plug)
	}

	query = "MATCH (areaName:area) <-[:RESIDES_IN]- (device:kasa) -[:OF_TYPE]-> (devicetype{name:'plug'}) "
	query += "RETURN areaName.name, device.ipAddress, device.name "
	result, _ = graph.Query(query)
	for result.Next() {
		r := result.Record()
		ipAddress, _ := r.Get("device.ipAddress")
		areaName, _ := r.Get("areaName.name")
		name, _ := r.Get("device.name")
		found := false
		for i, pl := range plugs {
			if pl.IPAddress == ipAddress {
				found = true
				plugs[i].AreaName = fmt.Sprintf("%v", areaName)
			}
		}
		if !found {
			plug := NewPlug(fmt.Sprintf("%v", ipAddress))
			plug.Name = fmt.Sprintf("%v", name)
			plug.AreaName = fmt.Sprintf("%v", areaName)
			plugs = append(plugs, plug)
		}
	}

	return plugs
}
