package stacks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/rootsami/terradrift/pkg/config"
	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-exec/tfexec"
)

type DriftSum struct {
	Drift   bool `json:"drift"`
	Add     int  `json:"add"`
	Change  int  `json:"change"`
	Destroy int  `json:"destroy"`
}

// StackScan scans a given stack only if stack exist in the list of stacks in the config file
// and returns a DriftSum that describes the drift between the stack's Terraform state
// and the state of its resources.
//
// name is the name of the stack to scan.
// workspace is the workspace of the stack to scan.
// configPath is the path to the config of the stacks to scan.
// extraBackendVars is a map of extra variables to pass to the backend when
// initializing the stack.
func StackScan(name, workspace, configPath string, extraBackendVars map[string]string) (*DriftSum, error) {

	config, err := config.ConfigLoader(configPath)
	if err != nil {
		return nil, err
	}

	stack, validStack := stackExists(name, config.Stacks)
	if validStack {

		response, _, err := StackInit(workspace, stack, extraBackendVars)
		if err != nil {
			log.WithFields(log.Fields{"stack": stack.Name}).Error(err)
			return nil, err
		}
		log.WithFields(log.Fields{"stack": stack.Name}).Info(fmt.Sprintf("%+v", *response))
		return response, nil

	} else {
		err := errors.New("ERROR: STACK WAS NOT FOUND")
		log.WithFields(log.Fields{"stack": name}).Error(err)
		return nil, err
	}
}

// StackInit initializes a stack and returns a DriftSum that describes the drift details
func StackInit(workspace string, stack config.Stack, extraBackendVars map[string]string) (*DriftSum, string, error) {

	// The path for terrafom binary
	execPath, tfver, err := install(stack, workspace)
	if err != nil {
		log.WithFields(log.Fields{"stack": stack.Name}).Error(err)
	}

	tf, err := tfexec.NewTerraform(workspace+stack.Path, execPath)
	if err != nil {
		log.WithFields(log.Fields{"stack": stack.Name}).Errorf("Running NewTerraform: %s", err)
		return nil, "", err
	}

	tfenv := setupEnv(stack.Name, extraBackendVars)
	tf.SetEnv(tfenv)

	response, err := stackPlan(workspace, stack, tf)
	if err != nil {
		log.WithFields(log.Fields{"stack": stack.Name}).Error(err)
		return nil, "", err
	}

	return response, tfver, nil

}

// stackPlan executes terraform plan for a given stack and returns a
// DriftSum. The DriftSum is used to determine if any resources have drifted
// from the Terraform state.
func stackPlan(workspace string, stack config.Stack, tf *tfexec.Terraform) (*DriftSum, error) {

	var response *DriftSum

	// Stacks come with two different structures:
	// 1. All resources for multiple stacks (environments) exist in one directory
	//    and backend initialization is done with path/to/backend.hcl
	// 2. Regular stack where all resources, tfvars and backend configs are in the same directory

	log.WithFields(log.Fields{"stack": stack.Name}).Debug("Initializing Terraform...")

	err := tf.Init(context.Background(), tfexec.Upgrade(false), tfexec.BackendConfig(stack.Backend))
	if err != nil {
		log.WithFields(log.Fields{"stack": stack.Name}).Error("Running Init")
		return nil, err
	}

	// Create TF Plan options
	planFile := workspace + stack.Path + "/" + stack.Name + ".plan"
	stackPlanFile := tfexec.Out(planFile)

	if len(stack.TFvars) > 0 {
		plan, err := tf.Plan(context.Background(), stackPlanFile, tfexec.VarFile(stack.TFvars))
		if err != nil {
			return nil, err
		}

		response, err = showPlan(plan, planFile, stack.Name, tf)
		if err != nil {
			return nil, err
		}

	} else {
		plan, err := tf.Plan(context.Background(), stackPlanFile)
		if err != nil {
			return nil, err
		}

		response, err = showPlan(plan, planFile, stack.Name, tf)
		if err != nil {
			return nil, err
		}
	}

	err = cleanUpPlanFile(planFile)
	if err != nil {
		return nil, err
	}

	return response, err
}

// showPlan shows the plan and returns the number of changes
func showPlan(plan bool, planFile string, name string, tf *tfexec.Terraform) (*DriftSum, error) {

	if plan {

		state, err := tf.ShowPlanFileRaw(context.Background(), planFile)
		if err != nil {
			return nil, err
		}

		summary, err := driftCalculator(state)
		if err != nil {
			return nil, err
		}

		return summary, nil

	} else {

		summary := &DriftSum{
			Drift:   false,
			Add:     0,
			Change:  0,
			Destroy: 0,
		}

		return summary, nil
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

// driftCalculator returns a detailed number of changes that was detected in the plan
func driftCalculator(state string) (*DriftSum, error) {

	re := regexp.MustCompile("Plan:[^0-9]*(?P<add>[0-9])[^0-9]*(?P<change>[0-9])[^0-9]*(?P<destroy>[0-9])")
	matches := re.FindStringSubmatch(state)

	addIndex := re.SubexpIndex("add")
	add, err := strconv.Atoi(matches[addIndex])
	if err != nil {
		return nil, err
	}

	changeIndex := re.SubexpIndex("change")
	change, err := strconv.Atoi(matches[changeIndex])
	if err != nil {
		return nil, err
	}

	destroyIndex := re.SubexpIndex("destroy")
	destroy, err := strconv.Atoi(matches[destroyIndex])
	if err != nil {
		return nil, err
	}

	DriftSum := &DriftSum{
		Drift:   true,
		Add:     add,
		Change:  change,
		Destroy: destroy,
	}

	return DriftSum, err
}

// cleanUpPlanFile removes the plan file after the plan has been reported
func cleanUpPlanFile(planFile string) error {

	err := os.Remove(planFile)
	if err != nil {
		return err
	}

	return nil
}
