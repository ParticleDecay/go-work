package cmd

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workCmd)
}

var workCmd = &cobra.Command{
	Use:   "on",
	Short: "load the specified project environment",
	Long: `Connects to a specific Go environment or
			launches a new one.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Error: A project name is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize database.
		db := initDatabase()
		result := db.Prep()
		if result == true {
			log.Debug("Database initialized")
		}

		db.SelectProject(args[0])
	},
}
