package main

import (
	"os"

	"github.com/gin-gonic/gin"
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

var workspace, response string

func main() {

	pwd, _ := os.Getwd()
	workspace = pwd + "/workspace/"

	gitClone(workspace, configLoader().Repository)

	route := gin.Default()

	route.GET("/api/plan", scanHandler)
	route.Run(":8080")
}
