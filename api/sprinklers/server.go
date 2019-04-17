package sprinklers

import (
	"context"
	"log"

	"github.com/dhiltgen/sprinklers/circuits"
)

type server struct {
	circuits []*circuits.Circuit
}

func NewSprinklerServiceServer() SprinklerServiceServer {
	activeCircuits, err := circuits.LoadCircuits()
	if err != nil {
		log.Fatal(err)
	}

	s := &server{
		circuits: activeCircuits,
	}

	return s
}

func (s *server) ListCircuits(ctx context.Context, _ *ListCircuitsRequest) (*ListCircuitsResponse, error) {
	ret := &ListCircuitsResponse{}
	for _, c := range s.circuits {
		ret.Items = append(ret.Items, convertCircuit(c))
	}
	return ret, nil
}

func convertCircuit(in *circuits.Circuit) *Circuit {
	return &Circuit{
		Name:             in.Name,
		WaterConsumption: float64(in.WaterConsumption),
		State:            in.State,
		//TimeRemaining:    in.TimeRemaining,
	}
}
