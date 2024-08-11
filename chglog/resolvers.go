package chglog

import "github.com/go-git/go-git/v5/plumbing/object"

type ArtifactHubChangelogKind string

const (
	// added, changed, deprecated, removed, fixed and security.
	// https://artifacthub.io/docs/topics/annotations/helm/
	KindAdded      ArtifactHubChangelogKind = "added"
	KindChanged    ArtifactHubChangelogKind = "changed"
	KindDeprecated ArtifactHubChangelogKind = "deprecated"
	KindRemoved    ArtifactHubChangelogKind = "removed"
	KindFixed      ArtifactHubChangelogKind = "fixed"
	KindSecurity   ArtifactHubChangelogKind = "security"
)

type ArtifactHubChangelogObjectLinks struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type ArtifactHubChangelogObject struct {
	Kind        ArtifactHubChangelogKind          `yaml:"kind"`
	Description string                            `yaml:"description"`
	Links       []ArtifactHubChangelogObjectLinks `yaml:"links,omitempty"`
}

type Resolver string

const (
	ConventionalCommits Resolver = "conventional-commits"
)

type ArtifactHubChangelogResolverInterface interface {
	ResolveChangelogEntry() *ArtifactHubChangelogObject

	isAdded() bool
	isChanged() bool
	isDeprecated() bool
	isFixed() bool
	isRemoved() bool
	isSecurity() bool

	parseDescription() string
}

func NewCommitResolver(c *object.Commit, r Resolver) ArtifactHubChangelogResolverInterface {
	switch r {
	default:
		return NewConventionalCommitsResolver(c)
	}
}
