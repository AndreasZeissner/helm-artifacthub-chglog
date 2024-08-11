package e2e

import (
	"testing"

	"github.com/AndreasZeissner/helm-artifacthub-chglog/chglog"
)

func TestGeneratingSimpleArtifactChangelog(t *testing.T) {
	var testCase = []struct {
		repoURL        string
		res            bool
		from           string
		to             string
		len            int
		subdirectories []string
	}{
		{
			repoURL:        "../fixtures/go-siva",
			from:           "v1.0.0",
			to:             "v1.7.0",
			len:            3,
			subdirectories: []string{},
		},
		{
			repoURL:        "../fixtures/conventional-changelog",
			from:           "v1.1.0",
			to:             "v0.0.4",
			len:            52,
			subdirectories: []string{},
		},
		{
			repoURL:        "../fixtures/cosmo",
			from:           "router@0.89.2",
			to:             "router@0.13.0",
			len:            214,
			subdirectories: []string{},
		},
		{
			repoURL:        "../fixtures/cosmo",
			from:           "router@0.89.2",
			to:             "router@0.13.0",
			len:            104,
			subdirectories: []string{"router"},
		},
	}

	for _, tt := range testCase {
		t.Run(tt.repoURL, func(t *testing.T) {
			logs := chglog.GenerateChangelogForRepo(tt.from, tt.to, tt.repoURL, tt.subdirectories)
			if tt.len != len(logs) {
				t.Errorf("want %d got %d", tt.len, len(logs))
			}
		})
	}
}
