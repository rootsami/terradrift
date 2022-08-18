package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// cloning the repository that contains all terraform stacks
func gitClone(workspace, repoUrl string) {

	_, present := os.LookupEnv("GITHUB_AUTH_TOKEN")
	if !present {
		log.Fatalf("ERROR: Could not find GITHUB_AUTH_TOKEN, make sure it is exported as an environment variable")
	}

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	log.Infof("Cloning repository %s", repoUrl)
	_, err := git.PlainClone(workspace, false, &git.CloneOptions{

		Auth: &http.BasicAuth{
			Username: "-", // Yes, this can be anything except an empty string
			Password: token,
		},
		URL: repoUrl,
		// Depth: 1, // disabled for now because gitPull fails after shallow clone
	})

	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

}

func gitPull(workspace string) {

	r, err := git.PlainOpen(workspace)
	if err != nil {
		log.Errorf("ERROR: PULL FAILED %s", err)
	}

	w, err := r.Worktree()
	if err != nil {
		log.Errorf("ERROR: PULL FAILED %s", err)
	}

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	err = w.Pull(&git.PullOptions{
		Auth: &http.BasicAuth{
			Username: "-", // Yes, this can be anything except an empty string
			Password: token,
		},
		Force: true,
	})
	if err != nil {
		log.Warnf("Pulling latest updates: %s", err)
	}

}
