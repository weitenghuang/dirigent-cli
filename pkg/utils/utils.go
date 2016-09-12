package utils

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/lookup"
	"github.com/docker/libcompose/project"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
)

func ParseDockerCompose(filePath string) (composeObject *project.Project, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Infof("Recovered in pkg/utils/utils.go libcomposeParser(filePath string) %#v\n", r)
			debug.PrintStack()
			err = r.(error)
		}
	}()
	context := &project.Context{}
	if filePath == "" {
		filePath = "docker-compose.yml"
	}
	context.ComposeFiles = []string{filePath}
	if context.EnvironmentLookup == nil {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		context.EnvironmentLookup = &lookup.ComposableEnvLookup{
			Lookups: []config.EnvironmentLookup{
				&lookup.EnvfileLookup{
					Path: filepath.Join(cwd, ".env"),
				},
				&lookup.OsEnvLookup{},
			},
		}
	}
	composeObject = project.NewProject(context, nil, nil)
	err = composeObject.Parse()
	if err != nil {
		log.Fatalf("Failed to load compose file", err)
		return nil, err
	}
	log.Infof("Post-parsing: %#v\n", composeObject)
	return composeObject, nil
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

func KubectlCreateCmd(filepath string) *exec.Cmd {
	return exec.Command("kubectl", "create", "--validate", "-f", filepath)
}
