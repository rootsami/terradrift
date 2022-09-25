package schedulers

import (
	"net/http"
	"time"

	"github.com/rootsami/terradrift/pkg/config"
	"github.com/rootsami/terradrift/pkg/git"
	log "github.com/sirupsen/logrus"

	"github.com/go-co-op/gocron"
)

func PullScheduler(workspace, token string, interval int) {
	job := gocron.NewScheduler(time.UTC)
	job.Every(interval).Seconds().Do(git.GitPull, workspace, token)

	job.StartAsync()
}

func ScanScheduler(hostname, port, protocol, configPath string, interval int) {

	stacks := config.ConfigLoader(configPath).Stacks
	for _, s := range stacks {

		url := protocol + "://" + hostname + ":" + port + "/api/plan?stack=" + s.Name
		job := gocron.NewScheduler(time.UTC)
		job.Every(interval).Seconds().Do(apiCaller, url)

		job.StartAsync()

	}
}

func apiCaller(url string) {

	_, err := http.Get(url)
	if err != nil {
		log.Error(err)
	}

}
