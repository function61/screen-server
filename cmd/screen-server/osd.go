package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// OSD (On-Screen Display - supports showing messages on the screens) drivers

type OsdDriver interface {
	DisplayMessage(ctx context.Context, screen Screen, message string) error
}

func showOsdMessage(ctx context.Context, screen Screen, message string) {
	// default: show the message for <timeout> seconds
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logIfError("OSD DisplayMessage", screen.Osd.DisplayMessage(ctx, screen, message))
}

type zenityOsdDriver struct{}

func (f *zenityOsdDriver) DisplayMessage(ctx context.Context, screen Screen, message string) error {
	zenity := exec.CommandContext(
		ctx,
		"zenity",
		"--info",
		"--text", message)
	zenity.Env = append(zenity.Env, "DISPLAY="+screen.XScreenNumberWithColon())

	return zenity.Run()
}

//nolint:unused
type firefoxOsdDriver struct{}

func (f *firefoxOsdDriver) DisplayMessage(ctx context.Context, screen Screen, message string) error {
	html := f.makeHtml(string(message))

	// TODO: this is not thread safe
	osdHtmlFilename := "/tmp/osd.html"
	defer os.Remove(osdHtmlFilename)

	_ = ioutil.WriteFile(osdHtmlFilename, []byte(html), 0600)

	firefox := exec.CommandContext(
		ctx,
		"firefox",
		"-no-remote", // dunno why, https://stackoverflow.com/questions/26276293/open-firefox-in-fullscreen-from-command-line
		"-p", "default",
		"-width", strconv.Itoa(screen.Opts.width),
		"-height", strconv.Itoa(screen.Opts.height),
		"file://"+osdHtmlFilename,
	)
	firefox.Env = append(firefox.Env, "DISPLAY="+screen.XScreenNumberWithColon())

	return firefox.Run()
}

func (f *firefoxOsdDriver) makeHtml(message string) string {
	return fmt.Sprintf(`<html>
<head>
	<title></title>
	<script>
		setTimeout(() => {
			window.close();
		}, 5000);
	</script>
	<style>
	body {
		background: #000000;
		color: #ffffff;
	}
	</style>
</head>

<body>
	%s
</body>
</html>
`, message)
}
