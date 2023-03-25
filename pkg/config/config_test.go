package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigLoader(t *testing.T) {
	// Loading configuration file for repository and stack properties

	staticConfig := &Config{
		Stacks: []Stack{
			{
				Name: "core-staging",
				Path: "gcp/core-staging"},
			{
				Name:    "core-production",
				Path:    "aws/core-production",
				TFvars:  "",
				Backend: ""},
			{
				Name:    "api-staging",
				Path:    "gcp/api",
				TFvars:  "environments/staging.tfvars",
				Backend: "environments/staging.hcl"},
			{
				Name:    "api-production",
				Path:    "gcp/api",
				TFvars:  "environments/production.tfvars",
				Backend: "environments/production.hcl"},
		},
	}
	want := staticConfig
	got, err := ConfigLoader("../../examples/", "../../examples/config.yaml")

	assert.NoError(t, err, "Unexpected error from ConfigLoader")
	assert.Equal(t, want, got, "Config does not match expected output")

}

func TestConfigGenerator(t *testing.T) {

	config, err := ConfigGenerator("../../examples/")
	want := &Config{
		Stacks: []Stack{
			{
				Name:    "aws-core-production",
				Path:    "aws/core-production",
				TFvars:  "",
				Backend: "",
			},

			{
				Name:    "gcp-api-environments-production",
				Path:    "gcp/api",
				TFvars:  "environments/production.tfvars",
				Backend: "environments/production.hcl",
			},

			{
				Name:    "gcp-api-environments-staging",
				Path:    "gcp/api",
				TFvars:  "environments/staging.tfvars",
				Backend: "environments/staging.hcl",
			},

			{
				Name:    "gcp-core-staging",
				Path:    "gcp/core-staging",
				TFvars:  "",
				Backend: "",
			},
		},
	}
	got := config

	assert.NoError(t, err, "Unexpected error from ConfigGenerator")
	assert.Equal(t, want, got, "Config does not match expected output")

}
