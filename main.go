package main

import (
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/codegangsta/cli"
	"github.com/untoldwind/gorrd/commands"
	"github.com/untoldwind/gorrd/config"
)

func main() {
	app := cli.NewApp()
	app.Name = "gorrd"
	app.Usage = "Go implementation of the rrdtool"
	app.Version = config.Version()
	app.Commands = []cli.Command{
		commands.CreateCommand,
		commands.UpdateCommand,
		commands.DumpCommand,
		commands.LastUpdateCommand,
	}

	if err := app.Run(os.Args); err != nil {
		log.Errorf("Failed to run command: %s", err.Error())
	}
}
