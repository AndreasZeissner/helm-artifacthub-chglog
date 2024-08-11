package main

import (
	"os"

	"github.com/AndreasZeissner/helm-artifacthub-chglog/cmd"
)

func main() {
	cmd.NewChglogCli().Run(os.Args)
}
