package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/weitenghuang/dirigent-cli/pkg/deploy"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy resource",
	Long:  `deploy resource defined in yaml/json file to target cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		err := deploy.Run("/opt/docker-compose.yml")
		if err != nil {
			log.Errorln(err)
		}
	},
}
