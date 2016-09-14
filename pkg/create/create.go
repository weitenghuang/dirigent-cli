package create

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/config"
	composeYaml "github.com/docker/libcompose/yaml"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/api"
	"github.com/weitenghuang/dirigent-cli/pkg/resource"
	"github.com/weitenghuang/dirigent-cli/pkg/utils"
	"k8s.io/kubernetes/pkg/util/intstr"
	"strconv"
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
		if !ok {
			log.Debugf("Unable to Find %v's Compose Service Config\n", name)
			continue
		}
		switch resourceType {
		case resource.ReplicationController:
			if resource.NotJobResource(name) {
				log.Infoln("Replication Controller File Creation Starts: ", name)
				if _, err := ReplicationController(name, composeServiceConfig); err != nil {
					log.Errorln("Error: ReplicationController File ", err, name, composeServiceConfig)
					return err
				}
			}
		case resource.Service:
			if resource.NotJobResource(name) {
				log.Infoln("Service File Creation Starts: ", name)
				if _, err := Service(name, composeServiceConfig); err != nil {
					log.Errorln("Error: Service File ", err, name, composeServiceConfig)
					return err
				}
			}
		case resource.Deployment:
			if resource.NotJobResource(name) {
				log.Infoln("Deployment File Creation Starts: ", name)
				if _, err := Deployment(name, composeServiceConfig); err != nil {
					log.Errorln("Error: Service File ", err, name, composeServiceConfig)
					return err
				}
			}
		case resource.Job:
			if !resource.NotJobResource(name) {
				log.Infoln("Job File Creation Starts: ", name)
				if _, err := Job(name, composeServiceConfig); err != nil {
					log.Errorln("Error: Service File ", err, name, composeServiceConfig)
					return err
				}
			}
		default:
			log.Infof("Resource type for %v is currently not supported yet.\n", name)
		}
	}

	return nil
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

// getServicePorts bind docker-compose container port mapping to Kubernetes service port mapping
// E.g.: docker-compose has 80:3000 as host=80, container=3000; Kubernetes service will map 80 to "Port", 3000 to "TargetPort"
func getServicePorts(composePorts []string) []api.ServicePort {
	var ports []api.ServicePort
	sep := ":"

	for _, cPort := range composePorts {
		var i32Ports []int32
		var servicePort int32
		var targetPort intstr.IntOrString
		portList := strings.Split(cPort, sep)

		for _, value := range portList {
			i64Port, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				log.Errorln("Invalid Port Value", err)
				continue
			}
			i32Ports = append(i32Ports, int32(i64Port))
		}

		if len(i32Ports) > 1 {
			servicePort = i32Ports[0]
			targetPort = intstr.IntOrString{Type: intstr.Int, IntVal: i32Ports[1]}
		} else if len(i32Ports) == 1 {
			servicePort = i32Ports[0]
			targetPort = intstr.IntOrString{Type: intstr.Int, IntVal: i32Ports[0]}
		} else {
			log.Errorln("Invalid Port Value From docker-compose File", cPort)
			continue
		}

		ports = append(ports, api.ServicePort{
			Name:       strings.Join([]string{"port-", strings.Replace(cPort, sep, "-", -1)}, ""),
			Protocol:   api.ProtocolTCP,
			Port:       servicePort,
			TargetPort: targetPort,
		})
	}
	return ports
}

func buildContainer(appName string, appConfig *config.ServiceConfig) api.Container {
	return api.Container{
		Name:            strings.Join([]string{appName, "-container"}, ""),
		Image:           appConfig.Image,
		ImagePullPolicy: api.PullIfNotPresent,
		Command:         []string(appConfig.Command),
		Env:             getEnvVarsFromCompose(appConfig.Environment),
		Ports:           getPortsFromCompose(appConfig.Ports),
	}
}

func buildPodTemplateSpec(appName string, appContainer *api.Container, podVolumes []api.Volume) api.PodTemplateSpec {
	podLabel := resource.DefaultPodLabel(appName, "latest")

	return api.PodTemplateSpec{
		ObjectMeta: api.ObjectMeta{
			Name:      podLabel,
			Namespace: resource.DefaultK8sNamespace,
			Labels:    map[string]string{resource.DefaultSelectorKey: podLabel},
		},
		Spec: api.PodSpec{
			Containers: []api.Container{*appContainer},
			Volumes:    podVolumes,
		},
	}
}

func attachVolumeToContainer(appName string, appConfig *config.ServiceConfig, appContainer *api.Container) []api.Volume {
	volumeLabel := strings.Join([]string{appName, "-storage"}, "")
	if appConfig.Volumes != nil && len(appConfig.Volumes.Volumes) > 0 {
		appContainer.VolumeMounts = []api.VolumeMount{
			api.VolumeMount{
				Name:      volumeLabel,
				MountPath: appConfig.Volumes.Volumes[0].Destination,
			},
		}

		return []api.Volume{
			api.Volume{
				Name: volumeLabel,
				VolumeSource: api.VolumeSource{
					EmptyDir: &api.EmptyDirVolumeSource{Medium: ""},
				},
			},
		}
	}
	return []api.Volume{}
}
