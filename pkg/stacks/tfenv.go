package stacks

import (
	"os"
)

// setupEnv sets all required environment variable for each tf execution
func setupEnv(stackName string, extraBackendVars map[string]string) (tfenv map[string]string) {

	// Add TF environment variable and append any additional vars provided from the flags
	// During the initialization, .terrafom directory collides with other environments' .terraform
	// causing lots of issues with local terraform.tfstate while performing terraform plan
	// solution would be export TF_DATA_DIR with customized .terraform naming to avoid the issue

	tfEnvVars := map[string]string{
		"TF_DATA_DIR": ".terraform." + stackName,
		"PATH":        os.Getenv("PATH"),
	}

	// User provided environment variables
	for key, value := range extraBackendVars {
		tfEnvVars[key] = value
	}

	return tfEnvVars
}
