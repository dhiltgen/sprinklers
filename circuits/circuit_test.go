package circuits

import (
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func WriteSomeCircuitDefs(t *testing.T) {
	CircuitFile = "/tmp/circuits.json"
	err := ioutil.WriteFile(CircuitFile, []byte(`
[
    {
        "gpio": 2,
        "name": "Test Circuit",
        "consumption": 36000.0
    }
]
`), 0644)
	require.NoError(t, err)
}

func TestNoFileLoadCircuits(t *testing.T) {
	dummyInit()
	log.Println("Loading circuits")
	CircuitFile = "garbage"
	circuits, err := LoadCircuits()
	require.Error(t, err)
	require.Contains(t, err.Error(), "no such file or directory")
	require.Nil(t, circuits)
}

func TestIsDirLoadCircuits(t *testing.T) {
	dummyInit()
	log.Println("Loading circuits")
	CircuitFile = "/tmp/"
	circuits, err := LoadCircuits()
	require.Error(t, err)
	require.Contains(t, err.Error(), "is a directory")
	require.Nil(t, circuits)
}

func TestNotJSONLoadCircuits(t *testing.T) {
	dummyInit()
	log.Println("Loading circuits")
	CircuitFile = "/tmp/circuits.json"
	err := ioutil.WriteFile(CircuitFile, []byte("not json data"), 0644)
	require.NoError(t, err)
	circuits, err := LoadCircuits()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid character")
	require.Nil(t, circuits)
}

func TestEmptyJSONLoadCircuits(t *testing.T) {
	dummyInit()
	log.Println("Loading circuits")
	CircuitFile = "/tmp/circuits.json"
	err := ioutil.WriteFile(CircuitFile, []byte("[]"), 0644)
	require.NoError(t, err)
	circuits, err := LoadCircuits()
	require.NoError(t, err)
	require.NotNil(t, circuits)
	require.Equal(t, 0, len(circuits))
}

func TestAccumulation(t *testing.T) {
	dummyInit()
	log.Println("Loading circuits")
	WriteSomeCircuitDefs(t)
	AccumulationUpdateInterval = 10 * time.Millisecond
	circuits, err := LoadCircuits()
	require.NoError(t, err)
	require.NotNil(t, circuits)
	require.Equal(t, 1, len(circuits))

	c := circuits[0]
	require.Equal(t, 0.0, c.accumulation)
	updater := Circuit{
		TimeRemaining: "1s",
		State:         true,
	}
	c.Update(&updater)
	waitForIt := time.NewTimer(1 * time.Second)
	<-waitForIt.C
	log.Printf("accumulation was %f\n", c.accumulation)
	require.True(t, c.accumulation > 9.0)
	require.True(t, c.accumulation < 11.0)
}
