package main

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rootsami/terradrift/pkg/git"
	"github.com/rootsami/terradrift/pkg/metric"
	"github.com/rootsami/terradrift/pkg/schedulers"
	"github.com/rootsami/terradrift/pkg/stacks"
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app              = kingpin.New("terradrift", "A tool to detect drifts in terraform IaC")
	hostname         = app.Flag("hostname", "hostname that apil will be exposed.").Default("localhost").String()
	port             = app.Flag("port", "port of the service api is listening on").Default("8080").String()
	scheme           = app.Flag("scheme", "The scheme of exposed endpoint http/https").Default("http").String()
	repository       = app.Flag("repository", "The git repository which include all terraform stacks ").Required().String()
	gitToken         = app.Flag("git-token", "Personal access token to access git repositories").Required().String()
	gitTimeout       = app.Flag("git-timeout", "Wait timeout for git repoistory to clone or pull updates").Default("120").Int()
	interval         = app.Flag("interval", "The interval for scan scheduler").Default("60").Int()
	configPath       = app.Flag("config", "Path for configuration file holding the stack information").Default("config.yaml").String()
	extraBackendVars = app.Flag("extra-backend-vars", "Extra backend environment variables ex. GOOGLE_CREDENTIALS OR AWS_ACCESS_KEY").StringMap()
	workspace        string
	promMetrics      *metric.Metrics
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

	err := git.GitClone(workspace, *gitToken, *repository, *gitTimeout)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize a non-global prometheus registry
	registery := prometheus.NewRegistry()
	promMetrics = metric.NewMetrics(registery)

	host := *scheme + "://" + *hostname + ":" + *port
	route := gin.Default()
	route.GET("/api/plan", scanHandler)
	route.GET("/api/sync", gitHandler)
	route.GET("/metrics", prometheusHandler(registery))

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		err := route.Run(":" + *port)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	go func() {
		err := schedulers.ScanScheduler(host, *configPath, *interval)
		if err != nil {
			log.Error(err)
		}
		wg.Done()
	}()

	go func() {
		err := schedulers.PullScheduler(host, *interval)
		if err != nil {
			log.Error(err)
		}
		wg.Done()
	}()
	wg.Wait()
}

// scanHandler is a handler function for scan endpoint and record metrics
// for changed resources based on the scan plan result
func scanHandler(c *gin.Context) {

	name := c.Query("stack")
	planResp, err := stacks.StackScan(name, workspace, *configPath, *extraBackendVars)

	if err == nil {

		// Record metrics for drifts in resources
		promMetrics.AddResources.With(prometheus.Labels{"stack": name}).Set(float64(planResp.Add))
		promMetrics.ChangeResources.With(prometheus.Labels{"stack": name}).Set(float64(planResp.Change))
		promMetrics.DestroyResources.With(prometheus.Labels{"stack": name}).Set(float64(planResp.Destroy))

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

// gitHandler is a handler function for git sync endpoint
func gitHandler(c *gin.Context) {

	status, err := git.GitPull(workspace, *gitToken, *gitTimeout)
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, status)
	}
}

// prometheusHandler returns a gin.HandlerFunc that serves prometheus metrics.
func prometheusHandler(reg *prometheus.Registry) gin.HandlerFunc {
	h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
