package storehandler

import (
	"RestKeyValueStore/arguments"
	"RestKeyValueStore/logger"
	"RestKeyValueStore/store"
	"RestKeyValueStore/tcpServer/tcpreader"
	"fmt"
)

func HandlePut(key, value string, s store.Store) string {
	s.Put(store.Key(key), store.Entry(value))
	logger.Info(fmt.Sprintf("Key: %s Value: %s added to store", key, value))
	return "ack"
}

func HandleGet(tcpReader tcpreader.TcpReader, s store.Store) string {
	key, err := tcpReader.Parse3PartArgument()
	if err != nil {
		return "err"
	}

	valueLenLimit, err := tcpReader.ParseResponseLengthArg()
	if err != nil {
		return "err"
	}

	value, err := s.Get(store.Key(key))

	if err != nil {
		return "nil"
	}

	if valueLenLimit != 0 && valueLenLimit < len(value) {
		value = value[:valueLenLimit]
	}

	logger.Info(fmt.Sprintf("Key: %s reteived from store", key))
	return "val" + arguments.Encode(value)
}

func HandleDelete(tcpReader tcpreader.TcpReader, s store.Store) string {
	key, err := tcpReader.Parse3PartArgument()
	if err != nil {
		return "err"
	}

	err = s.Delete(store.Key(key))

	if err != nil {
		return "ack"
	}

	logger.Info(fmt.Sprintf("Key: %s deleted from store", key))
	return "ack"
}
