package deploy

import (
	"github.com/docker/libcompose/config"
	"github.com/weitenghuang/dirigent-cli/pkg/create"
	"os"
	"os/exec"
)

// func DeployService(appKey string, appValue map[string]interface{}) error {
func Service(appName string, appConfig *config.ServiceConfig) error {
	serviceFile, err := create.Service(appName, appConfig)
	if err != nil {
		return err
	}

	// Start deployment process
	kubectlCreateCmd := exec.Command("kubectl", "create", "-f", serviceFile)
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
