package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/insomniadev/martian/internal/api"
	"github.com/insomniadev/martian/internal/brain"
	"github.com/insomniadev/martian/internal/database"
	"github.com/insomniadev/martian/internal/integrations/config"
	"github.com/insomniadev/martian/internal/integrations/hubitat"
	log "github.com/sirupsen/logrus"
)

var sigs = make(chan os.Signal, 1)

func main() {
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go handleAppClose()

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	log.Info("Martian is starting up")
	config.LoadConfiguration()

	brain.Brainiac.SayHello()

	// TODO: This should be set as a configuration something something
	hubitat.Instance.GetAllDeviceStatus()

	// Start up the API server
	api.StartApi()
}

func handleAppClose() {
	sig := <-sigs
	log.Println("Handling application termination")
	log.Println(sig)
	if stored := brain.Brainiac.StoreMemoryData(); !stored {
		log.Println("Failed to store the memory data")
	}
	database.MartianData.Connection.Close()
	os.Exit(0)
}
