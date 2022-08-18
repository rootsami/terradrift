package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

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

		installer := &releases.ExactVersion{
			Product: product.Terraform,
			Version: version.Must(version.NewVersion(stack.Version)),
		}

		log.Printf("%s: Installing Terraform %s ...", stack.Name, stack.Version)
		execPath, err := installer.Install(context.Background())
		if err != nil {
			log.Printf("%s: Error installing Terraform: %s", stack.Name, err)
			return "Error:", err
		}

		tf, err := tfexec.NewTerraform(workspace+stack.Path, execPath)
		if err != nil {
			log.Printf("%s: Error running NewTerraform: %s", stack.Name, err)
			return "Error:", err
		}

		response, err := stackPlan(stack, tf)
		if err != nil {
			log.Printf("%s: %s", stack.Name, err)
		}

		return response, err

	} else {
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
	log.Printf("%s: Initializing Terraform...", stack.Name)
	if len(stack.Backend) > 0 {

		// during the initialization, .terrafom directory collides with other environments' .terraform
		// causing lots of issue with local terraform.tfstate while performing terraform plan
		// solution would be export TF_DATA_DIR with customized .terraform naming to avoid the issue
		os.Setenv("TF_DATA_DIR", ".terraform."+stack.Name)
		err := tf.Init(context.Background(), tfexec.Upgrade(false), tfexec.BackendConfig(stack.Backend))
		if err != nil {
			log.Printf("%s: Error running Init: %s", stack.Name, err)
			return "Error:", err
		}

	} else {

		err := tf.Init(context.Background(), tfexec.Upgrade(false))
		if err != nil {
			log.Printf("%s: Error running Init: %s", stack.Name, err)
			return "Error:", err
		}
	}

	// Create TF Plan options
	tfplanPath := workspace + stack.Path + "/tfplan-" + stack.Name
	stackPlanOut := tfexec.Out(tfplanPath)

	if len(stack.TFvars) > 0 {
		plan, err := tf.Plan(context.Background(), stackPlanOut, tfexec.VarFile(stack.TFvars))
		if err != nil {
			log.Printf("%s: Error running Plan: %s", stack.Name, err)
			return "Error:", err
		}

		response, err = showPlan(plan, tfplanPath, stack.Name, tf)
		if err != nil {
			return "Error:", err
		}

	} else {
		plan, err := tf.Plan(context.Background(), stackPlanOut)
		if err != nil {
			log.Printf("%s: Error running Plan: %s", stack.Name, err)
			return "Error:", err
		}

		response, err = showPlan(plan, tfplanPath, stack.Name, tf)
		if err != nil {
			return "Error:", err
		}
	}

	return response, err
}

func showPlan(plan bool, tfplanPath string, name string, tf *tfexec.Terraform) (string, error) {

	var err error
	if plan {

		state, err := tf.ShowPlanFileRaw(context.Background(), tfplanPath)
		if err != nil {
			log.Printf("%s: Error running Show: %s", name, err)
			return "ERROR:", err
		}

		re := regexp.MustCompile("Plan:.*")
		summary := re.FindString(state)
		log.Printf("%s: CHANGES DETECTED... %s", name, summary)
		return fmt.Sprintf("CHANGES DETECTED... %s", summary), err

	} else {
		log.Printf("%s: No changes. Infrastructure matches the configuration.", name)
		return "No changes. Infrastructure matches the configuration.", err
	}

}
