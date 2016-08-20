package io

import (
	"log"
	"fmt"
	"github.com/paypal/gatt"
)

type Nuimo struct {
	deviceCallBack func(Device)
	ble_device gatt.Device
}

func ConnectNuimo(chan Nuimo) {

}



func (n Nuimo) Discover(f func(Device)) {

	var DefaultClientOptions = []gatt.Option{
		gatt.LnxMaxConnections(1),
		gatt.LnxDeviceID(-1, false),
	}

	n.deviceCallBack = f
	d, err := gatt.NewDevice(DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}
	n.ble_device = d
	// Register handlers.
	d.Handle(gatt.PeripheralDiscovered(n.onPeriphDiscovered))
	d.Init(n.onStateChanged)
	select {}
}

func (n Nuimo) onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("scanning...")
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		fmt.Println("stop scanning...")
		d.StopScanning()
	}
}

func (n Nuimo) onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if a.LocalName == "Nuimo" {
		n.deviceCallBack(Device{ble_device: n.ble_device, peripheral: p})
	}
}