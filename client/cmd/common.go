package cmd

import (
	"google.golang.org/grpc"

	"github.com/dhiltgen/sprinklers/api/sprinklers"
)

func getClient(server string) (sprinklers.SprinklerServiceClient, error) {
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := sprinklers.NewSprinklerServiceClient(conn)
	return client, nil
}
