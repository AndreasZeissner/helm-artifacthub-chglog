package chglog

import (
	"log"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func OpenRepo(repoURL string) *git.Repository {
	var repo *git.Repository
	var err error

	if strings.HasPrefix(repoURL, "https://") {
		repo, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL: repoURL,
		})
		if err != nil {
			panic(err)
		}
	} else {
		repo, err = git.PlainOpen(repoURL)
		if err != nil {
			panic(err)
		}
	}

	return repo
}

func ResolveTag(repo *git.Repository, tag string) (*plumbing.Reference, *object.Commit) {
	t, err := repo.Tag(tag)
	if err != nil {
		log.Fatalf("Failed to resolve tag object: %v", err)
	}

	var c *object.Commit
	tagObj, err := repo.TagObject(t.Hash())
	switch err {
	case plumbing.ErrObjectNotFound:
		// lightweight tag
		c, err = repo.CommitObject(t.Hash())

		if err != nil {
			log.Fatalf("Failed to get commit object: %v", err)
		}
	default:
		// Handle the case of an annotated tag
		if err != nil {
			log.Fatalf("Failed to get tag object: %v", err)
		}
		c, err = tagObj.Commit()
		if err != nil {
			log.Fatalf("Failed to get commit object from tag: %v", err)
		}
	}

	return t, c
}
