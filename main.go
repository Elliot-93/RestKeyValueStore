package main

import (
	"RestKeyValueStore/endpoints"
	"RestKeyValueStore/server"
	"RestKeyValueStore/server/middleware"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	var port int

	fmt.Println(os.Args)

	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()

	router := http.NewServeMux()

	for path, routeDetail := range endpoints.Routes {
		if routeDetail.AuthRequired {
			router.Handle(path,
				middleware.LoggingMiddleware(
					middleware.PanicHandler(
						middleware.Authenticate(routeDetail.Handler))))
		} else {
			router.Handle(path,
				middleware.LoggingMiddleware(
					middleware.PanicHandler(routeDetail.Handler)))
		}
	}

	server.Startup(port, router)
}
