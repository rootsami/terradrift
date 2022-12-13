package schedulers

import (
	"net/http"
	"time"

	"github.com/rootsami/terradrift/pkg/config"
	log "github.com/sirupsen/logrus"

	"github.com/go-co-op/gocron"
)

// The whole idea behind the schedular is to invoke all necessary api calls to the running server.
// It doesn't have to be tightly coupled with the server
// It always has to have the ability to be taken outside of the server to be ran seperately
// in case we decided to move the schedulars to a seperate app

func PullScheduler(host string, interval int) error {

	url := host + "/api/sync"
	job := gocron.NewScheduler(time.UTC)
	_, err := job.Every(interval).Seconds().Do(apiCaller, url)
	if err != nil {
		return err
	}

	job.StartAsync()

	return nil
}

func ScanScheduler(host, configPath string, interval int) error {

	stacks, err := config.ConfigLoader(configPath)
	if err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}

	for _, s := range stacks.Stacks {

		url := host + "/api/plan?stack=" + s.Name
		job := gocron.NewScheduler(time.UTC)
		_, err := job.Every(interval).Seconds().Do(apiCaller, url)
		if err != nil {
			log.WithFields(log.Fields{"stack": s.Name}).Error(err)
			return err
		}

		job.StartAsync()
	}
	return nil
}

func apiCaller(url string) error {

	_, err := http.Get(url)
	if err != nil {
		return err
	}

	return nil

}
