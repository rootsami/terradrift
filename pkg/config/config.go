package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Stacks []Stack `yaml:"stacks" json:"stacks"`
}

type Stack struct {
	Name    string `yaml:"name" json:"name"`
	Path    string `yaml:"path" json:"path"`
	TFvars  string `yaml:"tfvars,omitempty" json:"tfvars,omitempty"`
	Backend string `yaml:"backend,omitempty" json:"backend,omitempty"`
}

// ConfigLoader loads the configuration file and returns a Config struct
func ConfigLoader(workdir, configPath string) (*Config, error) {

	stackConfig, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(stackConfig, &config)
	if err != nil {
		return nil, err
	}

	err = ConfigValidator(workdir, &config)

	return &config, err
}

// ConfigValidator validates the configuration file and returns an error if TFvars or Backend file does not exist
func ConfigValidator(workdir string, cfg *Config) error {
	for _, stack := range cfg.Stacks {
		if stack.TFvars != "" {
			tfvarPath := workdir + stack.Path + "/" + stack.TFvars
			if _, err := os.Stat(tfvarPath); os.IsNotExist(err) {
				return err
			}
		}
		if stack.Backend != "" {
			backendPath := workdir + stack.Path + "/" + stack.Backend
			if _, err := os.Stat(backendPath); os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}
