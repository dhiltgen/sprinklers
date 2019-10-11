package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/urfave/cli"

	"github.com/dhiltgen/sprinklers/api/sprinklers"
)

const listTmpl = "{{.Name}}\t{{.Description}}\t{{.State}}\t{{.TimeRemaining}}\n"

// GetListCommand returns the list CLI command
func GetListCommand() cli.Command {
	t := template.Must(template.New("circuit").Parse(listTmpl))
	return cli.Command{
		Name:    "list",
		Usage:   "list the available sprinkler circuits",
		Aliases: []string{"ls"},

		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "all",
				Usage: "also show disabled circuits",
			},
			cli.BoolFlag{
				Name:  "quiet, q",
				Usage: "only report circuit name",
			},
			cli.StringFlag{
				Name:  "filter",
				Usage: "only show circuits matching the description string",
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
			if !c.Bool("quiet") {
				fmt.Fprintln(w, "NAME\tDESCRIPTION\tWATERING NOW\tTIME REMAINING")
			}
			filter := c.String("filter")
			for _, circuit := range resp.Items {
				if circuit.Disabled && !c.Bool("all") {
					continue
				}
				if filter != "" && !strings.Contains(circuit.Description, filter) {
					continue
				}
				if !c.Bool("quiet") {
					err := t.Execute(w, circuit)
					if err != nil {
						log.Fatalf("failed to render circuit: %#v: %s", circuit, err)
					}
				} else {
					fmt.Println(circuit.Name)
				}
			}

			w.Flush()
			return nil
		},
	}
}
