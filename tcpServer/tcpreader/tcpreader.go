package tcpreader

import (
	"RestKeyValueStore/logger"
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

type TcpReader struct {
	ioReader io.Reader
}

func New(ioReader io.Reader) TcpReader {
	return TcpReader{ioReader: ioReader}
}

func (r *TcpReader) ReadBytes(length int) (string, error) {
	bytes := make([]byte, length)

	for i := 0; i < length; i++ {
		buf := make([]byte, 1)
		_, err := r.ioReader.Read(buf)
		if err != nil {
			logger.Error(fmt.Sprintf("Read failure: %s", err.Error()))
			return string(bytes), err
		}

		bytes[i] = buf[0]
	}
	logger.Info(fmt.Sprintf("Received bytes '%s'", string(bytes)))
	return string(bytes), nil
}

func (r *TcpReader) ParseVerb() (string, error) {
	return r.ReadBytes(3)
}

func (r *TcpReader) Parse3PartArgument() (string, error) {
	lenOfLenArg, err := r.ReadBytes(1)
	if err != nil {
		logger.Error(ErrReadingPartOneOfArg)
		return "", ErrReadingPartOneOfArg
	}

	lenArg, err := strconv.Atoi(lenOfLenArg)
	if err != nil {
		logger.Error(ErrParsingPartOneOfArg)
		return "", ErrParsingPartOneOfArg
	}

	argLenString, err := r.ReadBytes(lenArg)
	if err != nil {
		logger.Error(ErrReadingPartTwoOfArg)
		return "", ErrReadingPartTwoOfArg
	}

	argLen, err := strconv.Atoi(argLenString)
	if err != nil {
		logger.Error(ErrParsingPartTwoOfArg)
		return "", ErrParsingPartTwoOfArg
	}

	arg, err := r.ReadBytes(argLen)
	if err != nil {
		logger.Error(ErrReadingPartThreeOfArg)
		return "", ErrReadingPartThreeOfArg
	}

	return arg, nil
}

func (r TcpReader) ParseResponseLengthArg() (int, error) {
	lenOfLenArgString, err := r.ReadBytes(1)
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

	argLenString, err := r.ReadBytes(lenOfLenArg)
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
