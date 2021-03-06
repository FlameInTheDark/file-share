package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

var (
	// build is the git version of this program. It is set using build flags in the makefile.
	build = "develop"
)

func main() {
	app := &cli.App{
		Name:    "api",
		Usage:   "run api service",
		Version: build,
		Commands: []*cli.Command{
			ApiCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
