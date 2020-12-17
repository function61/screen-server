package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xtest"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/function61/gokit/log/logex"
	"github.com/function61/screen-server/pkg/evdev"
)

// taken from https://sourcegraph.com/github.com/AdoptOpenJDK/openjdk-jdk11/-/blob/src/java.desktop/unix/classes/sun/awt/X11/XConstants.java#L88
// see also https://github.com/whot/libevdev/blob/master/include/linux/input-event-codes.h
// TODO: reference these from somewhere (also fix REL_X, REL_WHEEL etc. referenced from code below)
const (
	KeyPress      = 2
	KeyRelease    = 3
	ButtonPress   = 4
	ButtonRelease = 5
	MotionNotify  = 6

	// KeyPressMask   = 1 << 0
	// KeyReleaseMask = 1 << 1

	mouseLeft      = 1
	mouseMiddle    = 2
	mouseRight     = 3
	mouseWheelUp   = 4
	mouseWheelDown = 5
)

// for some reason these show up as Key events, while mouse scroll and movement are Rel events
// evdev code => xtest code
var mouseBtnKeyCodes = map[uint16]byte{
	// some values listed here https://who-t.blogspot.com/2016/09/understanding-evdev.html
	272: mouseLeft,
	273: mouseRight,
	274: mouseMiddle, // usually scroll click
}

func deliverInputEventsToX(
	ctx context.Context,
	input chan evdev.InputEvent,
	display string,
	logger *log.Logger,
) error {
	logl := logex.Levels(logger)

	// FIXME: not strictly safe/semantic at all if this function is called multiple times
	// FIXME: due to this we get "A read error is unrecoverable: EOF" that seem to be originating from
	//        this code, even though it is not. Global state..
	xgb.Logger = logger

	// TODO: use screen.getXConn(), but we've to implement centralized per-screen X connection
	//       teardown logic (+ maybe xevent.Main?) first
	xUtil, err := xgbutil.NewConnDisplay(display)
	if err != nil {
		return err
	}
	xConn := xUtil.Conn()
	// both close types (a) xConn.Close) b) xevent.Quit() ) yield panics :(

	// not sure if this is required
	go xevent.Main(xUtil)

	if err := xtest.Init(xConn); err != nil {
		return fmt.Errorf("xtest.Init: %w", err)
	}

	// TODO: can this change
	rootWin := xUtil.RootWin()

	handleKeyEvent := func(e evdev.InputEvent) {
		// for some reason mouse buttons are delivered as key codes
		if btnXtestCode, isMouseBtnKeyCode := mouseBtnKeyCodes[e.Code]; isMouseBtnKeyCode {
			switch e.Value {
			case evdev.KeyPress:
				xtest.FakeInput(xConn, ButtonPress, btnXtestCode, 0, rootWin, 0, 0, 0)
			case evdev.KeyRelease:
				xtest.FakeInput(xConn, ButtonRelease, btnXtestCode, 0, rootWin, 0, 0, 0)
			case evdev.KeyHold:
				// ignore holds (we don't need it in this context)
			default:
				logl.Debug.Printf("unknown InputEvent value: %d", e.Value)
			}
		} else {
			switch e.Value {
			case evdev.KeyPress:
				// +8 because of course https://stackoverflow.com/a/53551666/2176740
				xtest.FakeInput(xConn, KeyPress, byte(e.Code+8), 0, rootWin, 0, 0, 0)
			case evdev.KeyRelease:
				xtest.FakeInput(xConn, KeyRelease, byte(e.Code+8), 0, rootWin, 0, 0, 0)
			case evdev.KeyHold:
				// ignore holds (we don't need it in this context)
			default:
				logl.Debug.Printf("unknown InputEvent value: %d", e.Value)
			}
		}
	}

	handleRelativePointerMovement := func(e evdev.InputEvent) {
		amount := int16(e.Value)

		// with MotionNotify
		//   detail=0 => absolute
		//   detail=1 => relative
		relativeMovement := byte(1)

		mouseButtonPressAndRelease := func(btn byte) {
			xtest.FakeInput(xConn, ButtonPress, btn, 0, rootWin, 0, 0, 0)
			xtest.FakeInput(xConn, ButtonRelease, btn, 0, rootWin, 0, 0, 0)
		}

		switch e.Code {
		case 0: // REL_X
			xtest.FakeInput(xConn, MotionNotify, relativeMovement, 0, rootWin, amount, 0, 0)
		case 1: // REL_Y
			xtest.FakeInput(xConn, MotionNotify, relativeMovement, 0, rootWin, 0, amount, 0)
		case 6: // REL_HWHEEL
			logl.Debug.Println("ignoring REL_HWHEEL")
		case 8: // REL_WHEEL (vertical wheel, "scroll")
			// 2 = more up
			// 1 = little up
			// -1 = little down
			// -2 = more down
			// ...
			scrollAmount, up := func() (int32, bool) { // decompose to absolute & direction
				if e.Value < 0 {
					return -e.Value, false
				} else {
					return e.Value, true
				}
			}()

			for i := int32(0); i < scrollAmount; i++ {
				if up {
					mouseButtonPressAndRelease(mouseWheelUp)
				} else {
					mouseButtonPressAndRelease(mouseWheelDown)
				}
			}
		case 11, 12: // REL_WHEEL_HI_RES, REL_HWHEEL_HI_RES
			// ignore (scroll events seem to be delivered also)
		default:
			logl.Error.Printf("unknown code for relative movement: %d", e.Code)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-input:
			// these mappings rely on X11's keyboard layout to be evdev
			switch e.Type {
			case evdev.EvSyn: // "transaction boundary"
				// used to group together related events. we don't use it, so NOOP
			case evdev.EvKey:
				handleKeyEvent(e)
			case evdev.EvRel:
				handleRelativePointerMovement(e)
			case evdev.EvMsc: // miscellaneous - usually used to report raw scancodes
				// NOOP
			default:
				logl.Debug.Printf("unknown InputEvent type: %d", e.Type)
			}
		}
	}
}
