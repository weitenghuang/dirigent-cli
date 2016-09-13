package deploy

import (
	log "github.com/Sirupsen/logrus"
	"github.com/weitenghuang/dirigent-cli/pkg/resource"
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
		composeServiceConfig, ok := composeObject.ServiceConfigs.Get(name)
		notJob := resource.NotJobResource(name)
		if !ok {
			log.Debugf("Unable to Find %v's Compose Service Config\n", name)
			continue
		}
		if notJob {
			log.Infoln("Service Deployment Starts: ", name)
			if err := Service(name, composeServiceConfig); err != nil {
				log.Errorln("Error: Deploy Service ", err, name, composeServiceConfig)
				return err
			}
			log.Infoln("Deployment Resource Starts: ", name)
			if err := Deployment(name, composeServiceConfig); err != nil {
				log.Errorln("Error: Deploy Kubernetes Deployment ", err, name, composeServiceConfig)
				return err
			}
		} else {
			log.Infoln("Job Resource Starts: ", name)
			if err := Job(name, composeServiceConfig); err != nil {
				log.Errorln("Error: Deploy Kubernetes Job ", err, name, composeServiceConfig)
				return err
			}
		}
	}

	return nil
}
