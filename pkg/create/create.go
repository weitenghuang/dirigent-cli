package create

import (
	log "github.com/Sirupsen/logrus"
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
		notJob := utils.NotJobResource(name)
		if ok && notJob { // Regular resources
			switch resourceType {
			case resource.ReplicationController:
				log.Infoln("Replication Controller File Creation Starts: ", name)
				if _, err := ReplicationController(name, composeServiceConfig); err != nil {
					log.Errorln("Error: ReplicationController File ", err, name, composeServiceConfig)
					return err
				}
			case resource.Service:
				log.Infoln("Service File Creation Starts: ", name)
				if _, err := Service(name, composeServiceConfig); err != nil {
					log.Errorln("Error: Service File ", err, name, composeServiceConfig)
					return err
				}
			case resource.Deployment:
				log.Infoln("Deployment File Creation Starts: ", name)
				if _, err := Deployment(name, composeServiceConfig); err != nil {
					log.Errorln("Error: Service File ", err, name, composeServiceConfig)
					return err
				}
			}
		} else if ok && resourceType == resource.Job && !notJob { // Special Resource
			log.Infof("%v\n", resourceType)
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
