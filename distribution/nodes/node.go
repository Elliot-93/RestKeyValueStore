package nodes

import (
	"RestKeyValueStore/logger"
	"fmt"
	"net"
)

var (
	AvailablePorts = map[string]struct{}{
		"3001": {},
		"3002": {},
		//"3003": {},
		//"3004": {},
		//"3005": {},
	}
	ActiveNodeConnections map[string]net.Conn
	ListingPort           string
)

func init() {
	ActiveNodeConnections = make(map[string]net.Conn)
}

func SetListeningPort(port string) {
	delete(AvailablePorts, port)
	ListingPort = port
}

func DistributeUpdate(command string) {
	command = "upd" + command
	for port, node := range ActiveNodeConnections {
		logger.Info(fmt.Sprintf("upd %s pushed to port:%s", command, port))
		node.Write([]byte(command))
	}
}
