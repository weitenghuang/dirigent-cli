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
		composeServiceConfig, ok := composeObject.ServiceConfigs.Get(name)
		notJob := notJobReource(name)
		if ok && notJob { // Regular resources
			switch resourceType {
			case resource.ReplicationController:
				log.Infoln("Replication Controller File Creation Starts: ", name)
				if _, err := ReplicationController(name, composeServiceConfig); err != nil {
					log.Errorln("Error: ReplicationController File ", err, name, composeServiceConfig)
					return err
				}
			case resource.Service:
				log.Infoln("Service File Creation Starts: ", name)
				if _, err := Service(name, composeServiceConfig); err != nil {
					log.Errorln("Error: Service File ", err, name, composeServiceConfig)
					return err
				}
			}
		} else if ok && resourceType == resource.Job && !notJob { // Special Resource
			log.Infof("%v\n", resourceType)
		}
	}

	return nil
}

func notJobReource(appName string) bool {
	return !strings.Contains(appName, "-job") && !strings.Contains(appName, "-init")
}
