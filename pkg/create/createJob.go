package create

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/config"
	"github.com/ghodss/yaml"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/api"
	"github.com/weitenghuang/dirigent-cli/pkg/kubernetes/apis/batch"
	"github.com/weitenghuang/dirigent-cli/pkg/resource"
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func Job(appName string, appConfig *config.ServiceConfig) (string, error) {
	job := BuildJob(appName, appConfig)
	log.Infof("compose service %v config to Deployment: %#v\n", appName, job)
	jobYaml, err := yaml.Marshal(job)
	if err != nil {
		return "", err
	}
	jobFile := resource.DefaultJobFilePath(appName)
	if err := ioutil.WriteFile(jobFile, jobYaml, 0644); err != nil {
		return "", err
	}

	return jobFile, nil
}

func BuildJob(appName string, appConfig *config.ServiceConfig) batch.Job {
	jobLabel := resource.DefaultJobLabel(appName, "latest")
	podLabel := resource.DefaultPodLabel(appName, "latest")
	// Build single container
	appContainer := buildContainer(appName, appConfig)
	podVolumes := attachVolumeToContainer(appName, appConfig, &appContainer)
	podTemplateSpec := buildPodTemplateSpec(appName, &appContainer, podVolumes)
	podTemplateSpec.Spec.RestartPolicy = api.RestartPolicyNever

	log.Infof("Deployment Pod %v Template: %#v\n", resource.DefaultPodLabel(appName, "latest"), podTemplateSpec)

	return batch.Job{
		TypeMeta: unversioned.TypeMeta{Kind: "Job", APIVersion: "batch/v1"},
		ObjectMeta: api.ObjectMeta{
			Name:      jobLabel,
			Namespace: resource.DefaultK8sNamespace,
			Labels:    map[string]string{resource.DefaultSelectorKey: jobLabel},
		},
		Spec: batch.JobSpec{
			Selector: &unversioned.LabelSelector{
				MatchLabels: map[string]string{resource.DefaultSelectorKey: podLabel},
			},
			Template: podTemplateSpec,
		},
	}
}
