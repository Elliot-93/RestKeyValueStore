package endpoints

import (
	"net/http"
)

type PingHandler struct{}

const PingRoute = "/ping"

func (h PingHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")

	switch req.Method {
	case http.MethodGet:
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte("pong"))
	default:
		resp.WriteHeader(http.StatusNotFound)
	}
}
