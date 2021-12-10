package tcpServer

import (
	"RestKeyValueStore/arguments"
	"RestKeyValueStore/distribution/nodes"
	"RestKeyValueStore/logger"
	"RestKeyValueStore/store"
	"RestKeyValueStore/tcpServer/handler/storehandler"
	"RestKeyValueStore/tcpServer/tcpreader"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func Startup(port int, kvs store.KeyValueStore) {
	logger.Info(fmt.Sprintf("TCP Server started up on port: %d", port))
	listener, errListener := net.Listen("tcp", ":"+strconv.Itoa(port))

	if errListener != nil {
		logger.Fatal(fmt.Sprintf("connection failure: %s", errListener.Error()))
	}

	defer func() { _ = listener.Close() }()

	for {
		logger.Info("Waiting for client connection ...")

		c, err := listener.Accept()

		if err != nil {
			logger.Info(fmt.Sprintf("Connection error \nAddr: %s \nError: %s", c.RemoteAddr(), err.Error()))
			continue
		}

		logger.Info(fmt.Sprintf("Connection established: %s", c.RemoteAddr()))
		go func() {
			handleConnection(c, kvs)
			logger.Info(fmt.Sprintf("Connection closed: %s", c.RemoteAddr()))
		}()
	}
}

func handleConnection(rwc io.ReadWriteCloser, kvs store.KeyValueStore) {
	defer func() { _ = rwc.Close() }()

	tcpReader := tcpreader.New(rwc)

	for {
		verb, err := tcpReader.ParseVerb()
		if err != nil {
			return
		}

		var response string

		switch strings.ToLower(verb) {
		case "put":
			key, value, parseErr := tcpReader.ParseKeyValueArgs()
			if parseErr != nil {
				response = "err"
			}

			response = storehandler.HandlePut(key, value, kvs)

			nodes.DistributeUpdate(verb + arguments.Encode(key) + arguments.Encode(value))

		case "get":
			response = storehandler.HandleGet(tcpReader, kvs)
		case "del":
			response = storehandler.HandleDelete(tcpReader, kvs)
		case "bye":
			return
		default:
			logger.Fatal(fmt.Sprintf("Verb %s not supported", verb))
			return
		}

		logger.Info(fmt.Sprintf("Writing response '%s'", response))
		if _, writeErr := rwc.Write([]byte(response)); err != nil {
			logger.Fatal(fmt.Sprintf("Write response failure: %s", writeErr.Error()))
		}
	}
}
