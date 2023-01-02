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
		// obtain PATH from the environment
		"PATH": os.Getenv("PATH"),

		// Terraform
		"TF_DATA_DIR": ".terradrift." + stackName,

		// AWS
		// https://registry.terraform.io/providers/hashicorp/aws/latest/docs#environment-variables
		"AWS_ACCESS_KEY_ID":           os.Getenv("AWS_ACCESS_KEY_ID"),
		"AWS_SECRET_ACCESS_KEY":       os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"AWS_SESSION_TOKEN":           os.Getenv("AWS_SESSION_TOKEN"),
		"AWS_PROFILE":                 os.Getenv("AWS_PROFILE"),
		"AWS_CONFIG_FILE":             os.Getenv("AWS_CONFIG_FILE"),
		"AWS_SHARED_CREDENTIALS_FILE": os.Getenv("AWS_SHARED_CREDENTIALS_FILE"),

		// GCP
		// https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#full-reference
		"GOOGLE_APPLICATION_CREDENTIALS": os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		"GOOGLE_CREDENTIALS":             os.Getenv("GOOGLE_CREDENTIALS"),
		"GOOGLE_CLOUD_KEYFILE_JSON":      os.Getenv("GOOGLE_CLOUD_KEYFILE_JSON"),
		"GCLOUD_KEYFILE_JSON":            os.Getenv("GCLOUD_KEYFILE_JSON"),
	}

	// User provided environment variables
	for key, value := range extraBackendVars {
		tfEnvVars[key] = value
	}

	return tfEnvVars
}
