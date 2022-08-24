package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-co-op/gocron"
)

func scheduler() {

	stacks := configLoader().Stacks
	interval := configLoader().Interval
	for _, s := range stacks {

		job := gocron.NewScheduler(time.UTC)
		job.Every(interval).Seconds().Do(caller, s.Name)

		job.StartAsync()

	}

}

func caller(name string) {
	hostname := configLoader().Server.Hostname
	port := configLoader().Server.Port
	protocol := configLoader().Server.Protocol
	url := protocol + "://" + hostname + ":" + port + "/api/plan?stack=" + name

	_, err := http.Get(url)
	if err != nil {
		log.Error(err)
	}

}
