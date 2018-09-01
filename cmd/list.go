package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var shortOutput bool

func init() {
	listCmd.PersistentFlags().BoolVarP(&shortOutput, "short", "s", false, "display shortened output")
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
			var extraOutput string
			if shortOutput != true {
				extraOutput = fmt.Sprintf(" (%s)", project.Path)
			}
			fmt.Printf("%s%s\n", project.Name, extraOutput)
		}
	},
}
