package utils

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
)

func ParseDockerCompose(filename string) (*project.Project, error) {
	context := &docker.Context{}
	if filename == "" {
		filename = "docker-compose.yml"
	}
	context.ComposeFiles = []string{filename}
	composeObject := project.NewProject(&context.Context, nil, nil)
	err := composeObject.Parse()
	if err != nil {
		log.Fatalf("Failed to load compose file", err)
		return nil, err
	}
	log.Infof("Post-parsing: %#v\n", composeObject)
	return composeObject, nil
}
