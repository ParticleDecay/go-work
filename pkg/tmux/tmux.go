package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// LaunchEnvironment creates a new tmux session or connects to an existing one.
func LaunchEnvironment(sessionName string, path string, goroot string, gopath string) {
	quotedName := fmt.Sprintf("'%s'", sessionName)
	log.Debugf("Calling out to tmux command")
	tmuxes, err := exec.Command("tmux", "ls", "-F", "'#S'").Output()
	if err != nil {
		// Check for exit status 1 (tmux server not running)
		ignoreErr := false
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 1 {
					ignoreErr = true // we want to start it ourselves
					log.Debugf("No tmux sessions found")
				}
			}
		}
		if ignoreErr == false {
			log.Fatal(err)
		}
	}
	for _, session := range strings.Split(string(tmuxes), "\n") {
		if string(session) == quotedName {
			log.Debugf("Found existing tmux session at '%s'", sessionName)
			command := exec.Command("tmux", "a", "-t", sessionName)
			command.Stdin = os.Stdin
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			err = command.Run()
			if err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}
	}

	// If we're still here, we should create a new tmux session.
	log.Debug("tmux session not found, creating new session")
	os.Chdir(path)
	createCmd := exec.Command("tmux", "new", "-s", sessionName, "-c", path)
	createCmd.Stdin = os.Stdin
	createCmd.Stdout = os.Stdout
	createCmd.Stderr = os.Stderr
	err = createCmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		gorootCmd := exec.Command("tmux", "set-environment", "-t", sessionName, "GOROOT", goroot)
		err = gorootCmd.Run()
		if err != nil {
			log.Errorf("There was an error setting GOROOT: %s", err)
		}
	}()
	go func() {
		gopathCmd := exec.Command("tmux", "set-environment", "-t", sessionName, "GOPATH", gopath)
		err = gopathCmd.Run()
		if err != nil {
			log.Errorf("There was an error setting GOPATH: %s", err)
		}
	}()
	os.Exit(0)
}
