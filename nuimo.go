package nuimo

import (
	"log"
	"fmt"
	"strings"
	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/gatt"
	"github.com/currantlabs/ble/linux/hci"
	"github.com/currantlabs/ble/linux/hci/cmd"
	"encoding/binary"
)

const SERVICE_BATTERY_STATUS = "0000180F00001000800000805F9B34FB"
const SERVICE_DEVICE_INFO = "0000180A00001000800000805F9B34FB"
const SERVICE_LED_MATRIX = "F29B1523CB1940F3BE5C7241ECB82FD1"
const SERVICE_USER_INPUT = "F29B1525CB1940F3BE5C7241ECB82FD2"

const CHAR_BATTERY_LEVEL = "00002A1900001000800000805F9B34FB"
const CHAR_DEVICE_INFO = "00002A2900001000800000805F9B34FB"
const CHAR_LED_MATRIX = "F29B1524CB1940F3BE5C7241ECB82FD1"
const CHAR_INPUT_FLY = "F29B1526CB1940F3BE5C7241ECB82FD2"
const CHAR_INPUT_SWIPE = "F29B1527CB1940F3BE5C7241ECB82FD2"
const CHAR_INPUT_ROTATE = "F29B1528CB1940F3BE5C7241ECB82FD2"
const CHAR_INPUT_CLICK = "F29B1529CB1940F3BE5C7241ECB82FD2"

const DIR_LEFT = 0
const DIR_RIGHT = 1
const DIR_UP = 2
const DIR_BACKWARDS = 2
const DIR_DOWN = 3
const DIR_TOWARDS = 3
const DIR_UPDOWN = 4

const CLICK_DOWN = 1
const CLICK_UP = 0

type Nuimo struct {
	client ble.Client
	events chan Event
}

type Event struct {
	Key   string
	Value int64
	Raw   []byte
}

func Connect() (*Nuimo, error) {
	filter := func(a ble.Advertisement) bool {
		return strings.ToUpper(a.LocalName()) == "NUIMO"
	}

	// Set connection parameters. Only supported on Linux platform.
	d := gatt.DefaultDevice()
	if h, ok := d.(*hci.HCI); ok {
		if err := h.Option(hci.OptConnParams(
			cmd.LECreateConnection{
				LEScanInterval:        0x0004, // 0x0004 - 0x4000; N * 0.625 msec
				LEScanWindow:          0x0004, // 0x0004 - 0x4000; N * 0.625 msec
				InitiatorFilterPolicy: 0x00, // White list is not used
				PeerAddressType:       0x00, // Public Device Address
				PeerAddress:           [6]byte{}, //
				OwnAddressType:        0x00, // Public Device Address
				ConnIntervalMin:       0x0006, // 0x0006 - 0x0C80; N * 1.25 msec
				ConnIntervalMax:       0x0006, // 0x0006 - 0x0C80; N * 1.25 msec
				ConnLatency:           0x0000, // 0x0000 - 0x01F3; N * 1.25 msec
				SupervisionTimeout:    0x0048, // 0x000A - 0x0C80; N * 10 msec
				MinimumCELength:       0x0000, // 0x0000 - 0xFFFF; N * 0.625 msec
				MaximumCELength:       0x0000, // 0x0000 - 0xFFFF; N * 0.625 msec
			})); err != nil {
			log.Fatalf("can't set advertising param: %s", err)
		}
	}

	client, err := gatt.Discover(gatt.FilterFunc(filter))
	if err != nil {
		return nil, err
	}

	ch := make(chan Event)
	return &Nuimo{client: client, events: ch}, nil

}

func (n *Nuimo) Events() <-chan Event {
	if err := n.DiscoverServices(); err != nil {
		log.Fatalf("discover issue: %s", err)
	}

	return n.events
}

func (n *Nuimo) DiscoverServices() error {
	p, err := n.client.DiscoverProfile(true)
	if err != nil {
		return fmt.Errorf("can't discover services: %s\n", err)
	}

	for _, s := range p.Services {

		if s.UUID.Equal(ble.MustParse(SERVICE_BATTERY_STATUS)) {
			for _, c := range s.Characteristics {
				if (c.UUID.Equal(ble.MustParse(CHAR_BATTERY_LEVEL))) {
					n.client.Subscribe(c, false, n.battery)
				}
			}
		}

		if s.UUID.Equal(ble.MustParse(SERVICE_USER_INPUT)) {
			for _, c := range s.Characteristics {
				if (c.UUID.Equal(ble.MustParse(CHAR_INPUT_CLICK))) {
					n.client.Subscribe(c, false, n.click)
				}

				if (c.UUID.Equal(ble.MustParse(CHAR_INPUT_ROTATE))) {
					n.client.Subscribe(c, false, n.rotate)
				}

				if (c.UUID.Equal(ble.MustParse(CHAR_INPUT_SWIPE))) {
					n.client.Subscribe(c, false, n.swipe)
				}
				if (c.UUID.Equal(ble.MustParse(CHAR_INPUT_FLY))) {
					n.client.Subscribe(c, false, n.fly)
				}
			}
		}

		//fmt.Printf("Service: %s %s, Handle (0x%02X)\n", s.UUID.String(), ble.Name(s.UUID), s.Handle)
	}
	return nil
}

func (n *Nuimo) battery(req []byte) {
	uval, _ := binary.Uvarint(req)
	level := int64(uval)
	n.events <- Event{Key:"battery", Raw: req, Value: level}
}

func (n *Nuimo) click(req []byte) {
	uval, _ := binary.Uvarint(req)
	dir := int64(uval)
	switch dir {
	case CLICK_DOWN:
		n.events <- Event{Key:"press", Raw: req}
	case CLICK_UP:
		n.events <- Event{Key:"release", Raw: req}
	}
}

func (n *Nuimo) rotate(req []byte) {
	uval := binary.LittleEndian.Uint16(req)
	val := int64(int16(uval))
	n.events <- Event{Key:"rotate", Raw: req, Value: val}
}
func (n *Nuimo) swipe(req []byte) {
	uval, _ := binary.Uvarint(req)
	dir := int64(uval)
	n.events <- Event{Key:"swipe", Raw: req, Value: dir}

	switch dir {
	case DIR_LEFT:
		n.events <- Event{Key:"swipe_left", Raw: req}
	case DIR_RIGHT:
		n.events <- Event{Key:"swipe_right", Raw: req}
	case DIR_UP:
		n.events <- Event{Key:"swipe_up", Raw: req}
	case DIR_DOWN:
		n.events <- Event{Key:"swipe_down", Raw: req}
	}
}
func (n *Nuimo) fly(req []byte) {
	uval, _ := binary.Uvarint(req[0:1])
	dir := int(uval)
	uval, _ = binary.Uvarint(req[2:])
	distance := int64(uval)

	switch dir {
	case DIR_LEFT:
		n.events <- Event{Key:"fly_left", Raw: req, Value: distance}
	case DIR_RIGHT:
		n.events <- Event{Key:"fly_right", Raw: req, Value: distance}
	case DIR_BACKWARDS:
		n.events <- Event{Key:"fly_backwards", Raw: req, Value: distance}
	case DIR_TOWARDS:
		n.events <- Event{Key:"fly_towards", Raw: req, Value: distance}
	case DIR_UPDOWN:
		n.events <- Event{Key:"fly_updown", Raw: req, Value: distance}
	}
}

func (n *Nuimo) Disconnect() error {
	close(n.events)
	return n.client.CancelConnection()
}