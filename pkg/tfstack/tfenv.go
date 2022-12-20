package tfstack

import (
	"os"
)

// setupEnv sets all required environment variable for each tf execution
func setupEnv(stackName string, extraBackendVars map[string]string) (tfenv map[string]string) {

	// .terraform is a directory that Terraform uses to store internal data, including the
	// current state, configuration files, and any data from providers.
	// https://www.terraform.io/docs/commands/environment-variables.html#tf_data_dir
	// adding TF environment variable and append any additional vars provided from the flags
	// during the initialization, .terraform directory collides with other environments' .terraform
	// causing many issues with local terraform.tfstate while performing terraform plan.
	// TF_DATA_DIR with customized .terraform naming solves this issue.
	// rename .terraform directory to .terradrift.{stackName} for easy detection and cleanup

	tfEnvVars := map[string]string{
		"TF_DATA_DIR": ".terradrift." + stackName,
		"PATH":        os.Getenv("PATH"),
	}

	// User provided environment variables
	for key, value := range extraBackendVars {
		tfEnvVars[key] = value
	}

	return tfEnvVars
}
