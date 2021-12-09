package storehandler

import (
	"RestKeyValueStore/logger"
	"RestKeyValueStore/store"
	"RestKeyValueStore/tcpServer/reader"
	"errors"
	"fmt"
	"io"
	"strconv"
)

var (
	ErrReadingPartOneOfArg   = errors.New("error reading part 1 of argument, length of part 2")
	ErrParsingPartOneOfArg   = errors.New("error parsing part 1 of argument, must be int value")
	ErrReadingPartTwoOfArg   = errors.New("error reading part 2 of argument, length of part 3")
	ErrParsingPartTwoOfArg   = errors.New("error parsing part 2 of argument, must be int value")
	ErrReadingPartThreeOfArg = errors.New("error reading part 3 of argument, the value itself")
)

func HandlePut(r io.Reader, s store.Store) string {
	key, err := parseArgument(r)
	if err != nil {
		return "err"
	}

	value, err := parseArgument(r)
	if err != nil {
		return "err"
	}

	s.Put(store.Key(key), store.Entry(value))
	logger.Info(fmt.Sprintf("Key: %s Value: %s added to store", key, value))
	return "ack"
}

func HandleGet(r io.Reader, s store.Store) string {
	key, err := parseArgument(r)
	if err != nil {
		return "err"
	}

	valueLenLimit, err := parseGetLengthArgument(r)
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
	return "val" + encodeResponse(value)
}

func HandleDelete(r io.Reader, s store.Store) string {
	key, err := parseArgument(r)
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

func parseArgument(r io.Reader) (string, error) {
	lenOfLenArg, err := reader.ReadBytes(r, 1)
	if err != nil {
		logger.Error(ErrReadingPartOneOfArg)
		return "", ErrReadingPartOneOfArg
	}

	lenArg, err := strconv.Atoi(lenOfLenArg)
	if err != nil {
		logger.Error(ErrParsingPartOneOfArg)
		return "", ErrParsingPartOneOfArg
	}

	argLenString, err := reader.ReadBytes(r, lenArg)
	if err != nil {
		logger.Error(ErrReadingPartTwoOfArg)
		return "", ErrReadingPartTwoOfArg
	}

	argLen, err := strconv.Atoi(argLenString)
	if err != nil {
		logger.Error(ErrParsingPartTwoOfArg)
		return "", ErrParsingPartTwoOfArg
	}

	arg, err := reader.ReadBytes(r, argLen)
	if err != nil {
		logger.Error(ErrReadingPartThreeOfArg)
		return "", ErrReadingPartThreeOfArg
	}

	return arg, nil
}

func parseGetLengthArgument(r io.Reader) (int, error) {
	lenOfLenArgString, err := reader.ReadBytes(r, 1)
	if err != nil {
		logger.Error(ErrReadingPartOneOfArg)
		return 0, ErrReadingPartOneOfArg
	}

	lenOfLenArg, err := strconv.Atoi(lenOfLenArgString)
	if err != nil {
		logger.Error(ErrParsingPartOneOfArg)
		return 0, ErrParsingPartOneOfArg
	}

	if lenOfLenArg == 0 {
		return 0, nil
	}

	argLenString, err := reader.ReadBytes(r, lenOfLenArg)
	if err != nil {
		logger.Error(ErrReadingPartTwoOfArg)
		return 0, ErrReadingPartTwoOfArg
	}

	argLen, err := strconv.Atoi(argLenString)
	if err != nil {
		logger.Error(ErrParsingPartTwoOfArg)
		return 0, ErrParsingPartTwoOfArg
	}

	return argLen, nil
}

func encodeResponse(response string) string {
	var responseLength = len(response)
	var responseLengthLength = len(strconv.Itoa(responseLength))
	return fmt.Sprintf("%d%d%s", responseLengthLength, responseLength, response)
}
