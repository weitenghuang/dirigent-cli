package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	// "github.com/weitenghuang/dirigent-cli/pkg/deploy"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create cluster resource yaml file",
	Long:  `create cluster resource yaml file from docker-compose file`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("%#v", *cmd)
		// err := deploy.Run("/opt/docker-compose.yml")
		// if err != nil {
		// 	log.Errorln(err)
		// }
	},
}
