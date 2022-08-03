package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Repository string  `yaml:"repository"`
	Stacks     []Stack `yaml:"stacks"`
}

type Stack struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Path    string `yaml:"path"`
	TFvars  string `yaml:"tfvars"`
	Backend string `yaml:"backend"`
}

var workspace string

func main() {

	pwd, _ := os.Getwd()
	workspace = pwd + "/workspace/"

	// Loading configuration file for repository and stack properties
	// TODO: config validator
	stackConfig, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	loadConfig := yaml.Unmarshal(stackConfig, &config)
	if loadConfig != nil {
		log.Fatal(loadConfig)
	}

	gitClone(workspace, config.Repository)

	for _, s := range config.Stacks {
		stackScan(s)

	}
}
