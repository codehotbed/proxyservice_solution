package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"ProxyService/proxyservice/internal/register"
	"ProxyService/proxyservice/internal/resources"
	"ProxyService/proxyservice/openapi/autogen/proxyservice/server"
	"ProxyService/proxyservice/openapi/autogen/proxyservice/server/operations"
)

var (
	servicePort = 3000
)

func main() {
	//setup prometheous
	setupPrometheusHandler()

	//load swagger file
	swaggerDocumentation, err := loads.Analyzed(server.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	//create new service API
	openapi := operations.NewProxyServiceAPI(swaggerDocumentation)
	server := server.NewServer(openapi)
	defer server.Shutdown()

	//assign port which service will be running on
	server.Port = servicePort

	//create queue object and pass it to Router
	rabbitQueue := resources.NewQueueResource()

	// create router, and register all defined routes handlers
	router := register.NewRouter(openapi, rabbitQueue)
	router.RegisterRoutes()

	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func setupPrometheusHandler() {
	go func() {
		fmt.Println("Serving prometheus handler at http://[::]:9090/metrics")
		http.Handle("/metrics", promhttp.Handler())
		panic(http.ListenAndServe(":9090", nil))
	}()
}
