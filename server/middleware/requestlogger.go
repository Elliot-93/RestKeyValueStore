package middleware

import (
	"RestKeyValueStore/logger"
	"fmt"
	"net/http"
	"time"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func LoggingMiddleware(nextHandler http.Handler) http.Handler {
	fn := func(resp http.ResponseWriter, req *http.Request) {

		requestTime := time.Now().UTC()

		wrappedResponse := wrapResponseWriter(resp)

		nextHandler.ServeHTTP(wrappedResponse, req)

		logger.LogRequest(fmt.Sprintf(
			"Request logged TimestampUtc: %v \n\tMethod: %s \n\tURL: %s \n\tRemoteIP: %s "+
				"\n\tResponseCode: %d \n\tDuration: %v\n",
			requestTime.Format(time.RFC3339),
			req.Method,
			req.URL.String(),
			req.RemoteAddr,
			wrappedResponse.status,
			time.Since(requestTime)))
	}

	return http.HandlerFunc(fn)
}
