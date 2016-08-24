package resource

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
