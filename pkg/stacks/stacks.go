package stacks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/rootsami/terradrift/pkg/config"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

// The main function for installing the exact version of the stack, initiate and run terraform plan
func StackScan(name, workspace, configPath string) (string, error) {

	var response string
	config := config.ConfigLoader(configPath)

	stack, validStack := stackExists(name, config.Stacks)
	if validStack {

		// The path for terrafom binary
		execPath := install(stack)

		tf, err := tfexec.NewTerraform(workspace+stack.Path, execPath)
		if err != nil {
			log.WithFields(log.Fields{"stack": stack.Name, "version": stack.Version}).Errorf("Running NewTerraform: %s", err)
			return "", err
		}

		response, err := stackPlan(workspace, stack, tf)
		if err != nil {
			log.WithFields(log.Fields{"stack": stack.Name, "version": stack.Version}).Error(err)
		}

		return response, err

	} else {
		log.WithFields(log.Fields{"stack": name}).Error("STACK WAS NOT FOUND")
		err := fmt.Errorf("ERROR: STACK WAS NOT FOUND")
		return response, err
	}

}

// StackPlan is separated to be called indivisually to avoid downloading/installing terraform binaries of the same version.
// Initializing is part of the plan incase new modules added/upgraded to the stack code
func stackPlan(workspace string, stack config.Stack, tf *tfexec.Terraform) (string, error) {

	// Stacks come with two different structures:
	// 1. All resources for multiple stacks (environments) exist in one directory and backend initialization is done with environments/<name>.hcl
	// 2. Regular stack where all resources, tfvars and backend configs are in the same directory

	// during the initialization, .terrafom directory collides with other environments' .terraform
	// causing lots of issues with local terraform.tfstate while performing terraform plan
	// solution would be export TF_DATA_DIR with customized .terraform naming to avoid the issue

	log.WithFields(log.Fields{"stack": stack.Name, "version": stack.Version}).Info("Initializing Terraform...")

	var response string
	tfEnvVars := map[string]string{
		"TF_DATA_DIR": ".terraform." + stack.Name,
	}
	tf.SetEnv(tfEnvVars)

	err := tf.Init(context.Background(), tfexec.Upgrade(false), tfexec.BackendConfig(stack.Backend))
	if err != nil {
		log.WithFields(log.Fields{"stack": stack.Name, "version": stack.Version}).Error("Running Init")
		return "", err
	}

	// Create TF Plan options
	planFile := workspace + stack.Path + "/" + stack.Name + ".plan"
	stackPlanFile := tfexec.Out(planFile)

	if len(stack.TFvars) > 0 {
		plan, err := tf.Plan(context.Background(), stackPlanFile, tfexec.VarFile(stack.TFvars))
		if err != nil {
			log.WithFields(log.Fields{"stack": stack.Name, "version": stack.Version}).Error("Running Plan")
			return "", err
		}

		response, err = showPlan(plan, planFile, stack.Name, tf)
		if err != nil {
			return "", err
		}

	} else {
		plan, err := tf.Plan(context.Background(), stackPlanFile)
		if err != nil {
			log.WithFields(log.Fields{"stack": stack.Name, "version": stack.Version}).Error("Running Plan")
			return "", err
		}

		response, err = showPlan(plan, planFile, stack.Name, tf)
		if err != nil {
			return "", err
		}
	}

	return response, err
}

func showPlan(plan bool, planFile string, name string, tf *tfexec.Terraform) (string, error) {

	var err error
	if plan {

		state, err := tf.ShowPlanFileRaw(context.Background(), planFile)
		if err != nil {
			log.WithFields(log.Fields{"stack": name}).Errorf("Running show: %s", err)
			return "", err
		}

		re := regexp.MustCompile("Plan:.*")
		summary := re.FindString(state)
		log.WithFields(log.Fields{"stack": name, "summary": summary}).Info("CHANGES DETECTED...")
		return fmt.Sprintf("CHANGES DETECTED... %s", summary), err

	} else {
		log.WithFields(log.Fields{"stack": name, "summary": "No changes. Infrastructure matches the configuration."}).Info("NO CHANGES...")
		return "No changes. Infrastructure matches the configuration.", err
	}

}

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

// stackExists checks if the requested stack exists in the configration file
func stackExists(name string, stacks []config.Stack) (stack config.Stack, result bool) {
	result = false
	for _, stack := range stacks {
		if stack.Name == name {
			result = true
			return stack, result
		}
	}
	return stack, result
}
