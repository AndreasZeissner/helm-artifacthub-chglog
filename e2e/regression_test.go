package e2e

import (
	"testing"

	"github.com/AndreasZeissner/helm-artifacthub-chglog/chglog"
)

func TestGeneratingSimpleArtifactChangelog(t *testing.T) {
	var testCase = []struct {
		repoURL string
		res     bool
		from    string
		to      string
		len     int
	}{
		{
			repoURL: "../fixtures/go-siva",
			from:    "v1.0.0",
			to:      "v1.7.0",
			len:     3,
		},
		{
			repoURL: "../fixtures/conventional-changelog",
			from:    "v1.1.0",
			to:      "v0.0.4",
			len:     52,
		},
	}

	for _, tt := range testCase {
		t.Run(tt.repoURL, func(t *testing.T) {
			logs := chglog.GenerateChangelogForRepo(tt.from, tt.to, tt.repoURL, []string{})
			if tt.len != len(logs) {
				t.Errorf("want %d got %d", tt.len, len(logs))
			}
		})
	}
}
