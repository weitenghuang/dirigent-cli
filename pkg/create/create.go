package create

import (
	log "github.com/Sirupsen/logrus"
	"github.com/weitenghuang/dirigent-cli/pkg/resource"
	"github.com/weitenghuang/dirigent-cli/pkg/utils"
	"strings"
)

func ResourceFile(resourceType resource.ResourceType, composeFile string) error {
	composeObject, err := utils.ParseDockerCompose(composeFile)
	if err != nil {
		return err
	}

	composeServiceNames := composeObject.ServiceConfigs.Keys()
	for _, name := range composeServiceNames {
		if composeServiceConfig, ok := composeObject.ServiceConfigs.Get(name); ok {
			switch resourceType {
			case resource.ReplicationController:
				log.Infoln("Replication Controller File Creation Starts: ", name)
				if _, err := ReplicationController(name, composeServiceConfig); err != nil {
					log.Errorln("Error: ReplicationController File ", err, name, composeServiceConfig)
					return err
				}
			}
		}
	}

	return nil
}

func getPodSelectorLabel(appName string, version string) string {
	return strings.Join([]string{appName, "-", version, "-pod"}, "")
}

func getRCSelectorLabel(appName string, version string) string {
	return strings.Join([]string{appName, "-", version, "-rc"}, "")
}
