package cmd

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list saved projects",
	Long:  `List all projects previously saved by this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize database.
		db := initDatabase()
		result := db.Prep()
		if result == true {
			log.Debug("Database initialized")
		}

		projects := db.GetProjects("")
		for _, project := range projects {
			fmt.Printf("%s (%s)\n", project.Name, project.Path)
		}
	},
}
