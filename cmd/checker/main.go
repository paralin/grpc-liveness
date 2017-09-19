package main

import (
	"github.com/paralin/grpc-liveness/checker"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "checker"
	app.Usage = "Checks to see if a GRPC service is online or not, by calling a status service."
	app.HideVersion = true
	app.Commands = checker.CheckerSubCommands
	app.Flags = checker.CheckerFlags
	app.RunAndExitOnError()
}
