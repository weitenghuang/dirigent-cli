package deploy

import (
	"github.com/docker/libcompose/config"
	"github.com/weitenghuang/dirigent-cli/pkg/create"
	"github.com/weitenghuang/dirigent-cli/pkg/utils"
	"os"
)

func Job(appName string, appConfig *config.ServiceConfig) error {
	jobFile, err := create.Job(appName, appConfig)
	if err != nil {
		return err
	}
	// Start deployment process
	kubectlCreateCmd := utils.KubectlCreateCmd(jobFile)
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
