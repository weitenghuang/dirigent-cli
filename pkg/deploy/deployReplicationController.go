package deploy

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/config"
	composeYaml "github.com/docker/libcompose/yaml"
	"github.com/ghodss/yaml"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/api"
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func DeployReplicationController(appName string, appConfig *config.ServiceConfig) error {
	rc := BuildReplicationController(appName, appConfig)
	log.Infof("compose service %v config: %#v\n", appName, rc)
	rcYaml, err := yaml.Marshal(rc)
	if err != nil {
		return err
	}
	rcFile := strings.Join([]string{DefaultRcFile, "-", appName, YamlExtension}, "")
	if err := ioutil.WriteFile(rcFile, rcYaml, 0644); err != nil {
		return err
	}
	// Start deployment process
	kubectlCreateCmd := exec.Command("kubectl", "create", "-f", rcFile)
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

func BuildReplicationController(appName string, appConfig *config.ServiceConfig) api.ReplicationController {
	rcLabel := getRCSelectorLabel(appName, "latest")
	podLabel := getPodSelectorLabel(appName, "latest")

	// Build single container
	appContainer := api.Container{
		Name:            strings.Join([]string{appName, "-container"}, ""),
		Image:           appConfig.Image,
		ImagePullPolicy: api.PullIfNotPresent,
		Command:         []string(appConfig.Command),
		Env:             getEnvVarsFromCompose(appConfig.Environment),
		Ports:           getPortsFromCompose(appConfig.Ports),
	}

	podTemplateSpec := &api.PodTemplateSpec{
		ObjectMeta: api.ObjectMeta{
			Name:      podLabel,
			Namespace: DefaultK8sNamespace,
			Labels:    map[string]string{DefaultSelectorKey: podLabel},
		},
		Spec: api.PodSpec{
			Containers: []api.Container{appContainer},
		},
	}

	log.Infof("RC Pod %v Template: %#v\n", podLabel, *podTemplateSpec)

	return api.ReplicationController{
		TypeMeta: unversioned.TypeMeta{Kind: "ReplicationController", APIVersion: DefaultAPIVersion},
		ObjectMeta: api.ObjectMeta{
			Name:      rcLabel,
			Namespace: DefaultK8sNamespace,
			Labels:    map[string]string{DefaultSelectorKey: rcLabel},
		},
		Spec: api.ReplicationControllerSpec{
			Selector: map[string]string{DefaultSelectorKey: podLabel},
			Replicas: int32(1),
			Template: podTemplateSpec,
		},
	}
}

func getEnvVarsFromCompose(composeEnvs composeYaml.MaporEqualSlice) []api.EnvVar {
	var envs []api.EnvVar
	envMap := composeEnvs.ToMap()
	for envKey, envValue := range envMap {
		envs = append(envs, api.EnvVar{Name: envKey, Value: envValue})
	}
	return envs
}

// getPortsFromCompose currently ignores Docker host port value, only bind ContainerPort value
func getPortsFromCompose(composePorts []string) []api.ContainerPort {
	var ports []api.ContainerPort
	sep := ":"
	cProtocol := api.ProtocolTCP
	for _, cValue := range composePorts {
		var cPort string
		if strings.Contains(cValue, sep) {
			cPort = cValue[strings.Index(cValue, sep)+1:]
			cPort = strings.TrimSpace(cPort)
		} else {
			cPort = strings.TrimSpace(cValue)
		}
		i64Port, err := strconv.ParseInt(cPort, 10, 32)
		if err != nil {
			log.Errorln("Invalid Port Value", err)
			continue
		}
		ports = append(ports, api.ContainerPort{
			Name:          strings.Join([]string{"port-", cPort}, ""),
			ContainerPort: int32(i64Port),
			Protocol:      cProtocol,
		})
	}
	return ports
}
