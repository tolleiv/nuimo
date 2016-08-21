package main

import (
	"log"
	"github.com/tolleiv/nuimo"
)

var done = make(chan bool)

func main() {
	device, err := nuimo.Connect()
	defer device.Disconnect()
	if err != nil {
		log.Fatalf("can't connect: %s", err)
	}

	go func(events <-chan nuimo.Event) {
		log.Printf("Nuimo ready to receive events")
		for {
			event, more := <-events
			if more {
				log.Printf("%s %x %d", event.Key, event.Raw, event.Value)
			} else {
				log.Println("no events left")
				done <- true
				return
			}
		}

	}(device.Events())

	<-done
}
