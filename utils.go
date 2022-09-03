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
	err = yaml.Unmarshal(stackConfig, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
