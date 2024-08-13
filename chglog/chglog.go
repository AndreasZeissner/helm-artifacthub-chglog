package chglog

import (
	"fmt"
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

func extractTagPrefix(tag string) string {
	if idx := strings.Index(tag, "@"); idx != -1 {
		return tag[:idx]
	}
	return tag // If no "@" is found, return the whole tag (assuming no version in the tag)
}

func findPreviousTaggedCommit(repo *git.Repository, options *git.LogOptions, tagPrefix string) (*object.Commit, string, error) {
	// Get an iterator for the log starting from the given commit
	iter, err := repo.Log(options)
	if err != nil {
		return nil, "", err
	}
	defer iter.Close()

	var foundCommit *object.Commit
	var foundTag string
	skipFirstCommit := true

	// Iterate through the commits
	for {
		commit, err := iter.Next()
		if err != nil {
			if err == plumbing.ErrObjectNotFound {
				break
			}
			return nil, "", err
		}

		// Skip the first commit (which is the toCommit itself)
		if skipFirstCommit {
			skipFirstCommit = false
			continue
		}

		// Check if this commit has a tag directly associated with it
		tags, err := repo.Tags()
		if err != nil {
			return nil, "", err
		}

		err = tags.ForEach(func(ref *plumbing.Reference) error {
			// Ensure the tag name has the correct prefix
			currentTag := ref.Name().Short()
			currentTagPrefix := extractTagPrefix(currentTag)

			if currentTagPrefix != tagPrefix {
				return nil
			}

			tag, err := repo.TagObject(ref.Hash())
			if err != nil && err != plumbing.ErrObjectNotFound {
				return err
			}

			var target plumbing.Hash
			if err == plumbing.ErrObjectNotFound {
				target = ref.Hash()
			} else {
				target = tag.Target
			}

			if target == commit.Hash {
				foundCommit = commit
				foundTag = ref.Name().Short()
				return fmt.Errorf(StopIteration)
			}
			return nil
		})

		if err != nil && err.Error() == StopIteration {
			break
		} else if err != nil {
			return nil, "", err
		}
	}

	if foundCommit == nil {
		return nil, "", fmt.Errorf("no previous tag with prefix '%s' found", tagPrefix)
	}

	return foundCommit, foundTag, nil
}

func GenerateChangelogForRepo(from, to, repoURL string, subdirectories []string) []*ArtifactHubChangelogObject {
	chglogs := []*ArtifactHubChangelogObject{}
	repo := OpenRepo(repoURL)
	var fromCommit *object.Commit

	_, toCommit := ResolveTag(repo, to)

	options := &git.LogOptions{
		From: toCommit.Hash,
	}

	if len(subdirectories) > 0 {
		options = &git.LogOptions{
			From:       toCommit.Hash,
			PathFilter: filterSubirs(subdirectories),
		}
	}

	if from == "" {
		tagPrefix := extractTagPrefix(to)
		c, _, err := findPreviousTaggedCommit(repo, options, tagPrefix)
		if err != nil {
			log.Fatalf("Failed resolving latest commit: %v", err)
		}
		fromCommit = c
	} else {
		_, c := ResolveTag(repo, from)
		fromCommit = c
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
