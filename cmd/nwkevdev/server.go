package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/bendahl/uinput"
	"github.com/function61/gokit/app/evdev"
)

func runServer(ctx context.Context) error {
	fmt.Println("launching server")

	packetListener, err := net.ListenPacket("udp", ":6666")
	if err != nil {
		return err
	}
	defer packetListener.Close()

	go func() {
		<-ctx.Done()
		packetListener.Close() // double close intentional
	}()

	keyboard, err := uinput.CreateKeyboard(uinputDefaultPath, []byte("emulated-kb"))
	if err != nil {
		return fmt.Errorf("uinput device creation failed, see udev tips @ https://github.com/bendahl/uinput\n\t%w", err)
	}
	defer keyboard.Close()

	mouse, err := uinput.CreateMouse(uinputDefaultPath, []byte("emulated-mouse"))
	if err != nil {
		return fmt.Errorf("uinput device creation failed, see udev tips @ https://github.com/bendahl/uinput\n\t%w", err)
	}
	defer mouse.Close()

	// evdev events are very small & statically sized
	buffer := make([]byte, 1024)

	handleMouse := func(event *evdev.InputEvent) bool {
		isLeft := event.Code == 272
		isRight := event.Code == 273
		isMouse := isLeft || isRight

		if !isMouse {
			return false
		}

		switch event.Value {
		case evdev.KeyPress:
			if isLeft {
				_ = mouse.LeftPress()
			} else {
				_ = mouse.RightPress()
			}
		case evdev.KeyRelease:
			if isLeft {
				_ = mouse.LeftRelease()
			} else {
				_ = mouse.RightRelease()
			}
		case evdev.KeyHold:
			log.Println("unhandled KeyHold")
		default:
			// panic("wat")
			log.Printf("unhandled event code: %d", event.Code)
		}

		return true
	}

	for {
		n, _, err := packetListener.ReadFrom(buffer)
		if err != nil {
			return err
		}

		event, err := evdev.InputEventFromBytes(buffer[:n])
		if err != nil {
			return err
		}

		switch event.Type {
		case evdev.EvKey: // keyboard input
			if handleMouse(event) {
				break
			}

			switch event.Value {
			case evdev.KeyRelease:
				if err := keyboard.KeyUp(int(event.Code)); err != nil {
					return err
				}
			case evdev.KeyPress:
				if err := keyboard.KeyDown(int(event.Code)); err != nil {
					return err
				}
			case evdev.KeyHold:
				log.Println("unhandled KeyHold")
			default:
				// panic("wat")
				log.Printf("unhandled event code: %d", event.Code)
			}
		case evdev.EvRel: // relative mouse input
			switch evdev.Rel(event.Code) {
			case evdev.RelX:
				if err := mouse.Move(event.Value, 0); err != nil {
					return err
				}
			case evdev.RelY:
				if err := mouse.Move(0, event.Value); err != nil {
					return err
				}
			case evdev.RelWHEEL:
				if err:=mouse.Wheel(false,event.Value);err!=nil{
					return err
				}
			case evdev.RelWHEELHIRES:
				// hi-resolution extension to RelWHEEL (currently ignored)
			default:
				log.Printf("unknown relative code: %s", evdev.Rel(event.Code))
			}
		case evdev.EvSyn, evdev.EvMsc:
			// no-op, we don't need semantics of these with the uinput client library.
		default:
			log.Printf("unhandled event: %s", event)
		}
	}
}
