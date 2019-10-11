package sprinklers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dhiltgen/sprinklers/circuits"
	duration "github.com/golang/protobuf/ptypes/duration"
)

type server struct {
	circuits []*circuits.Circuit
}

// NewSprinklerServiceServer creates a new sprinkler server
// if dummy is set to true, a test server will be created and will not
// be hooked up to GPIO
func NewSprinklerServiceServer(dummy bool) SprinklerServiceServer {
	// TODO - for testing only
	if dummy {
		log.Printf("Running in dummy data mode...")
		circuits.DummyInit()
	}

	log.Printf("Starting sprinkler gRPC server")

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
	log.Printf("ListCircuits called")
	ret := &ListCircuitsResponse{}
	// TODO - sort, filter, and pagination
	for _, c := range s.circuits {
		ret.Items = append(ret.Items, convertCircuit(c))
	}
	return ret, nil
}

func (s *server) GetCircuit(ctx context.Context, input *GetCircuitRequest) (*Circuit, error) {
	log.Printf("GetCircuit called")
	descriptionMatches := []*Circuit{}
	for _, c := range s.circuits {
		circuit := convertCircuit(c)
		if input.Name != "" && circuit.Name == input.Name {
			return circuit, nil
		}
		if input.Description != "" && strings.Contains(circuit.Description, input.Description) {
			descriptionMatches = append(descriptionMatches, circuit)
		}
	}
	if len(descriptionMatches) > 1 {
		return nil, fmt.Errorf("ambiguous description matched %d circuits", len(descriptionMatches))
	} else if len(descriptionMatches) == 1 {
		return descriptionMatches[0], nil
	}
	return nil, fmt.Errorf("unable to located circuit %s", input.Name)

}

func (s *server) UpdateCircuit(ctx context.Context, input *Circuit) (*Circuit, error) {
	log.Printf("UpdateCircuit called")
	for _, c := range s.circuits {
		circuit := convertCircuit(c)
		if circuit.Name == input.Name {
			newCircuit := reverseConvert(input)
			//log.Printf("XXX updating with: %#v\n", newCircuit)
			c.Update(newCircuit)

			// update the status to reflect current state
			circuit = convertCircuit(c)
			return circuit, nil
		}
	}
	return nil, fmt.Errorf("unable to located circuit %s", input.Name)
}

func convertCircuit(in *circuits.Circuit) *Circuit {
	var remaining time.Duration
	var err error
	if in.TimeRemaining != "" {
		remaining, err = time.ParseDuration(in.TimeRemaining)
		if err != nil {
			log.Printf("WARNING - failed to parse duration :%s", in.TimeRemaining)
		}
	}
	return &Circuit{
		Name:             fmt.Sprintf("%d", in.GPIONumber),
		Description:      in.Name,
		WaterConsumption: float64(in.WaterConsumption),
		State:            in.State,
		TimeRemaining: &duration.Duration{
			Seconds: int64(remaining.Seconds()),
		},
	}
}

// The inverse of the above
func reverseConvert(in *Circuit) *circuits.Circuit {
	var remaining time.Duration
	if in.TimeRemaining != nil {
		remaining = time.Second * time.Duration(in.TimeRemaining.Seconds)
	}
	return &circuits.Circuit{
		// TODO - cheating and skipping fields we don't need for update...
		//Name:             fmt.Sprintf("%d", in.GPIONumber),
		//Description:      in.Name,
		//WaterConsumption: float64(in.WaterConsumption),
		State:         in.State,
		TimeRemaining: remaining.String(),
	}
}
