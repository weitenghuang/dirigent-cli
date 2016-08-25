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
	ReplicationController ResourceType = "rc"
	Service               ResourceType = "service"
)

func DefaultReplicationControllerFilePath(appName string) string {
	return strings.Join([]string{DefaultDeployPath, string(ReplicationController), "-", appName, YamlExtension}, "")
}

func DefaultServiceFilePath(appName string) string {
	return strings.Join([]string{DefaultDeployPath, string(Service), "-", appName, YamlExtension}, "")
}
