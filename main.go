package main

import (
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

	tcpServer.Startup(tcpPort, kvs)
}
