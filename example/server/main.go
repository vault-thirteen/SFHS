package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vault-thirteen/SFHS/server"
	"github.com/vault-thirteen/SFHS/server/settings"
)

func main() {
	cla, err := readCLA()
	mustBeNoError(err)

	var stn *settings.Settings
	stn, err = settings.NewSettingsFromFile(cla.ConfigurationFilePath)
	mustBeNoError(err)

	log.Println("Server is starting ...")
	var srv *server.Server
	srv, err = server.NewServer(stn)
	mustBeNoError(err)

	cerr := srv.Start()
	if cerr != nil {
		log.Fatal(cerr)
	}
	fmt.Println("HTTP Server: " + srv.GetListenDsn())
	fmt.Println("DB Client A: " + srv.GetDbDsnA())
	fmt.Println("DB Client B: " + srv.GetDbDsnB())

	serverMustBeStopped := srv.GetStopChannel()
	waitForQuitSignalFromOS(serverMustBeStopped)
	<-*serverMustBeStopped

	log.Println("Stopping the server ...")
	cerr = srv.Stop()
	if cerr != nil {
		log.Println(cerr)
	}
	log.Println("Server was stopped.")
	time.Sleep(time.Second)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func waitForQuitSignalFromOS(serverMustBeStopped *chan bool) {
	osSignals := make(chan os.Signal, 16)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range osSignals {
			switch sig {
			case syscall.SIGINT,
				syscall.SIGTERM:
				log.Println("quit signal from OS has been received: ", sig)
				*serverMustBeStopped <- true
			}
		}
	}()
}
