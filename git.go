package main

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// cloning the repository that contains all terraform stacks
func gitClone(workspace, repoUrl string) {

	_, present := os.LookupEnv("GITHUB_AUTH_TOKEN")
	if !present {
		log.Fatalf("Error: Could not find GITHUB_AUTH_TOKEN, make sure it is exported as an environment variable")
	}

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	log.Printf("Cloning repository %s", repoUrl)
	_, err := git.PlainClone(workspace, false, &git.CloneOptions{

		Auth: &http.BasicAuth{
			Username: "-", // Yes, this can be anything except an empty string
			Password: token,
		},
		URL: repoUrl,
	})

	if err != nil {
		log.Fatalf("error: %s", err)
	}

}

// TODO: A scheduled job for pulling updates from upstream
