package deploy

import (
	log "github.com/Sirupsen/logrus"
	"github.com/weitenghuang/dirigent-cli/pkg/utils"
)

const (
	YamlExtension       string = ".yml"
	DefaultPodFile      string = "/opt/deploy/pod"
	DefaultRcFile       string = "/opt/deploy/rc"
	DefaultServiceFile  string = "/opt/deploy/service"
	DefaultAPIVersion   string = "v1"
	DefaultK8sNamespace string = "default"
	DefaultSelectorKey  string = "name"
)

func Run(filename string) error {
	composeObject, err := utils.ParseDockerCompose(filename)
	if err != nil {
		return err
	}

	composeServiceNames := composeObject.ServiceConfigs.Keys()
	for _, name := range composeServiceNames {
		if composeServiceConfig, ok := composeObject.ServiceConfigs.Get(name); ok {
			log.Infoln("Service Deployment Starts: ", name)
			if err := Service(name, composeServiceConfig); err != nil {
				log.Errorln("Error: Deploy Service ", err, name, composeServiceConfig)
				return err
			}
			log.Infoln("Replication Controller Deployment Starts: ", name)
			if err := ReplicationController(name, composeServiceConfig); err != nil {
				log.Errorln("Error: Deploy ReplicationController ", err, name, composeServiceConfig)
				return err
			}
			log.Infoln("Deployment Done: ", name)
		}
	}

	return nil
}
