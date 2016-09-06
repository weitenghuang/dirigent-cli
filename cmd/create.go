package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/weitenghuang/dirigent-cli/pkg/create"
	"github.com/weitenghuang/dirigent-cli/pkg/resource"
)

var (
	resourceType string
)

const (
	DefaultComposeYamlPath string = "/opt/docker-compose.yml"
)

func init() {
	createCmd.Flags().StringVarP(&resourceType, "type", "t", "", "cluster's resource type to create")
}

var createCmd = &cobra.Command{
	Use:   "create --type ",
	Short: "create cluster resource yaml file",
	Long:  `create cluster resource yaml file from docker-compose file`,
	Run: func(cmd *cobra.Command, args []string) {
		if flagChanged(cmd.Flags(), "type") {
			var err error
			switch resourceType {
			case "replicationcontroller", "replication-controller", "replicationController", "RC", "rc":
				err = create.ResourceFile(resource.ReplicationController, DefaultComposeYamlPath)
			case "service":
				err = create.ResourceFile(resource.Service, DefaultComposeYamlPath)
			case "job":
				err = create.ResourceFile(resource.Job, DefaultComposeYamlPath)
			case "deployment":
				err = create.ResourceFile(resource.Deployment, DefaultComposeYamlPath)
			default:
				log.Warnf("Invalid resource type: %v", resourceType)
				cmd.Help()
			}
			if err != nil {
				log.Errorln(err)
			}
		} else {
			cmd.Help()
		}
	},
}
