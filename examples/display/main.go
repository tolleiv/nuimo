package main

import (
	"log"
	"github.com/tolleiv/nuimo"
	"time"
)

var done = make(chan bool)

func main() {
	device, err := nuimo.Connect()
	defer device.Disconnect()
	if err != nil {
		log.Fatalf("can't connect: %s", err)
	}

	for {
		matrix := nuimo.DisplayMatrix(
			0,0,0,0,0,0,0,0,0,
			0,1,1,0,0,0,1,1,0,
			1,0,0,1,0,1,0,0,1,
			1,0,0,1,0,1,0,0,1,
			0,1,1,1,0,0,1,1,0,
			0,0,0,1,0,0,0,0,0,
			0,0,1,0,0,0,0,0,0,
			0,0,0,0,0,0,0,0,0,
			0,0,0,0,0,0,0,0,0,
		)

		device.Display(matrix, 255, 10)

		time.Sleep(2 * time.Second)
	}
	<-done
}
