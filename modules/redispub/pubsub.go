package redispub

import (
	"encoding/json"

	"github.com/insomniadev/martian/integrations/config"
	"gopkg.in/redis.v2"
)

type PubSub struct {
	client *redis.Client
}

var Service *PubSub

func init() {
	url, port := config.LoadRedis()
	var client *redis.Client
	client = redis.NewTCPClient(&redis.Options{
		Addr:     url + ":" + port,
		Password: "",
		DB:       0,
		PoolSize: 10,
	})
	Service = &PubSub{client}
}

func (ps *PubSub) PublishString(channel, message string) *redis.IntCmd {
	return ps.client.Publish(channel, message)
}

func (ps *PubSub) Publish(channel string, message interface{}) *redis.IntCmd {
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	messageString := string(jsonBytes)
	return ps.client.Publish(channel, messageString)
}
