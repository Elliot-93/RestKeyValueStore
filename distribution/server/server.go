package server

import (
	"RestKeyValueStore/distribution/nodes"
	"RestKeyValueStore/logger"
	"RestKeyValueStore/store"
	"RestKeyValueStore/tcpServer/handler/storehandler"
	"RestKeyValueStore/tcpServer/tcpreader"
	"fmt"
	"net"
	"strings"
)

func Startup(kvs store.KeyValueStore) {

	var listener net.Listener

	for port, _ := range nodes.AvailablePorts {
		var err error
		listener, err = net.Listen("tcp", ":"+port)
		if err != nil {
			logger.Info(fmt.Sprintf("Distribution server listen failure: %s", err.Error()))
			continue
		} else {
			logger.Info(fmt.Sprintf("Distribution server listening on port: %s", port))
			nodes.SetListeningPort(port)
			break
		}
	}

	if listener == nil {
		logger.Fatal(fmt.Sprint("No available distribution ports, exiting"))
		panic("No available distribution ports")
	}

	defer func() { _ = listener.Close() }()

	go func() {
		for port, _ := range nodes.AvailablePorts {
			tryConnectToNode(port)
		}
	}()

	for {
		logger.Info("Waiting for distribution node connection ...")

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

func tryConnectToNode(port string) {
	c, err := net.Dial("tcp", "localhost:"+port)

	if err != nil {
		logger.Info(fmt.Sprintf("port %s not active err: %v", port, err))
		return
	}

	if _, conErr := c.Write([]byte("con" + nodes.ListingPort)); err != nil {
		logger.Error(conErr)
		return
	}

	resp := make([]byte, 3)
	_, readErr := c.Read(resp)

	if readErr != nil || string(resp) != "ack" {
		logger.Error(fmt.Sprintf("response should be ack but got: %s", resp))
		return
	}

	//todo: may need to keep alive
	//todo: if connection closed remove from active connections list and retry
	nodes.ActiveNodeConnections[port] = c
	logger.Info(fmt.Sprintf("connected to node at port: %s", port))
}

func handleConnection(conn net.Conn, kvs store.KeyValueStore) {
	defer func() { _ = conn.Close() }()

	tcpReader := tcpreader.New(conn)

	for {
		verb, err := tcpReader.ParseVerb()
		if err != nil {
			return
		}

		var response string

		switch strings.ToLower(verb) {
		case "con":
			remotePort, readPortErr := tcpReader.ReadBytes(4)
			if readPortErr != nil {
				logger.Error(fmt.Sprintf("Error parsing port err: %v", readPortErr))
				return
			}

			if _, portPresent := nodes.AvailablePorts[remotePort]; !portPresent {
				logger.Error(fmt.Sprintf("Distribution port %s not valid", remotePort))
				return
			}

			conn.Write([]byte("ack"))

			if _, portPresent := nodes.ActiveNodeConnections[remotePort]; portPresent {
				logger.Info(fmt.Sprintf("Already connected to node at port: %s", remotePort))
				return
			}

			tryConnectToNode(remotePort)

		case "upd":
			logger.Info("upd received")

			key, value, parseErr := tcpReader.ParseKeyValueArgs()
			if parseErr != nil {
				response = "err"
			}

			response = storehandler.HandlePut(key, value, kvs)

		default:
			logger.Fatal(fmt.Sprintf("Verb %s not supported", verb))
			return
		}

		logger.Info(fmt.Sprintf("Writing response '%s'", response))
		if _, writeErr := conn.Write([]byte(response)); err != nil {
			logger.Fatal(fmt.Sprintf("Write response failure: %s", writeErr.Error()))
		}
	}
}
