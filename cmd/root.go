package cmd

import (
	"errors"
	"os"
	"path/filepath"

	db "github.com/ParticleDecay/go-work/pkg/database"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	dbFile  string
	goroot  string
	gopath  string
	verbose bool
)

func init() {
	oldGoroot := os.Getenv("GOROOT")
	oldGopath := os.Getenv("GOPATH")
	rootCmd.PersistentFlags().StringVar(&dbFile, "dbfile", "", "database file (default \"$HOME/.gowork.db\")")
	rootCmd.PersistentFlags().StringVar(&goroot, "goroot", oldGoroot, "a custom GOROOT")
	rootCmd.PersistentFlags().StringVar(&gopath, "gopath", oldGopath, "a custom GOPATH")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "display debug messages")
}

func initDatabase() *db.Database {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	if dbFile == "" {
		dbFile = ".gowork.db"
	}

	return &db.Database{Path: filepath.Join(home, dbFile)}
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gowork",
	Short: "gowork creates tmux Go environments",
	Long:  `Creates or connects you to a tmux environment to work on Go projects.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose == true {
			log.SetLevel(log.DebugLevel)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Error: An action is required")
		}
		return nil
	},
}
