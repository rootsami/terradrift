package tfstack

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setupEnv(t *testing.T) {

	want := map[string]string{"AWS_ACCESS_KEY_ID": "",
		"AWS_CONFIG_FILE":                "",
		"AWS_PROFILE":                    "",
		"AWS_SECRET_ACCESS_KEY":          "",
		"AWS_SESSION_TOKEN":              "",
		"AWS_SHARED_CREDENTIALS_FILE":    "",
		"GCLOUD_KEYFILE_JSON":            "",
		"GOOGLE_APPLICATION_CREDENTIALS": "",
		"GOOGLE_CLOUD_KEYFILE_JSON":      "",
		"GOOGLE_CREDENTIALS":             "",
		"PATH":                           os.Getenv("PATH"),
		"TF_DATA_DIR":                    ".terradrift.testStack"}

	got := setupEnv("test-stack", want)

	assert.Equal(t, want, got, "Environment variables do not match expected output")
}
