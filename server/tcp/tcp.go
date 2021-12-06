package tcp

import (
	"RestKeyValueStore/logger"
	"RestKeyValueStore/store"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func Startup(port int) {
	logger.Info(fmt.Sprintf("TCP Server started up on port: %d", port))
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	if err != nil {
		logger.Fatal(fmt.Sprintf("connection failure: %s", err.Error()))
	}

	defer func() { _ = listener.Close() }()

	for {
		logger.Info("Waiting for client connection ...")
		if c, err := listener.Accept(); err != nil {
			logger.Info(fmt.Sprintf("Connection error \nAddr: %s \nError: %s", c.RemoteAddr(), err.Error()))
		} else {
			logger.Info(fmt.Sprintf("Connection established: %s", c.RemoteAddr()))
			go handleConnection(c)
		}
	}
}

func handleConnection(c net.Conn) {
	defer func() {
		_ = c.Close()
		logger.Info(fmt.Sprintf("Connection closed: %s", c.RemoteAddr()))
	}()

	verb, err := readBytes(c, 3)
	if err != nil {
		return
	}

	var response string

	switch verb {
	case "put":
		err = handlePut(c)
		if err != nil {
			response = err.Error()
		}
		response = "ack"
	default:
		logger.Fatal(fmt.Sprintf("Verb %s not supported", verb))
	}

	logger.Info(fmt.Sprintf("Writing response '%s'", response))
	if _, writeErr := c.Write([]byte(response)); err != nil {
		logger.Fatal(fmt.Sprintf("Write response failure: %s", writeErr.Error()))
	}
}

func readBytes(c net.Conn, length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := c.Read(bytes); err != nil {
		logger.Fatal(fmt.Sprintf("Read failure: %s", err.Error()))
		return "", err
	} else {
		logger.Info(fmt.Sprintf("Received message '%s'", string(bytes)))
		return strings.ToLower(string(bytes)), nil
	}
}

func handlePut(c net.Conn) error {
	key, err := parseArgument(c)
	if err != nil {
		return err
	}

	value, err := parseArgument(c)
	if err != nil {
		return err
	}

	store.Put(store.Key(key), store.Entry{Value: value, Owner: "TCP"}, false)
	logger.Info(fmt.Sprintf("Key: %s Value: %s added to store", key, value))
	return nil
}

func parseArgument(c net.Conn) (string, error) {
	lenOfLenArg, err := readBytes(c, 1)
	if err != nil {
		return "", errors.New("error reading part 1 length of length argument")
	}

	lenArg, err := strconv.Atoi(lenOfLenArg)
	if err != nil {
		return "", errors.New("error parsing length of length argument")
	}

	argLenString, err := readBytes(c, lenArg)
	if err != nil {
		return "", errors.New("error reading part 2 length argument")
	}

	argLen, err := strconv.Atoi(argLenString)
	if err != nil {
		return "", errors.New("error parsing length argument")
	}

	arg, err := readBytes(c, argLen)
	if err != nil {
		return "", errors.New("error reading part 3 argument")
	}

	return arg, nil
}
