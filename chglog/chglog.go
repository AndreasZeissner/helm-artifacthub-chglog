package chglog

import (
	"log"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

const (
	StopIteration = "StopIteration"
)

func filterSubirs(directories []string) func(string) bool {
	return func(path string) bool {
		var ok = false
		for _, dir := range directories {
			if strings.HasPrefix(path, dir) {
				ok = true
				break
			}
		}
		return ok
	}
}
func GenerateChangelogForRepo(from, to, repoURL string, subdirectories []string) []*ArtifactHubChangelogObject {
	chglogs := []*ArtifactHubChangelogObject{}
	repo := OpenRepo(repoURL)

	_, fromCommit := ResolveTag(repo, from)
	_, toCommit := ResolveTag(repo, to)

	options := &git.LogOptions{
		From: toCommit.Hash,
		All:  false,
	}

	if len(subdirectories) > 0 {
		options = &git.LogOptions{
			From:       toCommit.Hash,
			All:        false,
			PathFilter: filterSubirs(subdirectories),
		}
	}
	// Initialize a map to store commits reachable from the fromCommit
	visited := make(map[plumbing.Hash]bool)

	// Walk from `fromCommit` to mark all reachable commits
	err := walkCommits(repo, fromCommit, visited)
	if err != nil {
		log.Fatalf("Failed to walk commits: %v", err)
	}

	// Now iterate from `toCommit` and only print commits not in `visited`
	iter, err := repo.Log(options)
	if err != nil {
		log.Fatalf("Log commits: %v", err)
	}

	defer iter.Close()
	err = iter.ForEach(func(c *object.Commit) error {
		// Stop if we reach the `fromCommit`
		if visited[c.Hash] {
			return nil
		}

		resolver := NewConventionalCommitsResolver(c)
		chglog, err := resolver.ResolveChangelogEntry()
		if err != nil {
			return nil
		}
		chglogs = append(chglogs, chglog)

		return nil
	})

	return chglogs
}

// walkCommits walks through the commit history starting from `commit` and marks all reachable commits
func walkCommits(repo *git.Repository, commit *object.Commit, visited map[plumbing.Hash]bool) error {
	commitIter := object.NewCommitPreorderIter(commit, nil, nil)
	err := commitIter.ForEach(func(c *object.Commit) error {
		visited[c.Hash] = true
		return nil
	})
	return err
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
