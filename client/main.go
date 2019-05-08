package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/dhiltgen/sprinklers/client/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "Sprinklers"
	app.Usage = "manage sprinkler circuits"
	app.Commands = []cli.Command{
		cmd.GetListCommand(),
		cmd.GetUpdateCommand(),
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "server",
			Value: "sprinklers:1600",
			Usage: "specify the sprinkler server to use",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
