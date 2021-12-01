package endpoints

import (
	"RestKeyValueStore/logger"
	"net/http"
)

func ReturnBadRequest(resp http.ResponseWriter, err error) {
	logger.Error(err)
	resp.WriteHeader(http.StatusBadRequest)
}

func ReturnServerError(resp http.ResponseWriter, err error) {
	logger.Error(err)
	resp.WriteHeader(http.StatusInternalServerError)
}

func ReturnForbidden(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusForbidden)
	resp.Write([]byte("Forbidden"))
}

func ReturnOK(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("OK"))
}

func ReturnOKWithBody(resp http.ResponseWriter, body string) {
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(body))
}

func ReturnOKWithBodyBytes(resp http.ResponseWriter, body []byte) {
	resp.WriteHeader(http.StatusOK)
	resp.Write(body)
}

func ReturnKeyNotFound(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusNotFound)
	resp.Write([]byte("404 key not found"))
}
