package git

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// URL represents a valid Git repository URL
type URL struct {
	Scheme     string
	Host       string
	Account    string
	Repository string
}

// FullURL returns a usable URL for the repository.
func (u *URL) FullURL() string {
	if u.Scheme == "" || u.Host == "" || u.Account == "" || u.Repository == "" {
		log.Error("This Git repository URL has one or more missing parts")
		return ""
	}
	schemeHost := fmt.Sprintf("%s%s", u.Scheme, u.Host)
	accountJoin := ":"
	if strings.HasPrefix(schemeHost, "http") {
		accountJoin = "/"
	}
	withoutRepo := strings.Join([]string{schemeHost, u.Account}, accountJoin)
	withRepo := fmt.Sprintf("%s/%s.git", withoutRepo, u.Repository)
	return withRepo
}

// ParseURL parses a Git repository URL and returns a URL object.
func ParseURL(url string) *URL {
	gitURL := regexp.MustCompile(`^([^@]+@|https?://|)([^:/]+)[:/]([^/]+)/([^\.]+)`)
	u := &URL{}
	groups := gitURL.FindStringSubmatch(url)
	if len(groups) != 5 || groups[2] == "" || groups[3] == "" || groups[4] == "" {
		return nil
	}
	if groups[1] == "" {
		u.Scheme = "git://"
	} else {
		u.Scheme = groups[1]
	}
	u.Host, u.Account, u.Repository = groups[2], groups[3], groups[4]
	return u
}

// GetRepositoryName attempts to return the name of the Git repository URL
func GetRepositoryName(url string) string {
	repoName := regexp.MustCompile(`([^/]+)(\.git$|$)`)
	groups := repoName.FindStringSubmatch(url)
	if len(groups) < 2 || groups[1] == "" {
		return ""
	}
	return strings.TrimSuffix(groups[1], ".git")
}
