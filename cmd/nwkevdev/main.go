package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/function61/gokit/os/osutil"
)

const (
	uinputDefaultPath = "/dev/uinput"
)

const (
	addrTestingScreenServer = "127.0.0.1:6666"
	addrWorklaptop          = "192.168.1.244:6666"
)

func main() {
	ctx := osutil.CancelOnInterruptOrTerminate(nil)

	switch os.Args[1] {
	case "server":
		osutil.ExitIfError(runServer(ctx))
	case "client":
		osutil.ExitIfError(runClient(ctx))
	case "client-grab":
		osutil.ExitIfError(withPIDFile(pidFilePath, func() error {
			return runClientGrab(ctx, config{
				// endpointAddr: addrTestingScreenServer,
				endpointAddr: addrWorklaptop,

				devicesToAlwaysGrab: []string{
					// "/dev/input/by-id/usb-SEM_USB_Keyboard-event-kbd",
					"/dev/input/by-id/usb-Logitech_USB_Receiver-if01-event-mouse",
					"/dev/input/by-id/usb-Massdrop_Inc._CTRL_Keyboard_1642645373-event-kbd",
				},
			})
		}))
	case "toggle":
		osutil.ExitIfError(runClientToggle(ctx))
	default:
		osutil.ExitIfError(fmt.Errorf("unknown verb: %s", os.Args[1]))
	}
}

func withPIDFile(path string, task func() error) error {
	pid := strconv.Itoa(os.Getpid())
	if err := os.WriteFile(path, []byte(pid), 0644); err != nil {
		return err
	}

	defer os.Remove(path) // not checking error

	return task()
}
