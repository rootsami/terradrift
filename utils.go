package main

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

func configLoader() Config {

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
	return config
}

// stackExists checks if the requested stack exists in the configration file
func stackExists(name string, stacks []Stack) (stack Stack, result bool) {
	result = false
	for _, stack := range stacks {
		if stack.Name == name {
			result = true
			return stack, result
		}
	}
	return stack, result
}
