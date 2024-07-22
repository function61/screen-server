package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/function61/gokit/app/dynversion"
	"github.com/function61/gokit/app/evdev"
	"github.com/function61/gokit/log/logex"
	"github.com/function61/gokit/os/osutil"
	"github.com/function61/gokit/sync/taskrunner"
	"github.com/spf13/cobra"
)

const (
	instructionsUrl = "https://github.com/function61/screen-server"
)

func main() {
	app := &cobra.Command{
		Use:     os.Args[0],
		Short:   "Serves your screens.",
		Version: dynversion.Version,
	}

	app.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Runs the server",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			rootLogger := logex.StandardLogger()

			osutil.ExitIfError(run(
				osutil.CancelOnInterruptOrTerminate(rootLogger),
				rootLogger))
		},
	})

	osutil.ExitIfError(app.Execute())
}

type screenOptions struct {
	Description       string
	width             int
	height            int
	vncPort           int
	AttachInputDevice string
}

type Screen struct {
	Id            int
	XScreenNumber int
	Opts          screenOptions
	Osd           OsdDriver
	xUtilConn     *xgbutil.XUtil
	xUtilConnMu   sync.Mutex
}

func (s *Screen) XScreenNumberWithColon() string {
	return ":0"
}

func (s *Screen) VncWebsocketPath() string {
	websocketPath := fmt.Sprintf("/api/screen/%d/ws", s.Id)

	// for some reason the path doesn't start with '/'
	return url.QueryEscape(websocketPath[1:])
}

func run(ctx context.Context, logger *log.Logger) error {
	logl := logex.Levels(logger)

	optss := []screenOptions{}

	for i := 1; ; i++ {
		key := fmt.Sprintf("SCREEN_%d", i)

		serialized := os.Getenv(key)
		if serialized == "" {
			break
		}

		opts, err := parseScreenOpts(serialized)
		if err != nil {
			return fmt.Errorf("parseScreenOpts: %s: %w", key, err)
		}

		optss = append(optss, *opts)
	}

	// var osdDriver OsdDriver = &firefoxOsdDriver{}
	var osdDriver OsdDriver = &zenityOsdDriver{}

	screens := []*Screen{}
	for idx, opts := range optss {
		screen := &Screen{
			Id:            idx + 1, // 1,2,3...
			XScreenNumber: idx + 1, // 1,2,3...
			Opts:          opts,
			Osd:           osdDriver,
		}

		screens = append(screens, screen)
	}

	if screens := len(screens); screens != 1 {
		return fmt.Errorf("expected exactly one screen; got %d", screens)
	}

	// each screen task encapsulates three processes: Xvfb, x11vnc and openbox
	screenTasks := taskrunner.New(ctx, logger)

	if len(screens) == 0 {
		return fmt.Errorf(
			"Define at least one screen (configuration error)\nInstructions: %s",
			instructionsUrl)
	}
	screen := screens[0]

	screenTasks.Start(screen.Opts.Description, func(ctx context.Context) error {
		return runOneScreen(ctx, screen, logl, logger)
	})

	screenTasks.Start("webui", func(ctx context.Context) error {
		handler := newServerHandler(screens, logex.Prefix("webui", logger))

		return runServer(ctx, handler)
	})

	return screenTasks.Wait()
}

func runOneScreen(
	ctx context.Context,
	screen *Screen,
	logl *logex.Leveled,
	logger *log.Logger,
) error {
	xvfbReady := make(chan struct{})

	go func() {
		// dirty solution to assume Xvfb ready
		time.Sleep(2 * time.Second)

		close(xvfbReady)
	}()

	processes := taskrunner.New(ctx, logger)

	if screen.Opts.AttachInputDevice != "" {
		exists, err := osutil.Exists(screen.Opts.AttachInputDevice)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf(
				"Input-device-to-attach '%s' doesn't exist",
				screen.Opts.AttachInputDevice)
		}

		// this channel will receive input from the input device that we'll start scanning
		input := evdev.NewChan()

		inputDev, inputClose, err := evdev.OpenWithChan(screen.Opts.AttachInputDevice, input)
		if err != nil {
			return fmt.Errorf("failed opening input device %s: %w", screen.Opts.AttachInputDevice, err)
		}
		defer func() {
			if err := inputClose(); err != nil {
				logl.Error.Printf("evdev inputClose: %v", err)
			}
		}()

		// grabbed = input will only be processed by us
		processes.Start("evdev", inputDev.ScanInputGrabbed)

		processes.Start("x11-input-forwarder", func(ctx context.Context) error {
			<-xvfbReady

			return deliverInputEventsToX(
				ctx,
				input,
				screen.XScreenNumberWithColon(),
				logex.Prefix("x11-input-forwarder", logger))
		})
	}

	processes.Start("Xvfb", func(ctx context.Context) error {
		// this serves as a virtual display
		//nolint:gosec // ok
		xvfb := exec.CommandContext(
			ctx,
			"Xvfb",
			screen.XScreenNumberWithColon(),
			"-screen", "0", fmt.Sprintf("%dx%dx24", screen.Opts.width, screen.Opts.height))

		return xvfb.Run()
	})

	processes.Start("x11vnc", func(ctx context.Context) error {
		<-xvfbReady

		// this serves the virtual display over VNC
		//nolint:gosec // ok
		x11vnc := exec.CommandContext(
			ctx,
			"x11vnc",
			"-display", screen.XScreenNumberWithColon(),
			"-rfbport", strconv.Itoa(screen.Opts.vncPort),
			"-forever", // without this the process exits after first disconnect, WTF why
			"-xkb",     // use XKEYBOARD extension (improved keymapping)
			"-noxrecord",
			"-noxfixes",
			// "-noxdamage", // TODO: why was this optimization turned off?
			"-nopw",                             // don't display warning of serving VNC without password
			"-desktop", screen.Opts.Description, // VNC viewer might show this (TightVNC on Windows does)
			"-wait", "5", // screen poll [ms]
			"-shared", // allow simultaneous connections to this same display
			"-permitfiletransfer",
			"-tightfilexfer",
		)

		if os.Getenv("DEBUG") != "" {
			x11vnc.Stderr = os.Stderr
		}

		return x11vnc.Run()
	})

	processes.Start("openbox", func(ctx context.Context) error {
		<-xvfbReady

		// this serves as a window manager so the screen has a menu where the user can start
		// Firefox and a terminal
		openbox := exec.CommandContext(ctx, "openbox", "--startup", "firefox")
		openbox.Env = append(
			os.Environ(), // make sure possible TZ gets passed
			"DISPLAY="+screen.XScreenNumberWithColon())

		if os.Getenv("DEBUG") != "" {
			openbox.Stdout = os.Stdout
			openbox.Stderr = os.Stderr
		}

		return openbox.Run()
	})

	return processes.Wait()
}

var screenOptsParseRe = regexp.MustCompile("^([^,]+),([^,]+),([^,]+),([^,]+)(?:,([^,]+))?$")

func parseScreenOpts(serialized string) (*screenOptions, error) {
	screenDefParts := screenOptsParseRe.FindStringSubmatch(serialized)
	if screenDefParts == nil {
		return nil, errors.New("does not match format VNCPORT,WIDTH,HEIGHT,NAME,[INPUT_DEVICE]")
	}

	vncPort, err := strconv.Atoi(screenDefParts[1])
	if err != nil {
		return nil, fmt.Errorf("vncPort: %w", err)
	}

	width, err := strconv.Atoi(screenDefParts[2])
	if err != nil {
		return nil, fmt.Errorf("width: %w", err)
	}

	height, err := strconv.Atoi(screenDefParts[3])
	if err != nil {
		return nil, fmt.Errorf("height: %w", err)
	}

	return &screenOptions{
		vncPort:           vncPort,
		width:             width,
		height:            height,
		Description:       screenDefParts[4],
		AttachInputDevice: screenDefParts[5],
	}, nil
}
