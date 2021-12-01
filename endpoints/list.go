package endpoints

import (
	"RestKeyValueStore/store"
	"encoding/json"
	"net/http"
	"strings"
)

type ListHandler struct{}

const ListRoute = "/list/"

func (h ListHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	key := store.Key(strings.TrimPrefix(req.URL.Path, ListRoute))

	switch req.Method {
	case http.MethodGet:
		if key == "" {
			handleList(resp)
		} else {
			handleListKey(resp, key)
		}
	default:
		resp.WriteHeader(http.StatusNotFound)
	}
}

func handleList(resp http.ResponseWriter) {
	entrySummaries := store.GetAllEntrySummaries()

	marshalledResponse, err := json.Marshal(entrySummaries)

	if err != nil {
		ReturnServerError(resp, err)
	} else {
		ReturnOKWithBodyBytes(resp, marshalledResponse)
	}
}

func handleListKey(resp http.ResponseWriter, key store.Key) {
	entrySummary, err := store.GetEntrySummary(key)

	if err == store.ErrKeyNotFound {
		ReturnKeyNotFound(resp)
		return
	} else if err != nil {
		ReturnServerError(resp, err)
		return
	}

	marshalledResponse, marshalErr := json.Marshal(entrySummary)

	if marshalErr != nil {
		ReturnServerError(resp, marshalErr)
	} else {
		ReturnOKWithBodyBytes(resp, marshalledResponse)
	}
}
