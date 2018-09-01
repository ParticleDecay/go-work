package git

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// Repo is a Git repository and related actions.
type Repo struct {
	URL *URL
}

// Clone performs a Git clone of a remote repository.
func (r *Repo) Clone(tgtDir string) {
	if r.URL == nil {
		log.Fatal("You cannot clone a Git repository without a URL")
	}
	log.Debugf("Cloning repository %s to %s", r.URL.FullURL(), tgtDir)
	cmdline := exec.Command("git", "clone", r.URL.FullURL(), tgtDir)
	err := cmdline.Run()
	if err != nil {
		log.Fatal(err)
	}
}
