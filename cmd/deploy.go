package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weitenghuang/dirigent-cli/pkg/deploy"
)

func init() {
	RootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy resource",
	Long:  `deploy resource defined in yaml/json file to target cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deploy...")
		err := deploy.CreatePods("/opt/docker-compose.yml")
		if err != nil {
			fmt.Println(err)
		}
	},
}
