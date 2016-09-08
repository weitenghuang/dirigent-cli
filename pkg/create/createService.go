package create

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/config"
	"github.com/ghodss/yaml"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/api"
	"github.com/weitenghuang/dirigent-cli/pkg/resource"
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api/unversioned"
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
