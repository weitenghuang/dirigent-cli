package cmd

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
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

func flagChanged(flags *flag.FlagSet, key string) bool {
	flag := flags.Lookup(key)
	if flag == nil {
		return false
	}
	return flag.Changed
}
