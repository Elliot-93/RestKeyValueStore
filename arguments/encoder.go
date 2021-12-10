package arguments

import (
	"fmt"
	"strconv"
)

func Encode(response string) string {
	var responseLength = len(response)
	var responseLengthLength = len(strconv.Itoa(responseLength))
	return fmt.Sprintf("%d%d%s", responseLengthLength, responseLength, response)
}
