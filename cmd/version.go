package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of dirigent",
	Long:  `Print dirigent's current version number`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dirigent is on version v0.0.1")
	},
}
