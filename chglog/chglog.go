package chglog

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

const (
	StopIteration = "StopIteration"
)

func GenerateChangelogForRepo(from, to, repoURL string, subdirectories []string) []*ArtifactHubChangelogObject {
	chglogs := []*ArtifactHubChangelogObject{}
	repo := OpenRepo(repoURL)

	_, fromCommit := ResolveTag(repo, from)
	_, toCommit := ResolveTag(repo, to)

	commitIter, err := repo.Log(&git.LogOptions{
		From: fromCommit.Hash,
	})
	if err != nil {
		log.Fatalf("Failed to get commit iterator: %v", err)
	}

	err = commitIter.ForEach(func(c *object.Commit) error {
		if c.Hash == toCommit.Hash {
			return fmt.Errorf(StopIteration)
		}
		if len(subdirectories) > 0 {
			parent, err := c.Parent(0)
			if err != nil {
				return err
			}

			patch, err := c.Patch(parent)
			if err != nil {
				return err
			}
			for _, fileStat := range patch.Stats() {
				for _, subdir := range subdirectories {
					if strings.HasPrefix(fileStat.Name, subdir) {
						resolver := NewConventionalCommitsResolver(c)
						chglog, err := resolver.ResolveChangelogEntry()
						if err != nil {
							return nil
						}
						chglogs = append(chglogs, chglog)
						return nil
					}
				}
			}
		} else {
			resolver := NewConventionalCommitsResolver(c)
			chglog, err := resolver.ResolveChangelogEntry()
			if err != nil {
				return nil
			}
			chglogs = append(chglogs, chglog)
		}

		return nil
	})
	if err != nil && err.Error() != StopIteration {
		log.Fatalf("Error iterating commits: %v", err)
	}

	return chglogs
}

func NewCli() cli.ActionFunc {
	return func(cctx *cli.Context) error {
		from := cctx.String("from")
		to := cctx.String("to")
		repoURL := cctx.String("repoURL")
		output := cctx.String("output")
		subdirectories := cctx.StringSlice("paths")

		chglogs := GenerateChangelogForRepo(from, to, repoURL, subdirectories)

		file, err := os.Create(output)
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
		defer file.Close()
		encoder := yaml.NewEncoder(file)
		defer encoder.Close()

		err = encoder.Encode(chglogs)
		if err != nil {
			log.Fatalf("Failed to encode data to YAML: %v", err)
		}

		return nil
	}
}
