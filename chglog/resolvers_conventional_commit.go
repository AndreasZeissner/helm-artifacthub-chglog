package chglog

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// see e.g.: @commitlint/config-conventional
// https://www.conventionalcommits.org/en/v1.0.0/#specification
type ConventionalCommitsResolver struct {
	*object.Commit
}

func NewConventionalCommitsResolver(c *object.Commit) *ConventionalCommitsResolver {
	return &ConventionalCommitsResolver{c}
}

func (r *ConventionalCommitsResolver) isAdded() bool {
	return strings.HasPrefix(r.Message, "feat") ||
		strings.HasPrefix(r.Message, "feature")
	// Simply merge is way to unspecific
	// strings.HasPrefix(r.Message, "Merge")
}

func (r *ConventionalCommitsResolver) isChanged() bool {
	return strings.HasPrefix(r.Message, "changed") ||
		strings.HasPrefix(r.Message, "deps") ||
		strings.HasPrefix(r.Message, "perf") ||
		strings.HasPrefix(r.Message, "refactor") ||
		strings.Contains(r.Message, "BREAKING CHANGE")
}

func (r *ConventionalCommitsResolver) isDeprecated() bool {
	return strings.HasPrefix(r.Message, "deprecated")
}

func (r *ConventionalCommitsResolver) isRemoved() bool {
	return strings.HasPrefix(r.Message, "removed")
}

func (r *ConventionalCommitsResolver) isFixed() bool {
	return strings.HasPrefix(r.Message, "fix")
}

func (r *ConventionalCommitsResolver) isSecurity() bool {
	return strings.Contains(r.Message, "security")
}

func (r *ConventionalCommitsResolver) parseDescription() string {
	message := r.Message

	message = strings.ReplaceAll(message, "\n", " ")
	message = strings.TrimSpace(message)

	return message
}

func (r *ConventionalCommitsResolver) ResolveChangelogEntry() (*ArtifactHubChangelogObject, error) {
	switch {
	case r.isAdded():
		return &ArtifactHubChangelogObject{
			Kind:        KindAdded,
			Description: r.parseDescription(),
			Links:       []ArtifactHubChangelogObjectLinks{},
		}, nil
	case r.isChanged():
		return &ArtifactHubChangelogObject{
			Kind:        KindChanged,
			Description: r.parseDescription(),
			Links:       []ArtifactHubChangelogObjectLinks{},
		}, nil
	case r.isDeprecated():
		return &ArtifactHubChangelogObject{
			Kind:        KindDeprecated,
			Description: r.parseDescription(),
			Links:       []ArtifactHubChangelogObjectLinks{},
		}, nil
	case r.isRemoved():
		return &ArtifactHubChangelogObject{
			Kind:        KindRemoved,
			Description: r.parseDescription(),
			Links:       []ArtifactHubChangelogObjectLinks{},
		}, nil
	case r.isFixed():
		return &ArtifactHubChangelogObject{
			Kind:        KindFixed,
			Description: r.parseDescription(),
			Links:       []ArtifactHubChangelogObjectLinks{},
		}, nil
	case r.isSecurity():
		return &ArtifactHubChangelogObject{
			Kind:        KindSecurity,
			Description: r.parseDescription(),
			Links:       []ArtifactHubChangelogObjectLinks{},
		}, nil
	default:
		return nil, fmt.Errorf("unresolvable commit: %s", r.Message)
	}
}
