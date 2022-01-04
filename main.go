package main

import (
	"log"
	"time"

	"github.com/insomniadev/martian/graphql"
	"github.com/insomniadev/martian/modules/cache"
)

func main() {
	// testLocalCache()
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
		log.Println("value not found")
	} else {
		log.Println("the found value: " + value.(something).value)
	}
}
