package create

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/config"
	composeYaml "github.com/docker/libcompose/yaml"
	"github.com/ghodss/yaml"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/api"
	"github.com/weitenghuang/dirigent-cli/pkg/resource"
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"strconv"
	"strings"
)

func ReplicationController(appName string, appConfig *config.ServiceConfig) (string, error) {
	rc := BuildReplicationController(appName, appConfig)
	log.Infof("compose service %v config to RC: %#v\n", appName, rc)
	rcYaml, err := yaml.Marshal(rc)
	if err != nil {
		return "", err
	}
	rcFile := resource.DefaultReplicationControllerFilePath(appName)
	if err := ioutil.WriteFile(rcFile, rcYaml, 0644); err != nil {
		return "", err
	}

	return rcFile, nil
}

func BuildReplicationController(appName string, appConfig *config.ServiceConfig) api.ReplicationController {
	rcLabel := getRCSelectorLabel(appName, "latest")
	podLabel := getPodSelectorLabel(appName, "latest")
	volumeLabel := strings.Join([]string{appName, "-storage"}, "")
	// Build single container
	appContainer := api.Container{
		Name:            strings.Join([]string{appName, "-container"}, ""),
		Image:           appConfig.Image,
		ImagePullPolicy: api.PullIfNotPresent,
		Command:         []string(appConfig.Command),
		Env:             getEnvVarsFromCompose(appConfig.Environment),
		Ports:           getPortsFromCompose(appConfig.Ports),
	}
	var podVolumes []api.Volume
	if appConfig.Volumes != nil && len(appConfig.Volumes.Volumes) > 0 {
		appContainer.VolumeMounts = []api.VolumeMount{
			api.VolumeMount{
				Name:      volumeLabel,
				MountPath: appConfig.Volumes.Volumes[0].Destination,
			},
		}

		podVolumes = []api.Volume{
			api.Volume{
				Name: volumeLabel,
				VolumeSource: api.VolumeSource{
					EmptyDir: &api.EmptyDirVolumeSource{Medium: ""},
				},
			},
		}
	}
	podTemplateSpec := &api.PodTemplateSpec{
		ObjectMeta: api.ObjectMeta{
			Name:      podLabel,
			Namespace: resource.DefaultK8sNamespace,
			Labels:    map[string]string{resource.DefaultSelectorKey: podLabel},
		},
		Spec: api.PodSpec{
			Containers: []api.Container{appContainer},
			Volumes:    podVolumes,
		},
	}

	log.Infof("RC Pod %v Template: %#v\n", podLabel, *podTemplateSpec)

	return api.ReplicationController{
		TypeMeta: unversioned.TypeMeta{Kind: "ReplicationController", APIVersion: resource.DefaultAPIVersion},
		ObjectMeta: api.ObjectMeta{
			Name:      rcLabel,
			Namespace: resource.DefaultK8sNamespace,
			Labels:    map[string]string{resource.DefaultSelectorKey: rcLabel},
		},
		Spec: api.ReplicationControllerSpec{
			Selector: map[string]string{resource.DefaultSelectorKey: podLabel},
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
			Protocol:      api.ProtocolTCP,
		})
	}
	return ports
}
