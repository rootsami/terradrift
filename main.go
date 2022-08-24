package main

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Server     Server  `yaml:"server"`
	Interval   int     `yaml:"interval"`
	Repository string  `yaml:"repository"`
	Stacks     []Stack `yaml:"stacks"`
}

type Server struct {
	Hostname string `yaml:"hostname"`
	Port     string `yaml:"port"`
	Protocol string `yaml:"protocol"`
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
		FullTimestamp: true,
		// TimestampFormat: "%YYYY/%MM/%DD - %HH:%MM:%SS",
		TimestampFormat: "2006/01/02 - 15:04:05",
	})

	pwd, _ := os.Getwd()
	workspace = pwd + "/workspace/"

	gitClone(workspace, configLoader().Repository)

	scheduler()

	route := gin.Default()

	route.GET("/api/plan", scanHandler)
	route.Run(":" + configLoader().Server.Port)
}
