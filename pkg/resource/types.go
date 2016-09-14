package resource

import (
	"fmt"
	"strings"
)

const (
	YamlExtension       string = ".yml"
	DefaultDeployPath   string = "/opt/deploy/"
	DefaultAPIVersion   string = "v1"
	DefaultK8sNamespace string = "default"
	DefaultSelectorKey  string = "name"
)

type ResourceType string

const (
	Deployment            ResourceType = "deployment"
	Job                   ResourceType = "job"
	Pod                   ResourceType = "pod"
	ReplicationController ResourceType = "rc"
	Service               ResourceType = "service"
)

type EnvVar string

const (
	BaseFqdn EnvVar = "BASE_FQDN"
)

const (
	DefaultFqdn     string = "com.default"
	DefaultReplicas int32  = 1
)

type ClusterConfig struct {
	Replicas int32
}

func DefaultJobFilePath(appName string) string {
	return defaultFilePath(appName, Job)
}

func DefaultDeploymentFilePath(appName string) string {
	return defaultFilePath(appName, Deployment)
}

func DefaultReplicationControllerFilePath(appName string) string {
	return defaultFilePath(appName, ReplicationController)
}

func DefaultServiceFilePath(appName string) string {
	return defaultFilePath(appName, Service)
}

func defaultFilePath(appName string, resourceType ResourceType) string {
	return strings.Join([]string{DefaultDeployPath, string(resourceType), "-", appName, YamlExtension}, "")
}

func DefaultPodLabel(appName string, version string) string {
	return defaultLabel(appName, version, Pod)
}

func DefaultRCLabel(appName string, version string) string {
	return defaultLabel(appName, version, ReplicationController)
}

func DefaultJobLabel(appName string, version string) string {
	return defaultLabel(appName, version, Job)
}

func DefaultDeploymentLabel(appName string, version string) string {
	return defaultLabel(appName, version, Deployment)
}

func defaultLabel(name string, version string, resourceType ResourceType) string {
	return strings.Join([]string{name, "-", version, "-", string(resourceType)}, "")
}

func DefaultActiveDeadlineSeconds() *int64 {
	activeDeadlineSeconds := int64(300) // 300 seconds
	return &activeDeadlineSeconds
}

func NotJobResource(appName string) bool {
	return !strings.Contains(appName, "-job") && !strings.Contains(appName, "-init")
}

func StopJobResourceWithError(appName string) (stop bool, err error) {
	if notJob := NotJobResource(appName); !notJob {
		return true, fmt.Errorf("%#v should be a \"job\" resource.", appName)
	}
	return false, nil
}
