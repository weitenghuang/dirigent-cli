package create

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/config"
	"github.com/ghodss/yaml"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/api"
	"github.com/weitenghuang/dirigent-cli/pkg/resource"
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/util/intstr"
	"strconv"
	"strings"
)

func Service(appName string, appConfig *config.ServiceConfig) (string, error) {
	service := BuildService(appName, appConfig)
	log.Infof("compose service %v config to service: %#v\n", appName, service)
	serviceYaml, err := yaml.Marshal(service)
	if err != nil {
		return "", err
	}
	serviceFile := resource.DefaultServiceFilePath(appName)
	if err := ioutil.WriteFile(serviceFile, serviceYaml, 0644); err != nil {
		return "", err
	}
	return serviceFile, nil
}

func BuildService(appName string, appConfig *config.ServiceConfig) api.Service {
	servicePorts := getServicePorts(appConfig.Ports)
	label := strings.Join([]string{appName, "-service"}, "")
	return api.Service{
		TypeMeta: unversioned.TypeMeta{Kind: "Service", APIVersion: resource.DefaultAPIVersion},
		ObjectMeta: api.ObjectMeta{
			Name:      label,
			Namespace: resource.DefaultK8sNamespace,
			Labels:    map[string]string{resource.DefaultSelectorKey: label},
		},
		Spec: api.ServiceSpec{
			Selector: map[string]string{resource.DefaultSelectorKey: resource.DefaultPodLabel(appName, "latest")},
			Type:     api.ServiceTypeClusterIP,
			Ports:    servicePorts,
		},
	}
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
