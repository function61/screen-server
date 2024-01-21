package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/function61/gokit/app/evdev"
	"github.com/function61/gokit/log/logex"
	"github.com/function61/gokit/sync/taskrunner"
)

const (
	pidFilePath = "/run/nwkevdev/nwkevdev.pid"
)

func runClient(ctx context.Context) error {
	fmt.Println("launching client")

	conn, err := net.Dial("udp", "127.0.0.1:6666")
	if err != nil {
		return err
	}
	defer conn.Close()

	tick := time.NewTicker(1 * time.Second)

	on := true

	nextBlinkState := func() bool {
		defer func() {
			on = !on // blink
		}()

		return on
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
			binary.Write(conn, binary.LittleEndian, &evdev.InputEvent{
				Type: evdev.EvLed,
				Code: 0, // first led
				Value: func() int32 {
					if nextBlinkState() {
						return 0x01
					} else {
						return 0x00
					}
				}(),
			})

		}
	}
}

type config struct {
	endpointAddr string

	devicesToAlwaysGrab []string
}

func runClientGrab(ctx context.Context, conf config) error {
	log.Printf("launching client-grab -> %s", conf.endpointAddr)

	conn, err := net.Dial("udp", conf.endpointAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// this channel will receive input from all input devices that we'll start scanning
	allDevicesInput := evdev.NewChan()

	makeDevGrabberTask := func(devPath string) func(context.Context) error {
		return func(ctx context.Context) error {
			inputDev, inputClose, err := evdev.OpenWithChan(devPath, allDevicesInput)
			if err != nil {
				return fmt.Errorf("failed opening input device %s: %w", devPath, err)
			}
			defer func() {
				if err := inputClose(); err != nil {
					log.Printf("evdev inputClose: %v", err)
				}
			}()

			// grabbed = input will only be processed by us
			if err := inputDev.ScanInputGrabbed(ctx); err != nil {
				return err
			}

			return nil
		}
	}

	tasks := taskrunner.New(ctx, logex.StandardLogger())

	for idx, devPath := range conf.devicesToAlwaysGrab {
		log.Printf("grabbing %s", devPath)

		taskName := fmt.Sprintf("dev%d", idx)

		tasks.Start(taskName, makeDevGrabberTask(devPath))
	}

	tasksDone := tasks.Done()

	for {
		select {
		case err := <-tasksDone:
			return err
		case e := <-allDevicesInput:
			// just pipe the udev packet 1:1 over UDP
			if _, err := conn.Write(e.AsBytes()); err != nil {
				return err
			}
		}
	}
}

// determine if already running:
// - if true: stop existing
// - if false: start controlling remote
func runClientToggle(ctx context.Context) error {
	return toggleProcessFromPidFile(
		ctx,
		pidFilePath,
		exec.Command(os.Args[0], "client-grab"))
}

func toggleProcessFromPidFile(ctx context.Context,path string,cmd *exec.Cmd) error {
	pidStr, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) { // unexpected error
		return err
	}

	if os.IsNotExist(err) {
		log.Println("launching")

		if err := cmd.Start(); err != nil {
			return err
		}

		return nil
	} else {
		log.Println("asking to stop")

		pid, err := strconv.Atoi(string(pidStr))
		if err != nil {
			return err
		}
		proc, err := os.FindProcess(pid)
		if err != nil {
			return err
		}

		if err := proc.Signal(os.Interrupt); err != nil {
			return err
		}

		return nil
	}
}
