package api

import (
	_ "path/filepath"

	"github.com/feduxorg/letter-generator/letter_generator"
	git "github.com/libgit2/git2go/v34"
	log "github.com/sirupsen/logrus"
)

type Initializer struct{}

func (i *Initializer) Init(dir string, config letter_generator.Config) error {
	tmplSource := config.TemplateSource

	log.WithFields(log.Fields{
		"template_source": tmplSource,
		"destination":     dir,
	}).Debug("Starting to clone config repository")

	repo, err := git.Clone(tmplSource, dir, &git.CloneOptions{})

	if err != nil {
		log.WithFields(log.Fields{
			"source":      tmplSource,
			"destination": dir,
			"msg":         err.Error(),
			"status":      "failure",
		}).Debug("Cloning config repository")

		return err
	}

	remote := "origin"
	repo.Remotes.Delete(remote)

	if err != nil {
		log.WithFields(log.Fields{
			"repository": repo,
			"remote":     remote,
			"msg":        err.Error(),
			"status":     "failure",
		}).Warn("Removing remote failed")
	}

	log.WithFields(log.Fields{
		"source":      tmplSource,
		"destination": dir,
		"status":      "success",
	}).Debug("Cloning config repository")

	return nil
}
