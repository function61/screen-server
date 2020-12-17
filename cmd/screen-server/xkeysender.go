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
	"github.com/function61/screen-server/pkg/evdevtoxtesttranslator"
)

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

	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-input:
			// Xvfb doesn't support itself reading input from evdev (like X.org does), so we've to
			// write the piping code ourselves. Easiest is to use XTEST "fake input", since with that
			// we don't have to keep track of active window to know which window to send input to.
			//
			// For more info: https://joonas.fi/2020/12/attach-a-keyboard-to-a-docker-container/

			// translate evdev input codes into XTEST fake input codes. basically the numbers are
			// somewhat different.
			fakeInput := evdevtoxtesttranslator.Translate(e, logl)

			// some evdev events we just ignore or don't have a meaningful XTEST fake input for
			if fakeInput == nil {
				break
			}

			// some (rare) events need to be "mechanically" repeated many times. e.g. evdev tells
			// "mouse wheel scrolls down fast" => we don't have a way to do that with XTEST so we
			// translate that to many single events like "scroll down, scroll down, ..."
			for i := 0; i < fakeInput.Repeat(); i++ {
				// TODO: check error, while still not adding to latency? how does sync/async stuff
				//       work with X11 / xgb library?
				xtest.FakeInput(xConn, fakeInput.Type, fakeInput.Detail, 0, rootWin, fakeInput.MotionNotifyX, fakeInput.MotionNotifyY, 0)

				// some (rare) events like mouse "scroll down" comes from evdev as a single event, but
				// X expects "scroll down button press, scroll down button release", so we synthetize it here
				// TODO: is this really requied?
				if fakeInput.RepeatOnceForType != 0 {
					xtest.FakeInput(xConn, fakeInput.RepeatOnceForType, fakeInput.Detail, 0, rootWin, fakeInput.MotionNotifyX, fakeInput.MotionNotifyY, 0)
				}
			}
		}
	}
}
