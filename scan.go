package main

import (
	"context"
	"log"
	"os"
	"regexp"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

// func stackScan(stackName, stackVersion, stackWorkspace string) {
func stackScan(stack Stack) {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion(stack.Version)),
	}

	log.Printf("%s: Installing Terraform %s ...", stack.Name, stack.Version)
	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("%s: Error installing Terraform: %s", stack.Name, err)
	}

	tf, err := tfexec.NewTerraform(workspace+stack.Path, execPath)
	if err != nil {
		log.Fatalf("%s: Error running NewTerraform: %s", stack.Name, err)
	}

	// Stacks come with two different structures:
	// 1. All resources for multiple stacks (environments) exist in one directory and backend initialization is done with environments/<name>.hcl
	// 2. Regular stack where all resources, tfvars and backend configs are in the same directory
	log.Printf("%s: Initializing Terraform...", stack.Name)
	if len(stack.Backend) > 0 {

		// during the initialization, .terrafom directory collides with other environments' .terraform
		// causing lots of issue with local terraform.tfstate while performing terraform plan
		// solution would be export TF_DATA_DIR with customized .terraform naming to avoid the issue
		os.Setenv("TF_DATA_DIR", ".terraform."+stack.Name)
		err = tf.Init(context.Background(), tfexec.Upgrade(false), tfexec.BackendConfig(stack.Backend))
		if err != nil {
			log.Fatalf("%s: Error running Init: %s", stack.Name, err)
		}

	} else {

		err = tf.Init(context.Background(), tfexec.Upgrade(false))
		if err != nil {
			log.Fatalf("%s: Error running Init: %s", stack.Name, err)
		}
	}

	// TODO: stackPlan needs to run as a background scheduled job, a for loop is placed here just to prove the iteration works without downloading and initialzing with each run.
	// for i := 1; i <= 1; i++ {
	// 	result := stackPlan(stack, tf)
	// }
	stackPlan(stack, tf)

}

// stackPlan has been separated to be called indivisually with the schedule to avoid downloading/installing
// the required terraform version and initializing everytime the scan runs
func stackPlan(stack Stack, tf *tfexec.Terraform) {

	// Create TF Plan options
	tfplanPath := workspace + stack.Path + "/tfplan-" + stack.Name
	stackPlanOut := tfexec.Out(tfplanPath)

	if len(stack.TFvars) > 0 {
		plan, err := tf.Plan(context.Background(), stackPlanOut, tfexec.VarFile(stack.TFvars))
		if err != nil {
			log.Fatalf("%s: Error running Plan: %s", stack.Name, err)
		}

		showPlan(plan, tfplanPath, stack.Name, tf)

	} else {
		plan, err := tf.Plan(context.Background(), stackPlanOut)
		if err != nil {
			log.Fatalf("%s: Error running Plan: %s", stack.Name, err)
		}

		showPlan(plan, tfplanPath, stack.Name, tf)
	}
}

func showPlan(plan bool, tfplanPath string, name string, tf *tfexec.Terraform) {

	if plan {

		state, err := tf.ShowPlanFileRaw(context.Background(), tfplanPath)
		if err != nil {
			log.Fatalf("%s: Error running Show: %s", name, err)
		}

		re := regexp.MustCompile("Plan:.*")
		summary := re.FindString(state)
		log.Printf("%s: CHANGES DETECTED... %s", name, summary)

	} else {
		log.Printf("%s: No changes. Infrastructure matches the configuration.", name)
	}

}
