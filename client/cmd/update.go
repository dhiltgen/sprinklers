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
	duration "github.com/golang/protobuf/ptypes/duration"
)

const updateTmpl = "{{.Name}}\t{{.Description}}\t{{.State}}\t{{.TimeRemaining}}\n"

func GetUpdateCommand() cli.Command {
	t := template.Must(template.New("circuit").Parse(listTmpl))
	return cli.Command{
		Name:  "update",
		Usage: "list the available sprinkler circuits",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "stop",
				Usage: "immediately stop watering",
			},
			cli.BoolFlag{
				Name:  "start",
				Usage: "immediately start watering",
			},
			cli.DurationFlag{
				Name:  "stop-after",
				Usage: "start watering, and stop after specified duration",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return fmt.Errorf("you must specify a circuit by name/number")
			}
			name := c.Args().Get(0)
			client, err := getClient(c.GlobalString("server"))
			if err != nil {
				return fmt.Errorf("failed to connect to sprinkler server: %s", err)
			}

			// Validate flags are sane
			if c.Bool("start") && c.Bool("stop") {
				return fmt.Errorf("you must only specify one start or stop action")
			}

			circuit, err := client.GetCircuit(context.Background(), &sprinklers.GetCircuitRequest{
				Name: name,
			})
			if err != nil {
				log.Fatalf("failed to get circuit: %s", err)
			}
			if c.Bool("start") {
				circuit.State = true
				if d := c.Duration("stop-after"); d != 0 {
					circuit.TimeRemaining = &duration.Duration{
						Seconds: int64(d.Seconds()),
					}
				}
			} else if c.Bool("stop") {
				circuit.State = false
			}
			circuit, err = client.UpdateCircuit(context.Background(), circuit)
			if err != nil {
				log.Fatalf("failed to update circuit: %s", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
			fmt.Fprintln(w, "NAME\tDESCRIPTION\tWATERING NOW\tTIME REMAINING")
			err = t.Execute(w, circuit)
			w.Flush()
			return nil
		},
	}
}
