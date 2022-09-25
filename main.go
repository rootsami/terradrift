package main

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rootsami/terradrift/pkg/git"
	"github.com/rootsami/terradrift/pkg/schedulers"
	"github.com/rootsami/terradrift/pkg/stacks"
	log "github.com/sirupsen/logrus"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app        = kingpin.New("terradrift", "A tool to detect drifts in terraform IaC")
	hostname   = app.Flag("hostname", "hostname that apil will be exposed.").Default("localhost").String()
	port       = app.Flag("port", "port of the service api is listening on").Default("8080").String()
	protocol   = app.Flag("protocol", "The protocol of exposed endpoint http/https").Default("http").String()
	repository = app.Flag("repository", "The git repository which include all terraform stacks ").Required().String()
	gitToken   = app.Flag("git-token", "Personal access token to access git repositories").Required().String()
	interval   = app.Flag("interval", "The interval for scan scheduler").Default("60").Int()
	configPath = app.Flag("config", "Path for configuration file holding the stack information").Default("config.yaml").String()
	workspace  string
)

func init() {
	dir, err := ioutil.TempDir("", "terradrift")
	if err != nil {
		log.Fatal(err)
	}
	workspace = dir + "/"

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 - 15:04:05",
	})

}

func main() {

	kingpin.MustParse(app.Parse(os.Args[1:]))

	git.GitClone(workspace, *gitToken, *repository)

	route := gin.Default()
	route.GET("/api/plan", scanHandler)

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		route.Run(":" + *port)
		wg.Done()
	}()

	go func() {
		schedulers.ScanScheduler(*hostname, *port, *protocol, *configPath, *interval)
		wg.Done()
	}()

	go func() {
		schedulers.PullScheduler(workspace, *gitToken, *interval)
		wg.Done()
	}()
	wg.Wait()
}

func scanHandler(c *gin.Context) {

	name := c.Query("stack")
	planResp, err := stacks.StackScan(name, workspace, *configPath)

	if err == nil {

		c.JSON(200, planResp)
	} else {

		errorMessage := error.Error(err)
		if errorMessage == "ERROR: STACK WAS NOT FOUND" {

			// Given stack name was not found in the configuration
			c.JSON(404, errorMessage)
		} else if strings.Contains(errorMessage, "error acquiring the state lock") {

			// When there's a current terrafom plan in progress, terraform locks the state till it's finished.
			c.JSON(502, "Another plan is in-progress for the requested stack, please try again in few minutes.")

		} else {

			c.JSON(500, errorMessage)
		}
	}
}
