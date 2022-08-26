package main

import (
	"context"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

// The main function for installing the exact version of the stack, initiate and run terraform plan
func stackScan(name string) (string, error) {

	config := configLoader()

	stack, validStack := stackExists(name, config.Stacks)
	if validStack {

		// Checkout new commits/updates in repository
		gitPull(workspace)

		// TODO: Install should opt-out of this function to allow download only happen once. Where you pass the needed version
		// and it return the execPath
		installer := &releases.ExactVersion{
			Product: product.Terraform,
			Version: version.Must(version.NewVersion(stack.Version)),
		}

		log.WithFields(log.Fields{
			"stack":   stack.Name,
			"version": stack.Version,
		}).Info("Installing Terraform...")
		execPath, err := installer.Install(context.Background())
		if err != nil {
			log.WithFields(log.Fields{
				"stack":   stack.Name,
				"version": stack.Version,
			}).Errorf("Installing Terraform: %s", err)
			return "Error:", err
		}

		tf, err := tfexec.NewTerraform(workspace+stack.Path, execPath)
		if err != nil {
			log.WithFields(log.Fields{
				"stack":   stack.Name,
				"version": stack.Version,
			}).Errorf("Running NewTerraform: %s", err)
			return "Error:", err
		}

		response, err := stackPlan(stack, tf)
		if err != nil {
			log.WithFields(log.Fields{
				"stack":   stack.Name,
				"version": stack.Version,
			}).Error(err)
		}

		return response, err

	} else {
		log.WithFields(log.Fields{
			"stack": name,
		}).Error("STACK WAS NOT FOUND")
		err := fmt.Errorf("ERROR: STACK WAS NOT FOUND")
		return response, err
	}

}

// stackPlan has been separated to be called indivisually with the schedule to avoid downloading/installing
// the required terraform version.
// initializing is part of the plan incase new modules added/upgraded to the stack code
func stackPlan(stack Stack, tf *tfexec.Terraform) (string, error) {

	var err error
	// Stacks come with two different structures:
	// 1. All resources for multiple stacks (environments) exist in one directory and backend initialization is done with environments/<name>.hcl
	// 2. Regular stack where all resources, tfvars and backend configs are in the same directory

	// during the initialization, .terrafom directory collides with other environments' .terraform
	// causing lots of issue with local terraform.tfstate while performing terraform plan
	// solution would be export TF_DATA_DIR with customized .terraform naming to avoid the issue

	log.WithFields(log.Fields{
		"stack":   stack.Name,
		"version": stack.Version,
	}).Info("Initializing Terraform...")

	tfEnvVars := map[string]string{
		"TF_DATA_DIR": ".terraform." + stack.Name,
	}
	tf.SetEnv(tfEnvVars)

	if len(stack.Backend) > 0 {

		err := tf.Init(context.Background(), tfexec.Upgrade(false), tfexec.BackendConfig(stack.Backend))
		if err != nil {
			log.WithFields(log.Fields{
				"stack":   stack.Name,
				"version": stack.Version,
			}).Error("Running Init")
			return "Error:", err
		}

	} else {

		err := tf.Init(context.Background(), tfexec.Upgrade(false))
		if err != nil {
			log.WithFields(log.Fields{
				"stack":   stack.Name,
				"version": stack.Version,
			}).Error("Running Init")
			return "Error:", err
		}
	}

	// Create TF Plan options
	planFile := workspace + stack.Path + "/" + stack.Name + ".plan"
	stackPlanFile := tfexec.Out(planFile)

	if len(stack.TFvars) > 0 {
		plan, err := tf.Plan(context.Background(), stackPlanFile, tfexec.VarFile(stack.TFvars))
		if err != nil {
			log.WithFields(log.Fields{
				"stack":   stack.Name,
				"version": stack.Version,
			}).Error("Running Plan")
			return "Error:", err
		}

		response, err = showPlan(plan, planFile, stack.Name, tf)
		if err != nil {
			return "Error:", err
		}

	} else {
		plan, err := tf.Plan(context.Background(), stackPlanFile)
		if err != nil {
			log.WithFields(log.Fields{
				"stack":   stack.Name,
				"version": stack.Version,
			}).Error("Running Plan")
			return "Error:", err
		}

		response, err = showPlan(plan, planFile, stack.Name, tf)
		if err != nil {
			return "Error:", err
		}
	}

	return response, err
}

func showPlan(plan bool, planFile string, name string, tf *tfexec.Terraform) (string, error) {

	var err error
	if plan {

		state, err := tf.ShowPlanFileRaw(context.Background(), planFile)
		if err != nil {
			log.WithFields(log.Fields{
				"stack": name,
			}).Errorf("Running show: %s", err)
			return "ERROR:", err
		}

		re := regexp.MustCompile("Plan:.*")
		summary := re.FindString(state)
		log.WithFields(log.Fields{
			"stack":   name,
			"summary": summary,
		}).Info("CHANGES DETECTED...")
		return fmt.Sprintf("CHANGES DETECTED... %s", summary), err

	} else {
		log.WithFields(log.Fields{
			"stack":   name,
			"summary": "No changes. Infrastructure matches the configuration.",
		}).Info("NO CHANGES...")
		return "No changes. Infrastructure matches the configuration.", err
	}

}
