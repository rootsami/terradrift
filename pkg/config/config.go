package config

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Stacks []Stack `yaml:"stacks"`
}

type Stack struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Path    string `yaml:"path"`
	TFvars  string `yaml:"tfvars"`
	Backend string `yaml:"backend"`
}

func ConfigLoader(path string) Config {

	// Loading configuration file for repository and stack properties
	// TODO: config validator
	stackConfig, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(stackConfig, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
