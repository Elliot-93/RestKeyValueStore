package endpoints

import (
	"RestKeyValueStore/authentication"
	"RestKeyValueStore/store"
	"io"
	"net/http"
	"strings"
)

type StoreHandler struct{}

const StoreRoute = "/store/"

func (h StoreHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	key := store.Key(strings.TrimPrefix(req.URL.Path, StoreRoute))
	user := req.Context().Value(authentication.AuthUsernameCtxKey).(string)

	switch req.Method {
	case http.MethodPut:
		handleStorePut(resp, req, key, user)
	case http.MethodGet:
		handleStoreGet(resp, req, key)
	case http.MethodDelete:
		handleStoreDelete(resp, key, user)
	default:
		resp.WriteHeader(http.StatusNotFound)
	}
}

func handleStorePut(resp http.ResponseWriter, req *http.Request, key store.Key, user string) {
	valStringBuilder := new(strings.Builder)
	_, errReadBody := io.Copy(valStringBuilder, req.Body)

	if errReadBody != nil {
		ReturnBadRequest(resp, errReadBody)
		return
	}

	err := store.Put(key, store.Entry{Value: valStringBuilder.String(), Owner: user})

	switch err {
	case nil:
		ReturnOK(resp)
	case store.ErrKeyBelongsToOtherUser:
		ReturnForbidden(resp)
	default:
		ReturnServerError(resp, err)
	}
}

func handleStoreGet(resp http.ResponseWriter, req *http.Request, key store.Key) {
	valStringBuilder := new(strings.Builder)
	_, errReadBody := io.Copy(valStringBuilder, req.Body)

	if errReadBody != nil {
		ReturnBadRequest(resp, errReadBody)
		return
	}

	value, err := store.Get(key)

	switch err {
	case nil:
		ReturnOKWithBody(resp, value)
	case store.ErrKeyNotFound:
		ReturnKeyNotFound(resp)
	default:
		ReturnServerError(resp, err)
	}
}

func handleStoreDelete(resp http.ResponseWriter, key store.Key, user string) {
	err := store.Delete(key, user)

	switch err {
	case nil:
		ReturnOK(resp)
	case store.ErrKeyNotFound:
		ReturnKeyNotFound(resp)
	case store.ErrKeyBelongsToOtherUser:
		ReturnForbidden(resp)
	default:
		ReturnServerError(resp, err)
	}
}
