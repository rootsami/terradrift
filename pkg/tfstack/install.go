package tfstack

import (
	"context"
	"errors"
	"os"
	"regexp"
	"sync"

	"github.com/rootsami/terradrift/pkg/config"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

var mutexes = make(map[string]*sync.Mutex) // map of mutexes, one per stack path

// install should recieve a terraform version and return the execution path
func install(stack config.Stack, workspace string) (string, string, error) {

	// Check if terraform version is defined in the stack files
	v, err := detectTFVersion(workspace + stack.Path)
	if err != nil {
		log.WithField("stack", stack.Name).Debug("Terraform version not defined in the stack")
		return "", "", err
	}

	// To make sure returned value doesn't include '>='
	tfver := regexp.MustCompile(`[^a-zA-Z0-9. ]+`).ReplaceAllString(v, "")

	execPathDir := os.TempDir() + tfver
	execPath, err := downloadBinary(execPathDir, tfver)
	if err != nil {
		return "", "", err
	}

	return execPath, tfver, nil

}

// Detect terraform version based on its definition in tf files
func detectTFVersion(path string) (string, error) {

	module, err := tfconfig.LoadModule(path)
	if err != nil {
		return "", err
	}

	// Check if terraform version is defined in the module
	if len(module.RequiredCore) == 0 {
		return "", err
	}

	tfversion := module.RequiredCore[0]

	return tfversion, nil

}

// Download terraform binary based on the version defined in the stack files and return the execution path
func downloadBinary(dir, tfver string) (string, error) {

	// Create a mutex for each version of terraform
	// This prevents multiple parallel executions from trying to download the same version
	mutex, ok := mutexes[tfver]
	if !ok {
		mutex = &sync.Mutex{}
		mutexes[tfver] = mutex
	}

	mutex.Lock()
	defer mutex.Unlock()

	execPath := dir + "/terraform"

	// Check if binary already exists in the specified temp directory
	if _, err := os.Stat(execPath); !errors.Is(err, os.ErrNotExist) {
		log.WithFields(log.Fields{"version": tfver}).Debug("Skipping download, Terraform binary found...")

		return execPath, nil

	} else {
		// Create temp directory and download terraform binary to it
		os.MkdirAll(dir, os.ModePerm)
		installer := &releases.ExactVersion{
			Product:    product.Terraform,
			Version:    version.Must(version.NewVersion(tfver)),
			InstallDir: dir,
		}

		log.WithFields(log.Fields{"version": tfver}).Debug("Downloading Terraform...")

		execPath, err := installer.Install(context.Background())
		if err != nil {
			return "", err
		}

		return execPath, nil
	}

}
