package create

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/config"
	"github.com/ghodss/yaml"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/api"
	"github.com/weitenghuang/dirigent-cli/pkg/resource"
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func ReplicationController(appName string, appConfig *config.ServiceConfig) (string, error) {
	if stop, err := resource.StopJobResourceWithError(appName); stop && err != nil {
		return "", err
	}

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
	rcLabel := resource.DefaultRCLabel(appName, "latest")
	podLabel := resource.DefaultPodLabel(appName, "latest")
	// Build single container
	appContainer := buildContainer(appName, appConfig)
	podVolumes := attachVolumeToContainer(appName, appConfig, &appContainer)
	podTemplateSpec := buildPodTemplateSpec(appName, &appContainer, podVolumes)

	log.Infof("RC Pod %v Template: %#v\n", podLabel, podTemplateSpec)

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
			Template: &podTemplateSpec,
		},
	}
}
