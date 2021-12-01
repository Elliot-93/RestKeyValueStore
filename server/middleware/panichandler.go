package middleware

import (
	"RestKeyValueStore/logger"
	"fmt"
	"net/http"
)

// PanicHandler logs any panics that occur during handling a request
func PanicHandler(nextHandler http.Handler) http.Handler {
	fn := func(resp http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("panic: %+v", err))
				resp.WriteHeader(http.StatusInternalServerError)
			}
		}()

		nextHandler.ServeHTTP(resp, req)
	}

	return http.HandlerFunc(fn)
}
