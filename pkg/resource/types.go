package resource

import (
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
