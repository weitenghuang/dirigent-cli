package deploy

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"strings"
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
	composeObject, err := ParseDockerCompose(filename)
	if err != nil {
		return err
	}

	composeServiceNames := composeObject.ServiceConfigs.Keys()
	for _, name := range composeServiceNames {
		if composeServiceConfig, ok := composeObject.ServiceConfigs.Get(name); ok {
			log.Infoln("Replication Controller Deployment Starts: ", name)
			if err := DeployReplicationController(name, composeServiceConfig); err != nil {
				log.Errorln("Error: Deploy ReplicationController ", err, name, composeServiceConfig)
				return err
			}
			log.Infoln("Service Deployment Starts: ", name)
			if err := DeployService(name, composeServiceConfig); err != nil {
				log.Errorln("Error: Deploy Service ", err, name, composeServiceConfig)
				return err
			}
			log.Infoln("Deployment Done: ", name)
		}
	}

	return nil
}

func ParseDockerCompose(filename string) (*project.Project, error) {
	context := &docker.Context{}
	if filename == "" {
		filename = "docker-compose.yml"
	}
	context.ComposeFiles = []string{filename}
	composeObject := project.NewProject(&context.Context, nil, nil)
	err := composeObject.Parse()
	if err != nil {
		log.Fatalf("Failed to load compose file", err)
		return nil, err
	}
	log.Infof("Post-parsing: %#v\n", composeObject)
	return composeObject, nil
}

func DocodeDockerComposeYaml(filename string) (map[string]interface{}, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}

	dockerComposeYaml, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var dockerCompose map[string]interface{}
	if err := yaml.Unmarshal(dockerComposeYaml, &dockerCompose); err != nil {
		return nil, err
	}
	log.Infof("%#v\n", dockerCompose)
	return dockerCompose, nil
}

func getPodSelectorLabel(appName string, version string) string {
	return strings.Join([]string{appName, "-", version, "-pod"}, "")
}

func getRCSelectorLabel(appName string, version string) string {
	return strings.Join([]string{appName, "-", version, "-rc"}, "")
}
