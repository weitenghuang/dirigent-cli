package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "dirigent",
	Short: "Short description",
	Long:  `Long description ...`,
}
