package e2e

import (
	"testing"

	"github.com/AndreasZeissner/helm-artifacthub-chglog/chglog"
)

type testCase struct {
	repoURL        string
	from           string
	to             string
	len            int
	subdirectories []string
}

func TestGeneratingSimpleArtifactChangelog(t *testing.T) {
	var testCase = []testCase{
		{
			repoURL:        "../fixtures/go-siva",
			from:           "v1.0.0",
			to:             "v1.7.0",
			len:            4,
			subdirectories: []string{},
		},
		{
			repoURL:        "../fixtures/conventional-changelog",
			from:           "v0.0.4",
			to:             "v1.1.0",
			len:            90,
			subdirectories: []string{},
		},
		{
			repoURL:        "../fixtures/cosmo",
			from:           "router@0.13.0",
			to:             "router@0.89.2",
			len:            532,
			subdirectories: []string{},
		},
		{
			repoURL:        "../fixtures/cosmo",
			from:           "router@0.13.0",
			to:             "router@0.89.2",
			len:            179,
			subdirectories: []string{"router/"},
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
