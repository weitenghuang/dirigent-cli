package deploy

import (
	// "fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/api"
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

func CreatePods(filename string) error {
	dockerCompose, err := DocodeDockerComposeYaml(filename)
	if err != nil {
		return err
	}
	for serviceKey, serviceValue := range dockerCompose["services"].(map[string]interface{}) {
		pod := BuildPod(serviceKey, serviceValue.(map[string]interface{}))
		podYaml, err := yaml.Marshal(pod)
		if err != nil {
			return err
		}
		podYamlFile := "/opt/deploy/pod.yml"
		if err := ioutil.WriteFile(podYamlFile, podYaml, 0644); err != nil {
			return err
		}

		kubectlCreateCmd := exec.Command("kubectl", "create", "-f", podYamlFile)
		kubectlCreateCmd.Stdout = os.Stdout
		kubectlCreateCmd.Stderr = os.Stderr
		err = kubectlCreateCmd.Start()
		if err != nil {
			return err
		}
		err = kubectlCreateCmd.Wait()
		if err != nil {
			return err
		}
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

func BuildPod(serviceKey string, serviceValue map[string]interface{}) api.Pod {
	// log.Infoln("Key:", serviceKey, "Value:", serviceValue)
	var cVolumes, cPorts []string
	var cImage string
	for configKey, configValues := range serviceValue {
		// log.Infoln("ConfigKey:", configKey, "ConfigValues:", configValues)

		if cValue, ok := configValues.(string); ok && configKey == "image" {
			cImage = cValue
		} else if configKey == "volumes" && reflect.TypeOf(configValues).Kind() == reflect.Slice {
			for _, cValue := range configValues.([]interface{}) {
				cVolumes = append(cVolumes, cValue.(string))
			}
		} else if configKey == "ports" && reflect.TypeOf(configValues).Kind() == reflect.Slice {
			for _, cValue := range configValues.([]interface{}) {
				cPorts = append(cPorts, cValue.(string))
				log.Infoln(cPorts)
			}
		}
	}

	pod := api.Pod{ // Basic fields.
		TypeMeta:   unversioned.TypeMeta{Kind: "Pod", APIVersion: "v1"},
		ObjectMeta: api.ObjectMeta{Name: serviceKey, Namespace: "default"},
		Spec: api.PodSpec{
			//TODO: remove ugly hard-coded volume name
			Volumes:       []api.Volume{{Name: strings.Trim(strings.Split(cVolumes[0], ":")[0], "./"), VolumeSource: api.VolumeSource{EmptyDir: &api.EmptyDirVolumeSource{}}}},
			Containers:    []api.Container{{Name: serviceKey, Image: cImage, ImagePullPolicy: api.PullIfNotPresent}},
			RestartPolicy: api.RestartPolicyAlways,
			DNSPolicy:     api.DNSClusterFirst,
		},
	}
	log.Infof("%#v\n", pod)
	return pod
}
