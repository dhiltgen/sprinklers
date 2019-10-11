package sprinklers

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/dhiltgen/sprinklers/circuits"

	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	WriteSomeCircuitDefs(t)

	srv := NewSprinklerServiceServer(true)
	require.NotNil(t, srv)
	resp, err := srv.ListCircuits(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, len(resp.Items), 1)
}

func WriteSomeCircuitDefs(t *testing.T) {
	circuits.CircuitFile = "/tmp/circuits.json"
	err := ioutil.WriteFile(circuits.CircuitFile, []byte(`
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
