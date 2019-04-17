package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/dhiltgen/sprinklers/api/sprinklers"
)

// TODO - replace with a "real" client

func main() {
	conn, err := grpc.Dial("localhost:1600", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %s", err)
	}
	client := sprinklers.NewSprinklerServiceClient(conn)

	circuits, err := client.ListCircuits(context.Background(), &sprinklers.ListCircuitsRequest{})
	if err != nil {
		log.Fatalf("failed to get circuits: %s", err)
	}
	fmt.Printf("Circuits: %#v", circuits)
}
