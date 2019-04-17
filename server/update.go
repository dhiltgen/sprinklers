package server

import (
	//"io"
	"fmt"
	//"log"
	"net/http"

	"github.com/emicklei/go-restful"

	"github.com/dhiltgen/sprinklers/circuits"
)

func UpdateCircuit(request *restful.Request, response *restful.Response) {
	circuit := new(circuits.Circuit)
	err := request.ReadEntity(&circuit)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	for _, c := range activeCircuits {
		if c.GPIONumber == circuit.GPIONumber {
			err := c.Update(circuit)
			if err != nil {
				response.WriteError(http.StatusInternalServerError, err)
				return
			}
			response.WriteEntity(c)
			return
		}
	}

	response.WriteError(http.StatusInternalServerError, fmt.Errorf("Unable to locate circuit"))
}
