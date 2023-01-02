package git

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// cloning the repository that contains all terraform stacks
func GitClone(workdir, token, repoUrl string, timeout int) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	log.Infof("cloning repository %s", repoUrl)
	_, err := git.PlainCloneContext(ctx, workdir, false, &git.CloneOptions{

		Auth: &http.BasicAuth{
			Username: "-", // Yes, this can be anything except an empty string
			Password: token,
		},
		URL: repoUrl,
	})

	if err != nil {
		return err
	}

	return nil
}

func GitPull(workdir, token string, timeout int) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	log.Info("pulling latest updates from upstream")

	r, err := git.PlainOpen(workdir)
	if err != nil {
		log.Errorf("pull failed %s", err)
		return "", err
	}

	w, err := r.Worktree()
	if err != nil {
		log.Errorf("pull failed %s", err)
		return "", err
	}

	err = w.PullContext(ctx, &git.PullOptions{
		Auth: &http.BasicAuth{
			Username: "-", // Yes, this can be anything except an empty string
			Password: token,
		},
		Force: true,
	})

	if err == git.NoErrAlreadyUpToDate {
		log.Infof("pulling latest updates: %s", err)
		return error.Error(err), nil

	} else if err != nil {
		log.Errorf("pull failed %s", err)
		return "", err

	} else {
		log.Info("pulling latest updates: Success!")
		return "Success!", nil
	}
}
