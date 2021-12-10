package main

import (
	distributionServer "RestKeyValueStore/distribution/server"
	"RestKeyValueStore/store"
	"RestKeyValueStore/tcpServer"
	"flag"
	"fmt"
	"os"
)

func main() {
	var tcpPort int

	fmt.Println(os.Args)

	flag.IntVar(&tcpPort, "tcpPort", 1234, "tcpPort to listen on")
	flag.Parse()

	kvs := store.New()

	go distributionServer.Startup(kvs)

	tcpServer.Startup(tcpPort, kvs)
}
