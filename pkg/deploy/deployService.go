package deploy

import (
	// "fmt"
	log "github.com/Sirupsen/logrus"
	// "github.com/ghodss/yaml"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/api"
	// "io/ioutil"
	"k8s.io/kubernetes/pkg/api/unversioned"
	// "os"
	// "os/exec"
	// "reflect"
	"k8s.io/kubernetes/pkg/util/intstr"
	"strings"
)

func DeployService(appKey string, appValue map[string]interface{}) error {
	service := BuildService(appKey, appValue)
	log.Infoln(service)
	return nil
}

func BuildService(appKey string, appValue map[string]interface{}) api.Service {
	label := strings.Join([]string{appKey, "-service"}, "")
	return api.Service{
		TypeMeta: unversioned.TypeMeta{Kind: "Service", APIVersion: DefaultAPIVersion},
		ObjectMeta: api.ObjectMeta{
			Name:      label,
			Namespace: DefaultK8sNamespace,
			Labels:    map[string]string{DefaultSelectorKey: label},
		},
		Spec: api.ServiceSpec{
			Selector: map[string]string{DefaultSelectorKey: getPodSelectorLabel(appKey, "latest")},
			Type:     api.ServiceTypeClusterIP,
			Ports: []api.ServicePort{
				{Name: "p", Protocol: "TCP", Port: 28015, TargetPort: intstr.FromInt(28015)},
			},
		},
	}
}
