package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ParticleDecay/go-work/pkg/database"
	"github.com/ParticleDecay/go-work/pkg/git"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var installDir string

func init() {
	addCmd.PersistentFlags().StringVarP(&installDir, "dir", "d", "", "target parent directory for new project (default \"$GOPATH/src\")")
	rootCmd.AddCommand(addCmd)
}

func addHandler(cmd *cobra.Command, args []string) {
	// Initialize database.
	db := initDatabase()
	result := db.Prep()
	if result == true {
		log.Debug("Database initialized")
	}

	// Make sure the project dir exists.
	name := args[0]
	// If URL provided, use that to determine full path.
	var fullPath string
	var gitRemote *git.URL
	gitRemote = git.ParseURL(name)
	projectName := git.GetRepositoryName(name)
	if installDir == "" && gitRemote != nil { // default install into GOPATH
		fullPath = filepath.Join(gopath, "src", gitRemote.Host, gitRemote.Account, gitRemote.Repository)
	} else {
		// We don't have enough information about package, so assume name is the package
		if installDir == "" {
			installDir = filepath.Join(gopath, "src")
		}
		fullPath = filepath.Join(installDir, projectName)
	}

	// Check for existence of project.
	var projects []*database.Project
	projectResults := db.GetProjects(projectName)
	for _, project := range projectResults {
		// Only grab matching projects.
		if project.Path == fullPath {
			projects = append(projects, project)
		}
	}
	if len(projects) == 0 {
		dirExists := true
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			os.MkdirAll(fullPath, 0755)
			dirExists = false
		}
		result := db.NewProject(projectName, fullPath, goroot, gopath)
		if result == true {
			fmt.Printf("Project added. Launch it with the 'on' command.\n")
		}

		// Clone the repo.
		if gitRemote != nil && dirExists == false {
			log.Debug(fmt.Sprintf("Using git remote '%s'", gitRemote.FullURL()))
			repo := git.Repo{URL: gitRemote}
			repo.Clone(fullPath)
		}
	} else {
		fmt.Printf("'%s' project already exists at %s.\n", projectName, fullPath)
		os.Exit(1)
	}
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new project environment",
	Long:  `Adds a new project environment if one does not exist.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Error: A URL or a project name is required")
		}
		return nil
	},
	Run: addHandler,
}
