package create

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/config"
	"github.com/ghodss/yaml"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/api"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/apis/extensions"
	"github.com/weitenghuang/dirigent-cli/pkg/resource"
	"github.com/weitenghuang/dirigent-cli/pkg/utils"
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func Deployment(appName string, appConfig *config.ServiceConfig) (string, error) {
	if stop, err := utils.StopJobResourceWithError(appName); stop && err != nil {
		return "", err
	}

	deployment := BuildDeployment(appName, appConfig)
	log.Infof("compose service %v config to Deployment: %#v\n", appName, deployment)
	deploymentYaml, err := yaml.Marshal(deployment)
	if err != nil {
		return "", err
	}
	deploymentFile := resource.DefaultDeploymentFilePath(appName)
	if err := ioutil.WriteFile(deploymentFile, deploymentYaml, 0644); err != nil {
		return "", err
	}

	return deploymentFile, nil
}

func BuildDeployment(appName string, appConfig *config.ServiceConfig) extensions.Deployment {
	deploymentLabel := resource.DefaultDeploymentLabel(appName, "latest")
	podLabel := resource.DefaultPodLabel(appName, "latest")
	// Build single container
	appContainer := buildContainer(appName, appConfig)
	podVolumes := attachVolumeToContainer(appName, appConfig, &appContainer)
	podTemplateSpec := buildPodTemplateSpec(appName, &appContainer, podVolumes)

	log.Infof("Deployment Pod %v Template: %#v\n", resource.DefaultPodLabel(appName, "latest"), podTemplateSpec)

	return extensions.Deployment{
		TypeMeta: unversioned.TypeMeta{Kind: "Deployment", APIVersion: "extensions/v1beta1"},
		ObjectMeta: api.ObjectMeta{
			Name:      deploymentLabel,
			Namespace: resource.DefaultK8sNamespace,
			Labels:    map[string]string{resource.DefaultSelectorKey: deploymentLabel},
		},
		Spec: extensions.DeploymentSpec{
			Selector: &unversioned.LabelSelector{
				MatchLabels: map[string]string{resource.DefaultSelectorKey: podLabel},
			},
			Replicas: int32(1),
			Template: podTemplateSpec,
		},
	}
}
