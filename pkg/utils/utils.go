package utils

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
	"io/ioutil"
	"os"
	"runtime/debug"
)

func ParseDockerCompose(filePath string) (composeObject *project.Project, err error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	sanitizedCompose := []byte(os.ExpandEnv(string(content)))
	tmpfile, err := ioutil.TempFile("/tmp", "sanitized-")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(sanitizedCompose); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return libcomposeParser(tmpfile.Name())
}

func libcomposeParser(filePath string) (composeObject *project.Project, err error) {
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
	composeObject = project.NewProject(&context.Context, nil, nil)
	err = composeObject.Parse()
	if err != nil {
		log.Fatalf("Failed to load compose file", err)
		return nil, err
	}
	log.Infof("Post-parsing: %#v\n", composeObject)
	return composeObject, nil
}
