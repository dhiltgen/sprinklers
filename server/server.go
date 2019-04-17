package server

import (
	"io"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/dhiltgen/sprinklers/circuits"
)

var activeCircuits []*circuits.Circuit

func New() *restful.WebService {
	var err error
	activeCircuits, err = circuits.LoadCircuits()
	if err != nil {
		log.Fatal(err)
	}
	service := new(restful.WebService)
	service.
		Path("/circuits").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	service.Route(service.GET("/").To(ListCircuits))
	service.Route(service.GET("/{circuit-id}").To(FindCircuit))
	http.Handle("/metrics", promhttp.Handler())

	service.Route(service.POST("").To(UpdateCircuit))

	return service
}

func Serve() {
	container := restful.NewContainer()
	svc := New()
	container.Add(svc)
	container.Handle("/metrics", promhttp.Handler())
	server := &http.Server{Addr: ":80", Handler: container}

	log.Println("Server starting")
	log.Fatal(server.ListenAndServe())
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world")
}
