package stacks

import (
	"context"
	"encoding/json"
	"fmt"
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

// The main function for installing the exact version of the stack, initiate and run terraform plan
func StackScan(name, workspace, configPath string, extraBackendVars map[string]string) (string, error) {

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

		tfenv := setupEnv(stack.Name, extraBackendVars)
		tf.SetEnv(tfenv)

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

	var response string

	// Stacks come with two different structures:
	// 1. All resources for multiple stacks (environments) exist in one directory and backend initialization is done with environments/<name>.hcl
	// 2. Regular stack where all resources, tfvars and backend configs are in the same directory

	log.WithFields(log.Fields{"stack": stack.Name, "version": stack.Version}).Info("Initializing Terraform...")

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

	if plan {

		state, err := tf.ShowPlanFileRaw(context.Background(), planFile)
		if err != nil {
			log.WithFields(log.Fields{"stack": name}).Errorf("Running show: %s", err)
			return "", err
		}

		rawSummary := driftCalculator(state)

		jsonSummary, err := json.Marshal(rawSummary)
		if err != nil {
			log.WithFields(log.Fields{"stack": name}).Error(err)
		}

		log.WithFields(log.Fields{"stack": name}).Info(string(jsonSummary))
		return (string(jsonSummary)), err

	} else {

		rawSummary := &DriftSum{
			Drift:   false,
			Add:     0,
			Change:  0,
			Destroy: 0,
		}

		jsonSummary, err := json.Marshal(rawSummary)
		if err != nil {
			log.WithFields(log.Fields{"stack": name}).Error(err)
		}

		log.WithFields(log.Fields{"stack": name}).Info(string(jsonSummary))
		return (string(jsonSummary)), err
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
func driftCalculator(state string) DriftSum {

	re := regexp.MustCompile("Plan:[^0-9]*(?P<add>[0-9])[^0-9]*(?P<change>[0-9])[^0-9]*(?P<destroy>[0-9])")
	matches := re.FindStringSubmatch(state)

	addIndex := re.SubexpIndex("add")
	add, err := strconv.Atoi(matches[addIndex])
	if err != nil {
		log.Error(err)
	}

	changeIndex := re.SubexpIndex("change")
	change, err := strconv.Atoi(matches[changeIndex])
	if err != nil {
		log.Error(err)
	}

	destroyIndex := re.SubexpIndex("destroy")
	destroy, err := strconv.Atoi(matches[destroyIndex])
	if err != nil {
		log.Error(err)
	}

	DriftSum := &DriftSum{
		Drift:   true,
		Add:     add,
		Change:  change,
		Destroy: destroy,
	}

	return *DriftSum
}
