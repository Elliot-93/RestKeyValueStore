package reader

import (
	"RestKeyValueStore/logger"
	"fmt"
	"io"
)

func ReadBytes(r io.Reader, length int) (string, error) {
	bytes := make([]byte, length)

	for i := 0; i < length; i++ {
		buf := make([]byte, 1)
		_, err := r.Read(buf)
		if err != nil {
			logger.Error(fmt.Sprintf("Read failure: %s", err.Error()))
			return string(bytes), err
		}

		bytes[i] = buf[0]
	}
	logger.Info(fmt.Sprintf("Received bytes '%s'", string(bytes)))
	return string(bytes), nil
}
