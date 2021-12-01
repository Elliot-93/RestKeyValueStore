package endpoints

import (
	"RestKeyValueStore/authentication"
	"RestKeyValueStore/server"
	"errors"
	"net/http"
)

type ShutdownHandler struct{}

const ShutdownRoute = "/shutdown"

var ErrNonAdminRequestedShutdown = errors.New("non admin requested shutdown")

func (h ShutdownHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	user := req.Context().Value(authentication.AuthUsernameCtxKey)

	switch req.Method {
	case http.MethodGet:
		if user == authentication.Admin {
			ReturnOK(resp)
			server.Shutdown()
		} else {
			ReturnForbidden(resp, ErrNonAdminRequestedShutdown)
			return
		}
	default:
		resp.WriteHeader(http.StatusNotFound)
	}
}
