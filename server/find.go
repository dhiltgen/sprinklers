package server

import (
	"fmt"
	"strconv"
	//"io"
	//"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

func FindCircuit(request *restful.Request, response *restful.Response) {
	idStr := request.PathParameter("circuit-id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WriteError(http.StatusBadRequest, fmt.Errorf("ID not valid: %s", err))
	}

	if id < 0 || id >= len(activeCircuits) {
		response.WriteError(http.StatusBadRequest, fmt.Errorf("ID out of range"))
	}
	response.WriteEntity(activeCircuits[id])
}
