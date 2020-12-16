package evdev

import (
	"syscall"
	"time"
	"unsafe"
)

// https://python-evdev.readthedocs.io/en/latest/apidoc.html
const (
	KeyRelease = 0
	KeyPress   = 1
	KeyHold    = 2
)

const (
	// marker to separate events. Events may be separated in time or in space, such as with the multitouch protocol.
	EvSyn EventType = 0x00
	// state changes of keyboards, buttons, or other key-like devices.
	EvKey EventType = 0x01
	// relative axis value changes, e.g. moving the mouse 5 units to the left.
	EvRel EventType = 0x02
	// absolute axis value changes, e.g. describing the coordinates of a touch on a touchscreen.
	EvAbs EventType = 0x03
	// miscellaneous input data that do not fit into other types.
	EvMsc EventType = 0x04
	// binary state input switches.
	EvSw EventType = 0x05
	// turn LEDs on devices on and off.
	EvLed EventType = 0x11
	// output sound to devices.
	EvSnd EventType = 0x12
	// for autorepeating devices.
	EvRep EventType = 0x14
	// send force feedback commands to an input device.
	EvFf EventType = 0x15
	// special type for power button and switch input.
	EvPwr EventType = 0x16
	// receive force feedback device status.
	EvFfStatus EventType = 0x17
)

// EventType are groupings of codes under a logical input construct.
// Each type has a set of applicable codes to be used in generating events.
// See the Ev section for details on valid codes for each type
type EventType uint16

// eventsize is size of structure of InputEvent
var eventsize = int(unsafe.Sizeof(InputEvent{}))

// InputEvent is the keyboard event structure itself
type InputEvent struct {
	Time  syscall.Timeval
	Type  EventType
	Code  uint16
	Value int32 // 1=press, 0=release
}

func (i *InputEvent) TimevalToTime() time.Time {
	sec, nsec := i.Time.Unix()
	return time.Unix(sec, nsec)
}

// KeyString returns representation of pressed key as string
// eg enter, space, a, b, c...
func (i *InputEvent) KeyString() string {
	return keyCodeMap[i.Code]
}

// KeyPress is the value when we press the key on keyboard
func (i *InputEvent) KeyPress() bool {
	return i.Value == 1
}

// KeyRelease is the value when we release the key on keyboard
func (i *InputEvent) KeyRelease() bool {
	return i.Value == 0
}
