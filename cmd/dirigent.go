package cmd

import (
	"github.com/spf13/cobra"
)

func Execute() error {
	AddCommands()
	return DirigentCmd.Execute()
}

var DirigentCmd = &cobra.Command{
	Use:   "dirigent",
	Short: "Short description",
	Long:  `Long description ...`,
}

func AddCommands() {
	DirigentCmd.AddCommand(deployCmd)
	DirigentCmd.AddCommand(versionCmd)
	DirigentCmd.AddCommand(createCmd)
}
