package io

import (
	"log"
	"github.com/paypal/gatt"
)

type Device struct {
	ble_device      gatt.Device
	peripheral      gatt.Peripheral

	connectCallback func(err error)
}

func (d *Device) ID() string {
	return d.peripheral.ID()
}

func (d *Device) Connect(f func(error)) {
	d.connectCallback = f
	d.ble_device.StopScanning()
	d.ble_device.Handle(
		gatt.PeripheralConnected(d.onConnected),
		gatt.PeripheralDisconnected(d.onDisconnected),
	)

	d.peripheral.Device().Connect(d.peripheral)
}

func (d *Device) onConnected(p gatt.Peripheral, err error) {
	log.Println("Connected...");

	if err := p.SetMTU(500); err != nil {
		log.Printf("Failed to set MTU, err: %s\n", err)
	} else {
		log.Printf("MTU set\n")
	}
	if d.connectCallback != nil {
		d.connectCallback(err)
	}
/*
	log.Printf("Services...")
	ss, err := d.peripheral.DiscoverServices(nil)
	log.Printf("Services %d", len(ss))


	if err != nil {
		log.Printf("Failed to discover services, err: %s\n", err)
		return
	} else {
		log.Printf("Discover service\n")
	}

	for _, s := range ss {
		msg := "Service: " + s.UUID().String()
		if len(s.Name()) > 0 {
			msg += " (" + s.Name() + ")"
		}
		log.Printf(msg)
	}*/

}
func (d *Device) onDisconnected(p gatt.Peripheral, err error) {
	d.peripheral.Device().CancelConnection(p)
	log.Println("Disconnected " + d.peripheral.ID())
	d.peripheral.Device().Connect(d.peripheral)
}