package stacks

import (
	"context"
	"errors"
	"os"

	"github.com/rootsami/terradrift/pkg/config"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
)

// install should recieve a terraform version and return the execution path
func install(stack config.Stack) (execPath string) {

	execPathDir := "/tmp/binaries/" + stack.Version
	execPath = execPathDir + "/terraform"

	if _, err := os.Stat(execPath); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(execPathDir, os.ModePerm)
		installer := &releases.ExactVersion{
			Product:    product.Terraform,
			Version:    version.Must(version.NewVersion(stack.Version)),
			InstallDir: execPathDir,
		}

		log.WithFields(log.Fields{"stack": stack.Name, "version": stack.Version}).Info("Installing Terraform...")

		execPath, err := installer.Install(context.Background())
		if err != nil {
			log.WithFields(log.Fields{"stack": stack.Name, "version": stack.Version}).Errorf("Installing Terraform: %s", err)
		}
		return execPath

	} else {

		log.WithFields(log.Fields{"stack": stack.Name, "version": stack.Version}).Info("Skipping download, Terraform binary found...")

		return execPath
	}

}
