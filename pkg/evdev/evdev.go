// Read input events from evdev (Linux input device subsystem)
package evdev

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// https://github.com/gvalkov/golang-evdev implements what we do, but requires Cgo

/* Hunting the numeric code for EVIOCGRAB was an absolute pain (none of these specify it legibly):

- https://github.com/gvalkov/golang-evdev/search?q=EVIOCGRAB&unscoped_q=EVIOCGRAB
- https://github.com/pfirpfel/node-exclusive-keyboard/blob/master/lib/eviocgrab.cpp
- https://docs.rs/ioctl/0.3.1/src/ioctl/.cargo/registry/src/github.com-1ecc6299db9ec823/ioctl-0.3.1/src/platform/linux.rs.html#396
- https://github.com/torvalds/linux/blob/4a3033ef6e6bb4c566bd1d556de69b494d76976c/include/uapi/linux/input.h#L183
*/
const (
	EVIOCGRAB = 1074021776 // found out with help of evdev.EVIOCGRAB (github.com/gvalkov/golang-evdev)
)

// Sidenote: code for enumerating input devices: https://github.com/MarinX/keylogger/blob/master/keylogger.go

func NewChan() chan InputEvent {
	return make(chan InputEvent)
}

// ScanInput() may close the given channel
func ScanInput(ctx context.Context, inputDevicePath string, ch chan InputEvent) error {
	return scanInput(ctx, inputDevicePath, ch, false)
}

// same as ScanInput() but grabbed means exclusive access (we'll be the only one receiving the events)
// https://stackoverflow.com/a/1698686
// https://stackoverflow.com/a/1550320
func ScanInputGrabbed(ctx context.Context, inputDevicePath string, ch chan InputEvent) error {
	return scanInput(ctx, inputDevicePath, ch, true)
}

func scanInput(ctx context.Context, inputDevicePath string, ch chan InputEvent, grab bool) error {
	inputDevice, err := os.Open(inputDevicePath)
	if err != nil {
		return err
	}
	defer inputDevice.Close()

	// other programs won't receive input while we have this file handle open
	if err := grabExclusiveInputDeviceAccess(inputDevice); err != nil {
		return err
	}

	readingErrored := make(chan error, 1)

	go func() {
		for {
			e, err := readOneInputEvent(inputDevice)
			if err != nil {
				readingErrored <- fmt.Errorf("readOneInputEvent: %w", err)
				close(ch)
				break
			}

			if e != nil {
				ch <- *e
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			// this triggers close of inputDevice which will result in return of readOneInputEvent()
			// whose error will be returned to readingErrored, but that's ok because it's buffered
			// and will now not be read because we just exited gracefully
			return nil
		case err := <-readingErrored:
			return err
		}
	}
}

func isRoot() bool {
	return syscall.Getuid() == 0 && syscall.Geteuid() == 0
}

func readOneInputEvent(inputDevice *os.File) (*InputEvent, error) {
	buffer := make([]byte, eventsize)
	n, err := inputDevice.Read(buffer)
	if err != nil {
		return nil, err
	}
	// no input, dont send error
	if n <= 0 {
		return nil, nil
	}
	return eventFromBuffer(buffer)
}

func grabExclusiveInputDeviceAccess(inputDevice *os.File) error {
	if err := unix.IoctlSetInt(int(inputDevice.Fd()), EVIOCGRAB, 1); err != nil {
		return fmt.Errorf("grabExclusiveInputDeviceAccess: IOCTL(EVIOCGRAB): %w", err)
	}

	return nil
}

func eventFromBuffer(buffer []byte) (*InputEvent, error) {
	event := &InputEvent{}
	err := binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, event)
	return event, err
}

func timevalToTime(tv syscall.Timeval) time.Time {
	sec, nsec := tv.Unix()
	return time.Unix(sec, nsec)
}
