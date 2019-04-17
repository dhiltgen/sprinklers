package main

// TODO - this doesn't belong here...

import (
	"log"
	"net"

	"github.com/dhiltgen/sprinklers/api/sprinklers"

	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":1600")
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	grpcServer := grpc.NewServer()
	server := sprinklers.NewSprinklerServiceServer()

	sprinklers.RegisterSprinklerServiceServer(grpcServer, server)
	grpcServer.Serve(lis)
}
