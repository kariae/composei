package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"github.com/kariae/composei/commands"
	"github.com/kariae/composei/version"
)

func main() {
	// Launch CLI
	app := cli.NewApp()

	app.Name = "Composei"
	app.Version = version.Version.String()
	app.Usage = "Composei is an interactive command line tool build with golang that helps you create your docker compose file."

	app.Commands = []cli.Command{
		commands.GenerateCommand,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
