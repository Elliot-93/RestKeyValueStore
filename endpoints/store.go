package endpoints

import (
	"RestKeyValueStore/authentication"
	"RestKeyValueStore/store"
	"errors"
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

	err := store.Put(key, store.Entry{Value: valStringBuilder.String(), Owner: user}, strings.EqualFold(user, authentication.Admin))

	if err != nil {
		switch {
		case errors.Is(err, store.ErrKeyBelongsToOtherUser):
			ReturnForbidden(resp, err)
		default:
			ReturnServerError(resp, err)
		}
		return
	}

	ReturnOK(resp)
}

func handleStoreGet(resp http.ResponseWriter, req *http.Request, key store.Key) {
	valStringBuilder := new(strings.Builder)
	_, errReadBody := io.Copy(valStringBuilder, req.Body)

	if errReadBody != nil {
		ReturnBadRequest(resp, errReadBody)
		return
	}

	value, err := store.Get(key)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrKeyNotFound):
			ReturnKeyNotFound(resp, err)
		default:
			ReturnServerError(resp, err)
		}
	}

	ReturnOKWithBody(resp, value)
}

func handleStoreDelete(resp http.ResponseWriter, key store.Key, user string) {
	err := store.Delete(key, user, strings.EqualFold(user, authentication.Admin))

	if err != nil {
		switch {
		case errors.Is(err, store.ErrKeyNotFound):
			ReturnKeyNotFound(resp, err)
		case errors.Is(err, store.ErrKeyBelongsToOtherUser):
			ReturnForbidden(resp, err)
		default:
			ReturnServerError(resp, err)
		}
	}

	ReturnOK(resp)
}
