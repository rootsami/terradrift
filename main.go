package main

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

	log.SetFormatter(&log.TextFormatter{
		// DisableColors: true,
		FullTimestamp: true,
		// TimestampFormat: "%YYYY/%MM/%DD - %HH:%MM:%SS",
		TimestampFormat: "2006/01/02 - 15:04:05",
	})

	pwd, _ := os.Getwd()
	workspace = pwd + "/workspace/"

	gitClone(workspace, configLoader().Repository)

	route := gin.Default()

	route.GET("/api/plan", scanHandler)
	route.Run(":8080")
}
