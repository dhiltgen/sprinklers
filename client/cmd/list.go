package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"text/template"

	"github.com/urfave/cli"

	"github.com/dhiltgen/sprinklers/api/sprinklers"
)

const listTmpl = "{{.Name}}\t{{.Description}}\t{{.State}}\t{{.TimeRemaining}}\n"

func GetListCommand() cli.Command {
	t := template.Must(template.New("circuit").Parse(listTmpl))
	return cli.Command{
		Name:  "list",
		Usage: "list the available sprinkler circuits",
		Flags: []cli.Flag{
			cli.BoolTFlag{
				Name:  "all",
				Usage: "also show disabled circuits",
			},
		},
		Action: func(c *cli.Context) error {
			client, err := getClient(c.GlobalString("server"))
			if err != nil {
				log.Fatalf("failed to connect: %s", err)
			}

			resp, err := client.ListCircuits(context.Background(), &sprinklers.ListCircuitsRequest{})
			if err != nil {
				log.Fatalf("failed to get circuits: %s", err)
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
			fmt.Fprintln(w, "NAME\tDESCRIPTION\tWATERING NOW\tTIME REMAINING")
			for _, circuit := range resp.Items {
				if circuit.Disabled && !c.Bool("all") {
					continue
				}
				err := t.Execute(w, circuit)
				if err != nil {
					log.Fatalf("failed to render circuit: %#v: %s", circuit, err)
				}
			}
			w.Flush()
			return nil
		},
	}
}
