package cmd

import (
	"github.com/AndreasZeissner/helm-artifacthub-chglog/chglog"
	"github.com/urfave/cli/v2"
)

func NewChglogCli() *cli.App {
	return &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "from",
				Value:    "",
				Usage:    "tag to end with will be to - 1 in not provided",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "to",
				Value:    "here",
				Usage:    "latest tag from which to look back",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "repoURL",
				Value:    "https://github.com/src-d/go-siva",
				Usage:    "repo to generate the changelog for",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "output",
				Value:    "output.yaml",
				Usage:    "file to write the changelog to",
				Required: false,
			},
			&cli.StringSliceFlag{
				Name:     "paths",
				Usage:    "list of subdirectories from which to take the commits",
				Required: false,
			},
		},
		Name:   "chglog",
		Usage:  "generate chglog for artifactoryhub",
		Action: chglog.NewCli(),
	}
}
