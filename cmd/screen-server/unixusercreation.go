package main

import (
	"log"
	"os/exec"

	"github.com/function61/gokit/os/osutil"
)

// each screen has its own Unix user. Firefox proved to be problematic if trying to be ran from two
// different X sessions with same user, and dicking around with Firefox profiles would've been harder
// than actually using separate users
func createUserIfNotExists(screen *Screen) error {
	homeDirExists, err := osutil.Exists(screen.Homedir())
	if err != nil {
		return err
	}

	if homeDirExists {
		return nil
	}

	// was is not exist => user does not exist => create
	log.Printf("setting up user %s", screen.Username())

	// unfortunately, syntax of "$ adduser" is different for Debian/Alpine
	isAlpine, err := osutil.Exists("/etc/alpine-release")
	if err != nil {
		return err
	}

	if isAlpine {
		return exec.Command(
			"adduser",
			"-G", "alpine",
			"-s", "/bin/sh",
			"-D", // don't assign a password
			screen.Username(),
		).Run()
	} else {
		return exec.Command(
			"adduser",
			"--disabled-password",
			"--disabled-login",
			screen.Username(),
		).Run()
	}
}
