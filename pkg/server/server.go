package server

import (
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/rootsami/terradrift/pkg/git"
	"github.com/rootsami/terradrift/pkg/metric"
	"github.com/rootsami/terradrift/pkg/schedulers"

	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

var promMetrics *metric.Metrics

type Server struct {
	Workdir          string
	GitToken         string
	GitTimeout       int
	ConfigPath       string
	ExtraBackendVars map[string]string
	Interval         int
	Repository       string
	Scheme           string
	Hostname         string
	Port             string
	Debug            bool
}

func (s Server) Start() error {

	// Enable debug mode if debug flag is set
	if s.Debug {
		gin.SetMode(gin.DebugMode)
		log.SetLevel(log.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	err := git.GitClone(s.Workdir, s.GitToken, s.Repository, s.GitTimeout)
	if err != nil {
		return err
	}

	// Initialize a non-global prometheus registry
	registery := prometheus.NewRegistry()
	promMetrics = metric.NewMetrics(registery)

	host := s.Scheme + "://" + s.Hostname + ":" + s.Port
	route := gin.Default()
	route.GET("/api/plan", s.scanHandler)
	route.GET("/api/sync", s.gitHandler)
	route.GET("/metrics", prometheusHandler(registery))

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		err := route.Run(":" + s.Port)
		if err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	go func() {
		err := schedulers.ScanScheduler(host, s.Workdir, s.ConfigPath, s.Interval)
		if err != nil {
			log.Error(err)
		}
		wg.Done()
	}()

	go func() {
		err := schedulers.PullScheduler(host, s.Interval)
		if err != nil {
			log.Error(err)
		}
		wg.Done()
	}()
	wg.Wait()

	return err
}
