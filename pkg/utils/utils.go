package utils

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/lookup"
	"github.com/docker/libcompose/project"
	"path/filepath"
	// "io/ioutil"
	"os"
	"runtime/debug"
)

func ParseDockerCompose(filePath string) (composeObject *project.Project, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Infof("Recovered in pkg/utils/utils.go libcomposeParser(filePath string) %#v\n", r)
			debug.PrintStack()
			err = r.(error)
		}
	}()
	context := &docker.Context{}
	if filePath == "" {
		filePath = "docker-compose.yml"
	}
	context.ComposeFiles = []string{filePath}
	if context.Context.EnvironmentLookup == nil {
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
	composeObject = project.NewProject(&context.Context, nil, nil)
	err = composeObject.Parse()
	if err != nil {
		log.Fatalf("Failed to load compose file", err)
		return nil, err
	}
	log.Infof("Post-parsing: %#v\n", composeObject)
	return composeObject, nil
}
