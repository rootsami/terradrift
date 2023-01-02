package config

import (
	"io/ioutil"

	"testing"

	"gopkg.in/yaml.v2"
)

func TestConfigLoader(t *testing.T) {
	// Loading configuration file for repository and stack properties
	stackConfig, err := ioutil.ReadFile("../../examples/config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(stackConfig, &config)
	if err != nil {
		t.Fatal(err)
	}

	err = ConfigValidator("../../examples/", &config)
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigGenerator(t *testing.T) {

	config, err := ConfigGenerator("../../examples/")
	if err != nil {
		t.Fatal(err)
	}

	err = ConfigValidator("../../examples/", config)
	if err != nil {
		t.Fatal(err)
	}

}
