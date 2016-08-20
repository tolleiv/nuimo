package main

import (
	"log"
	"github.com/tolleiv/nuimo-go/io"
)

func main() {

	ch := make(chan io.Nuimo)
	io.ConnectNuimo(ch)
	nuimo := <- ch


	/*
	n := new(io.Nuimo)
	n.Discover(func(d io.Device) {

		d.Connect(func(err error) {
			if err == nil {
				log.Println("Connected " + d.ID());
			}
		})

		log.Println("Found " + d.ID())

	})*/
}
