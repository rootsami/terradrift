package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Stacks []Stack `yaml:"stacks"`
}

type Stack struct {
	Name    string `yaml:"name" json:"name"`
	Path    string `yaml:"path" json:"path"`
	TFvars  string `yaml:"tfvars,omitempty" json:"tfvars,omitempty"`
	Backend string `yaml:"backend,omitempty" json:"backend,omitempty"`
}

func ConfigLoader(path string) (*Config, error) {

	// Loading configuration file for repository and stack properties
	// TODO: config validator
	stackConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(stackConfig, &config)
	if err != nil {
		return nil, err
	}

	return &config, err
}
