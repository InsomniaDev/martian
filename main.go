package main

import (
	"os"
	"time"

	"github.com/insomniadev/martian/brain"
	"github.com/insomniadev/martian/graphql"
	"github.com/insomniadev/martian/modules/cache"
	log "github.com/sirupsen/logrus"
)

func main() {
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	log.Info("Martian is starting up")
	// testLocalCache()
	brain.Brainiac.SayHello()

	graphql.Graphql()

	// mainBrain.Init()
	// harmony.RetrieveAllNodes()
	// go graphql.Graphql()
	// server.Start()
}

func testLocalCache() {

	localCache := cache.LocalCache{}
	localCache.Init()

	type something struct {
		value string
	}
	testing := something{value: "successful"}
	localCache.Set("new", testing)

	// wait for value to pass through buffers
	time.Sleep(10 * time.Millisecond)

	value, wasRetrieved := localCache.Get("new")
	if !wasRetrieved {
		log.Info("value not found")
	} else {
		log.Info("the found value: " + value.(something).value)
	}
}
