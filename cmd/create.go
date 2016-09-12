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
			switch resourceType {
			case "replicationcontroller", "replication-controller", "replicationController", "RC", "rc":
				err := create.ResourceFile(resource.ReplicationController, DefaultComposeYamlPath)
				if err != nil {
					log.Errorln(err)
				}
			case "service":
				err := create.ResourceFile(resource.Service, DefaultComposeYamlPath)
				if err != nil {
					log.Errorln(err)
				}
			case "job":
				err := create.ResourceFile(resource.Job, DefaultComposeYamlPath)
				if err != nil {
					log.Errorln(err)
				}
			case "deployment":
				err := create.ResourceFile(resource.Deployment, DefaultComposeYamlPath)
				if err != nil {
					log.Errorln(err)
				}
			case "all":
				resources := []resource.ResourceType{resource.Service, resource.Job, resource.Deployment}
				for _, value := range resources {
					err := create.ResourceFile(value, DefaultComposeYamlPath)
					if err != nil {
						log.Errorln(err)
					}
				}
			default:
				log.Warnf("Invalid resource type: %v", resourceType)
				cmd.Help()
			}
		} else {
			cmd.Help()
		}
	},
}
