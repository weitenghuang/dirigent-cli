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
	"strconv"
	"strings"
)

func DeployPod(appKey string, appValue map[string]interface{}) error {
	pod := BuildPod(appKey, appValue)
	log.Infof("%#v\n", pod)
	podYaml, err := yaml.Marshal(pod)
	if err != nil {
		return err
	}
	podFile := strings.Join([]string{DefaultPodFile, "-", appKey, YamlExtension}, "")
	if err := ioutil.WriteFile(podFile, podYaml, 0644); err != nil {
		return err
	}
	// Start deployment process
	kubectlCreateCmd := exec.Command("kubectl", "create", "-f", podFile)
	kubectlCreateCmd.Stdout = os.Stdout
	kubectlCreateCmd.Stderr = os.Stderr
	if err := kubectlCreateCmd.Start(); err != nil {
		return err
	}
	if err := kubectlCreateCmd.Wait(); err != nil {
		return err
	}
	return nil
}

// BuildPod takes docker-compose's service to build a single "container" pod
func BuildPod(appKey string, appValue map[string]interface{}) api.Pod {
	var cCommands, cVolumes []string
	var cPorts []api.ContainerPort
	var cEnvs []api.EnvVar
	var cImage string
	// Loop through all configurations of current app
	for configKey, configValues := range appValue {
		if cValue, ok := configValues.(string); ok && configKey == "image" {
			cImage = cValue
		} else if configKey == "volumes" && reflect.TypeOf(configValues).Kind() == reflect.Slice {
			for _, cValue := range configValues.([]interface{}) {
				sValue := cValue.(string)
				volumeList := strings.Split(sValue, ":")
				if len(volumeList) < 2 {
					cVolumes = append(cVolumes, sValue)
				} else {
					cVolumes = append(cVolumes, volumeList[1])
				}
			}
			log.Infoln(appKey, " container volumes:", cVolumes)
		} else if configKey == "ports" && reflect.TypeOf(configValues).Kind() == reflect.Slice {
			cPorts = getContainerPorts(configValues.([]interface{}))
		} else if configKey == "environment" && reflect.TypeOf(configValues).Kind() == reflect.Map {
			cEnvs = getEnvVars(configValues.(map[string]interface{}))
		} else if configKey == "command" {
			switch configValues.(type) {
			default:
				cCommands = strings.Split(configValues.(string), " ")
			}
		}
	}
	// Build single container
	appContainer := api.Container{
		Name:            strings.Join([]string{appKey, "-container"}, ""),
		Image:           cImage,
		ImagePullPolicy: api.PullIfNotPresent,
		Command:         cCommands,
		Env:             cEnvs,
		Ports:           cPorts,
	}

	label := getPodSelectorLabel(appKey, "latest")
	return api.Pod{ // Basic fields.
		TypeMeta: unversioned.TypeMeta{Kind: "Pod", APIVersion: DefaultAPIVersion},
		ObjectMeta: api.ObjectMeta{
			Name:      label,
			Namespace: DefaultK8sNamespace,
			Labels:    map[string]string{DefaultSelectorKey: label},
		},
		Spec: api.PodSpec{
			Volumes: []api.Volume{
				{
					Name:         strings.Join([]string{appKey, "-volume"}, ""),
					VolumeSource: api.VolumeSource{EmptyDir: &api.EmptyDirVolumeSource{}},
				},
			},
			Containers:    []api.Container{appContainer},
			RestartPolicy: api.RestartPolicyAlways,
			DNSPolicy:     api.DNSClusterFirst,
		},
	}
}

// getContainerPorts currently ignores Docker host port value, only bind ContainerPort value
func getContainerPorts(raw []interface{}) []api.ContainerPort {
	var ports []api.ContainerPort
	for _, cValue := range raw {
		cPort := cValue.(string)
		portList := strings.Split(cPort, ":") // 8080:8080
		if len(portList) > 1 {
			//portList[0] is as public/host port
			cPort = portList[1]
		}
		i64Port, err := strconv.ParseInt(cPort, 10, 32)
		if err != nil {
			log.Errorln("Invalid Port Value", err)
			continue
		}
		ports = append(ports, api.ContainerPort{
			Name:          strings.Join([]string{"port-", cPort}, ""),
			ContainerPort: int32(i64Port),
		})
	}
	return ports
}

func getEnvVars(raw map[string]interface{}) []api.EnvVar {
	var envs []api.EnvVar
	for envKey, envValue := range raw {
		var sValue string
		switch envValue.(type) {
		case int32:
			sValue = strconv.FormatInt(envValue.(int64), 10)
		case float64:
			sValue = strconv.FormatFloat(envValue.(float64), 'f', -1, 64)
		default:
			sValue = envValue.(string)
		}
		envs = append(envs, api.EnvVar{Name: envKey, Value: sValue})
	}
	return envs
}
