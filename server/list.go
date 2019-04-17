package server

import (
	//"io"
	//"log"

	"github.com/emicklei/go-restful"
)

func ListCircuits(request *restful.Request, response *restful.Response) {
	response.WriteEntity(activeCircuits)
}
