package deploy

import (
	// "fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/api"
	"io/ioutil"
	"os"
)

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

	log.Infoln(dockerCompose)

	rc := api.ReplicationController{ObjectMeta: api.ObjectMeta{Namespace: "other"}}

	log.Infoln(rc)

	return dockerCompose, nil
}
