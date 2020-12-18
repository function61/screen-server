// Translates evdev input codes into "fake input" for X11's XTEST extension
package evdevtoxtesttranslator

import (
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

const (
	MotionNotifyAbsolute = 0
	MotionNotifyRelative = 1
)

type FakeInput struct {
	// these are arguments to xtest.FakeInput()
	Type          byte  // ButtonPress | ButtonRelease | KeyPress | KeyRelease | ...
	Detail        byte  // <key code> | <button code> | <mouse button code> | MotionNotifyAbsolute | MotionNotifyRelative | ...
	MotionNotifyX int16 // when Type=MotionNotify
	MotionNotifyY int16 // when Type=MotionNotify

	RepeatOnceForType byte
	repeat            int // our own construct. empty repeat value (=0) means 1
}

func (f FakeInput) Repeat() int {
	if f.repeat == 0 {
		return 1
	} else {
		return f.repeat
	}
}

// for some reason these show up as Key events, while mouse scroll and movement are Rel events
// evdev code => xtest code
var mouseBtnKeyCodes = map[uint16]byte{
	// some values listed here https://who-t.blogspot.com/2016/09/understanding-evdev.html
	272: mouseLeft,
	273: mouseRight,
	274: mouseMiddle, // usually scroll click
}

// Translate evdev input event into XTEST fake input
// these mappings rely on X11's keyboard layout to be evdev
func Translate(e evdev.InputEvent, logl *logex.Leveled) *FakeInput {
	switch e.Type {
	case evdev.EvSyn: // "transaction boundary" - used to group together related events.
		return nil // NOOP
	case evdev.EvKey:
		// keyboard key codes and mouse button key codes are conceptually in the same group in evdev,
		// but XTEST treats them differently (keys vs. buttons)
		if btnXtestCode, isMouseBtnKeyCode := mouseBtnKeyCodes[e.Code]; isMouseBtnKeyCode {
			return handleButtonEvent(e, btnXtestCode, logl)
		} else {
			return handleKeyEvent(e, logl)
		}
	case evdev.EvRel:
		return handleRelativePointerMovement(e, logl)
	case evdev.EvMsc: // miscellaneous - usually used to report raw scancodes
		return nil // NOOP
	default:
		logl.Debug.Printf("unknown InputEvent type: %d", e.Type)
		return nil
	}
}

func handleKeyEvent(e evdev.InputEvent, logl *logex.Leveled) *FakeInput {
	switch e.Value {
	case evdev.KeyPress:
		// +8 because of course https://stackoverflow.com/a/53551666/2176740
		return &FakeInput{Type: KeyPress, Detail: byte(e.Code + 8)}
	case evdev.KeyRelease:
		return &FakeInput{Type: KeyRelease, Detail: byte(e.Code + 8)}
	case evdev.KeyHold:
		// ignore holds (we don't need it in this context)
		return nil
	default:
		logl.Debug.Printf("unknown InputEvent value: %d", e.Value)
		return nil
	}
}

func handleButtonEvent(e evdev.InputEvent, btnXtestCode byte, logl *logex.Leveled) *FakeInput {
	switch e.Value {
	case evdev.KeyPress:
		return &FakeInput{Type: ButtonPress, Detail: btnXtestCode}
	case evdev.KeyRelease:
		return &FakeInput{Type: ButtonRelease, Detail: btnXtestCode}
	case evdev.KeyHold:
		// ignore holds (we don't need it in this context)
		return nil
	default:
		logl.Debug.Printf("unknown InputEvent value: %d", e.Value)
		return nil
	}
}

func handleRelativePointerMovement(e evdev.InputEvent, logl *logex.Leveled) *FakeInput {
	switch e.Code {
	case 0: // REL_X
		return &FakeInput{Type: MotionNotify, Detail: MotionNotifyRelative, MotionNotifyX: int16(e.Value)}
	case 1: // REL_Y
		return &FakeInput{Type: MotionNotify, Detail: MotionNotifyRelative, MotionNotifyY: int16(e.Value)}
	case 6: // REL_HWHEEL
		logl.Debug.Println("ignoring REL_HWHEEL")
		return nil
	case 8: // REL_WHEEL (vertical wheel, "scroll")
		// 2 = more up
		// 1 = little up
		// -1 = little down
		// -2 = more down
		// ...
		scrollAmount, up := func() (int, bool) { // decompose to absolute & direction
			if e.Value < 0 {
				return int(-e.Value), false
			} else {
				return int(e.Value), true
			}
		}()

		if up {
			return &FakeInput{Type: ButtonPress, Detail: mouseWheelUp, repeat: scrollAmount, RepeatOnceForType: ButtonRelease}
		} else {
			return &FakeInput{Type: ButtonPress, Detail: mouseWheelDown, repeat: scrollAmount, RepeatOnceForType: ButtonRelease}
		}
	case 11, 12: // REL_WHEEL_HI_RES, REL_HWHEEL_HI_RES
		// ignore (scroll events seem to be delivered also)
		return nil
	default:
		logl.Error.Printf("unknown code for relative movement: %d", e.Code)
		return nil
	}
}