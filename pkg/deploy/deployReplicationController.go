package deploy

import (
	"github.com/docker/libcompose/config"
	"github.com/weitenghuang/dirigent-cli/pkg/create"
	"github.com/weitenghuang/dirigent-cli/pkg/utils"
	"os"
	"os/exec"
)

func ReplicationController(appName string, appConfig *config.ServiceConfig) error {
	if stop, err := utils.StopJobResourceWithError(appName); stop && err != nil {
		return err
	}

	rcFile, err := create.ReplicationController(appName, appConfig)
	if err != nil {
		return err
	}
	// Start deployment process
	kubectlCreateCmd := exec.Command("kubectl", "create", "-f", rcFile)
	kubectlCreateCmd.Stdout = os.Stdout
	kubectlCreateCmd.Stderr = os.Stderr
	if err := kubectlCreateCmd.Start(); err != nil {
		return err
	}
	if err := kubectlCreateCmd.Wait(); err != nil {
		return err
	}
	return nil
}
