package main

import (
	"RestKeyValueStore/config"
	"RestKeyValueStore/endpoints"
	"RestKeyValueStore/server/http"
	"RestKeyValueStore/server/middleware"
	"RestKeyValueStore/server/tcp"
	"flag"
	"fmt"
	netHttp "net/http"
	"os"
)

func main() {
	var httpPort int
	var tcpPort int
	var depth int

	fmt.Println(os.Args)

	flag.IntVar(&httpPort, "httpPort", 8080, "httpPort to listen on")
	flag.IntVar(&tcpPort, "tcpPort", 8090, "tcpPort to listen on")
	flag.IntVar(&depth, "depth", 100, "httpPort to listen on")
	flag.Parse()

	config.Populate(httpPort, depth)

	router := netHttp.NewServeMux()

	for path, routeDetail := range endpoints.Routes {
		if routeDetail.AuthRequired {
			router.Handle(path,
				middleware.PanicHandler(
					middleware.LoggingMiddleware(
						middleware.Authenticate(routeDetail.Handler))))
		} else {
			router.Handle(path,
				middleware.PanicHandler(
					middleware.LoggingMiddleware(routeDetail.Handler)))
		}
	}

	go tcp.Startup(tcpPort)

	http.Startup(router)
}
