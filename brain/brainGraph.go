package brain

import (
	"fmt"
	"strconv"

	"github.com/gomodule/redigo/redis"

	rg "github.com/redislabs/redisgraph-go"
)

const (
	tempGraph = "temporary"
	timeGraph = "time"
	permGraph = "automations"
)

func (b *Brain) cleanTempDatabase() {
	conn, _ := redis.Dial("tcp", b.redisURL+":"+b.redisPort)
	defer conn.Close()

	graph := rg.GraphNew(tempGraph, conn)
	query := "MATCH (a) DELETE a"
	_, err := graph.Query(query)
	if err != nil {
		fmt.Println(err)
	}
}

func (b *Brain) checkForTimeAutomations(eventTime string) (exists bool) {
	b.AutomationEvent = nil

	conn, _ := redis.Dial("tcp", b.redisURL+":"+b.redisPort)
	defer conn.Close()

	exists = false
	graph := rg.GraphNew(permGraph, conn)

	query := fmt.Sprintf("MATCH (a:event) <-[b:RELATES]-> (c:time{military:'%s'})  RETURN a.deviceId, a.state, a.device order by b.weight", eventTime)

	resp, err := graph.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	for resp.Next() {
		r := resp.Record()
		deviceIDString, _ := r.Get("a.deviceId")
		deviceID, _ := strconv.Atoi(fmt.Sprintf("%v", deviceIDString))
		deviceState, _ := r.Get("a.state")
		deviceType, _ := r.Get("a.device")
		b.AutomationEvent = append(b.AutomationEvent, Event{
			ID:    deviceID,
			Type:  fmt.Sprintf("%v", deviceType),
			Value: fmt.Sprintf("%v", deviceState),
		})
		exists = true
	}
	return
}

func (b *Brain) checkForEventAutomations() (exists bool) {
	b.AutomationEvent = nil

	conn, _ := redis.Dial("tcp", b.redisURL+":"+b.redisPort)
	defer conn.Close()

	exists = false
	graph := rg.GraphNew(permGraph, conn)

	query := fmt.Sprintf("MATCH (a:event{deviceId:'%s',state:'%s',device:'%s'}) -[b:RELATES]-> (c) RETURN c.deviceId, c.state, c.device order by b.weight", strconv.Itoa(b.CurrentEvent.ID), b.CurrentEvent.Value, b.CurrentEvent.Type)

	resp, err := graph.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	for resp.Next() {
		r := resp.Record()
		deviceIDString, _ := r.Get("c.deviceId")
		deviceID, _ := strconv.Atoi(fmt.Sprintf("%v", deviceIDString))
		deviceState, _ := r.Get("c.state")
		deviceType, _ := r.Get("c.device")
		b.AutomationEvent = append(b.AutomationEvent, Event{
			ID:    deviceID,
			Type:  fmt.Sprintf("%v", deviceType),
			Value: fmt.Sprintf("%v", deviceState),
		})
		exists = true
	}
	return
}

func (b *Brain) storeEventGraph() int {
	conn, _ := redis.Dial("tcp", b.redisURL+":"+b.redisPort)
	defer conn.Close()

	graph := rg.GraphNew(tempGraph, conn)

	query := fmt.Sprintf("MERGE (a:event{deviceId:'%s',state:'%s',device:'%s'})", strconv.Itoa(b.LastEvent.ID), b.LastEvent.Value, b.LastEvent.Type)
	query += fmt.Sprintf(" MERGE (b:event{deviceId:'%s',state:'%s',device:'%s'})", strconv.Itoa(b.CurrentEvent.ID), b.CurrentEvent.Value, b.CurrentEvent.Type)
	query += " MERGE (a) -[c:RELATES]-> (b) ON CREATE SET c.weight=1 ON MATCH SET c.weight=c.weight+1 RETURN c.weight"

	resp, err := graph.Query(query)
	if err != nil {
		fmt.Println("Placeholder")
	}
	var weight int
	for resp.Next() {
		r := resp.Record()
		weightString, _ := r.Get("c.weight")
		weight, _ = strconv.Atoi(fmt.Sprintf("%v", weightString))
	}

	// Move to the permanent automation store
	if weight > 3 {
		permG := rg.GraphNew(permGraph, conn)
		_, err := permG.Query(query)
		if err != nil {
			fmt.Println("Placeholder")
		}

		// DELETE the relationship
		query = fmt.Sprintf("MATCH (a:event{deviceId:'%s',state:'%s',device:'%s'}) -[c:RELATES]-> (b:event{deviceId:'%s',state:'%s',device:'%s'}) ", strconv.Itoa(b.LastEvent.ID), b.LastEvent.Value, b.LastEvent.Type, strconv.Itoa(b.CurrentEvent.ID), b.CurrentEvent.Value, b.CurrentEvent.Type)
		query += "DELETE c"
		_, err = graph.Query(query)
		if err != nil {
			fmt.Println("Placeholder")
		}
	}

	return weight
}

func (b *Brain) storeTimeGraph() int {
	conn, _ := redis.Dial("tcp", b.redisURL+":"+b.redisPort)
	defer conn.Close()

	timeInstanceForEvent := assembleTimeString(b.CurrentEvent.Time)

	graph := rg.GraphNew(tempGraph, conn)

	query := fmt.Sprintf("MERGE (a:event{deviceId:'%s',state:'%s',device:'%s'})", strconv.Itoa(b.CurrentEvent.ID), b.CurrentEvent.Value, b.CurrentEvent.Type)
	query += fmt.Sprintf(" MERGE (b:time{military:'%s'})", timeInstanceForEvent)
	query += " MERGE (a) <-[c:RELATES]-> (b) ON CREATE SET c.weight=1 ON MATCH SET c.weight=c.weight+1 RETURN c.weight"

	fmt.Println(query)
	resp, err := graph.Query(query)
	if err != nil {
		fmt.Println("Placeholder")
	}
	var weight int
	for resp.Next() {
		r := resp.Record()
		weightString, _ := r.Get("c.weight")
		weight, _ = strconv.Atoi(fmt.Sprintf("%v", weightString))
	}

	// Move to the permanent automation store
	if weight > 3 {
		permG := rg.GraphNew(permGraph, conn)
		_, err := permG.Query(query)
		if err != nil {
			fmt.Println("Placeholder")
		}

		// DELETE the relationship
		query = fmt.Sprintf("MATCH (a:event{deviceId:'%s',state:'%s',device:'%s'}) <-[c:RELATES]-> (b:time{military:'%s'}) ", strconv.Itoa(b.CurrentEvent.ID), b.CurrentEvent.Value, b.CurrentEvent.Type, timeInstanceForEvent)
		query += "DELETE c"
		fmt.Println(query)
		_, err = graph.Query(query)
		if err != nil {
			fmt.Println("Placeholder")
		}
	}
	return weight
}
