package deploy

import (
	"github.com/docker/libcompose/config"
	"github.com/weitenghuang/dirigent-cli/pkg/create"
	"github.com/weitenghuang/dirigent-cli/pkg/resource"
	"github.com/weitenghuang/dirigent-cli/pkg/utils"
	"os"
)

func Deployment(appName string, appConfig *config.ServiceConfig) error {
	if stop, err := resource.StopJobResourceWithError(appName); stop && err != nil {
		return err
	}

	deploymentFile, err := create.Deployment(appName, appConfig)
	if err != nil {
		return err
	}
	// Start deployment process
	kubectlCreateCmd := utils.KubectlCreateCmd(deploymentFile)
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
