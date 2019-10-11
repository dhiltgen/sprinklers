package circuits

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// CircuitDefinition defines the base settings for a circuit
type CircuitDefinition struct {
	// GPIONumber defines the GPIO Pin as defined by the CPU (not header)
	GPIONumber uint8 `json:"gpio"`

	// Name is the human readable name for the circuit controlled by this relay
	Name string `json:"name"`

	// WaterConsumption is the approximate amount of water consumed by the circuit per hour
	WaterConsumption float64 `json:"consumption"`

	// Disabled if set to true indicates this circuit is currently unusable
	Disabled bool `json:"disabled"`
}

// Circuit expands a CircuitDefinition to include current state
type Circuit struct {
	CircuitDefinition
	State         bool               `json:"state"`
	TimeRemaining string             `json:"remaining"`
	pin           Pin                `json:-`
	cancel        context.CancelFunc `json:-`

	// Metrics for how much water we've accumulated
	accumulation float64   `json:-`
	started      time.Time `json:-`
}

// CircuitFile defines where we load the circuit instance data from (json format)
var CircuitFile = "circuits.json"

// AccumulationUpdateInterval defines how frequently we update the accumulation of delivered water
var AccumulationUpdateInterval = 10 * time.Second

var accumulation = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "zone_water_inches",
		Help: "Watering zone accumulated inches of water",
	},
	[]string{"circuit", "name"},
)

func init() {
	prometheus.MustRegister(accumulation)
}

// LoadCircuits will load up the known circuits from the definition file
func LoadCircuits() ([]*Circuit, error) {

	data, err := ioutil.ReadFile(CircuitFile)
	if err != nil {
		return nil, err
	}

	// TODO load these as definitions then convert to full circuits
	var circuits []*Circuit
	err = json.Unmarshal(data, &circuits)
	if err != nil {
		return nil, fmt.Errorf("Failed to process %s: %s", CircuitFile, err)
	}

	err = pinsInit()
	// close with rpio.Close()
	if err != nil {
		return nil, fmt.Errorf("Failed to open GPIO device: %s", err)
	}

	// Now load all the circuits
	for i := range circuits {
		circuits[i].init()
	}

	return circuits, nil
}
