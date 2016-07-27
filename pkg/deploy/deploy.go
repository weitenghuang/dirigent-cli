package deploy

import (
	// "fmt"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"strings"
)

const (
	YamlExtension       string = ".yml"
	DefaultPodFile      string = "/opt/deploy/pod"
	DefaultServiceFile  string = "/opt/deploy/service"
	DefaultAPIVersion   string = "v1"
	DefaultK8sNamespace string = "default"
	DefaultSelectorKey  string = "name"
)

func Run(filename string) error {
	dockerCompose, err := DocodeDockerComposeYaml(filename)
	if err != nil {
		return err
	}

	dockerServices, ok := dockerCompose["services"]
	if !ok {
		return errors.New("Please check if docker-compose.yml includes \"services\".")
	}

	apps := dockerServices.(map[string]interface{})
	if len(apps) == 0 {
		return errors.New("Please check if docker-compose.yml \"services\" section has at least one service.")
	}

	for appKey, rawValue := range apps {
		appValue := rawValue.(map[string]interface{})
		log.Infoln("Pod Deployment Starts: ", appKey)
		if err := DeployPod(appKey, appValue); err != nil {
			log.Errorln("Error: Deploy Pod ", err, appKey, appValue)
		}
		log.Infoln("Service Deployment Starts: ", appKey)
		if err := DeployService(appKey, appValue); err != nil {
			log.Errorln("Error: Deploy Service ", err, appKey, appValue)
		}
		log.Infoln("Deployment Done: ", appKey)
	}
	return nil
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

func getPodSelectorLabel(appKey string, version string) string {
	return strings.Join([]string{appKey, "-", version, "-pod"}, "")
}
