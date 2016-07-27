package deploy

import (
	// "fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/api"
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/util/intstr"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
)

func DeployService(appKey string, appValue map[string]interface{}) error {
	service := BuildService(appKey, appValue)
	log.Infof("%#v\n", service)
	serviceYaml, err := yaml.Marshal(service)
	if err != nil {
		return err
	}
	serviceFile := strings.Join([]string{DefaultServiceFile, "-", appKey, YamlExtension}, "")
	if err := ioutil.WriteFile(serviceFile, serviceYaml, 0644); err != nil {
		return err
	}
	// Start deployment process
	kubectlCreateCmd := exec.Command("kubectl", "create", "-f", serviceFile)
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

func BuildService(appKey string, appValue map[string]interface{}) api.Service {
	var servicePorts []api.ServicePort
	for configKey, configValues := range appValue {
		if configKey == "ports" && reflect.TypeOf(configValues).Kind() == reflect.Slice {
			servicePorts = getServicePorts(configValues.([]interface{}))
		}
	}

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
			Ports:    servicePorts,
		},
	}
}

// getServicePorts bind docker-compose container port mapping to Kubernetes service port mapping
// E.g.: docker-compose has 80:3000 as host=80, container=3000; Kubernetes service will map 80 to "Port", 3000 to "TargetPort"
func getServicePorts(raw []interface{}) []api.ServicePort {
	var ports []api.ServicePort

	for _, cValue := range raw {
		var i32Ports []int32
		var servicePort int32
		var targetPort intstr.IntOrString
		cPort := cValue.(string)
		portList := strings.Split(cPort, ":") // 80:3000

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
			log.Errorln("Invalid Port Value From docker-compose File", cValue)
			continue
		}

		ports = append(ports, api.ServicePort{
			Name:       strings.Join([]string{"port-", strings.Replace(cPort, ":", "-", -1)}, ""),
			Protocol:   "TCP",
			Port:       servicePort,
			TargetPort: targetPort,
		})
	}
	return ports
}
