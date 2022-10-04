package git

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// cloning the repository that contains all terraform stacks
func GitClone(workspace, token, repoUrl string) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(120)*time.Second)
	defer cancel()
	log.Infof("Cloning repository %s", repoUrl)
	_, err := git.PlainCloneContext(ctx, workspace, false, &git.CloneOptions{

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

func GitPull(workspace, token string) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(120)*time.Second)
	defer cancel()

	r, err := git.PlainOpen(workspace)
	if err != nil {
		log.Errorf("ERROR: PULL FAILED %s", err)
	}

	w, err := r.Worktree()
	if err != nil {
		log.Errorf("ERROR: PULL FAILED %s", err)
	}

	err = w.PullContext(ctx, &git.PullOptions{
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
