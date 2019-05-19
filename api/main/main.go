package main

// TODO - this doesn't belong here...

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/dhiltgen/sprinklers/api/sprinklers"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func main() {
	app := cli.NewApp()
	app.Name = "Sprinkler Server"
	app.Usage = "manage sprinkler circuits"
	app.Action = func(c *cli.Context) error {
		dummy := c.GlobalBool("dummy")
		port := c.GlobalString("port")
		metricsPort := c.GlobalString("metrics-port")
		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			log.Fatalf("failed to listen: %s", err)
		}
		grpcServer := grpc.NewServer(
			grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		)
		server := sprinklers.NewSprinklerServiceServer(dummy)

		sprinklers.RegisterSprinklerServiceServer(grpcServer, server)
		grpc_prometheus.Register(grpcServer)

		http.Handle("/metrics", promhttp.Handler())
		go func() {
			log.Println("Server metrics server")
			log.Fatal(http.ListenAndServe(":"+metricsPort, nil))
		}()

		return grpcServer.Serve(lis)
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "dummy",
			Usage: "use dummy data instead of real circuits",
		},
		cli.StringFlag{
			Name:  "port",
			Value: "1600",
			Usage: "the gRPC port number",
		},
		cli.StringFlag{
			Name:  "metrics-port",
			Value: "1601",
			Usage: "the metrics port number",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
