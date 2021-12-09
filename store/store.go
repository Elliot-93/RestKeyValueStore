// Package store. A concurrent key entries store.
package store

import (
	"RestKeyValueStore/logger"
	"errors"
)

type Key string
type Entry string

type kvsResult struct {
	entry Entry
	err   error
}

type kvsRequest struct {
	requestType string
	key         Key
	entry       Entry
	resultChan  chan kvsResult
}

type KeyValueStore struct {
	store             map[Key]Entry
	kvsRequestChannel chan kvsRequest
}

const (
	GetReq    = "Get"
	PutReq    = "Put"
	DeleteReq = "Delete"
)

var (
	ErrKeyNotFound            = errors.New("entry not present for key")
	ErrUnsupportedStoreAction = errors.New("unsupported store action requested")
)

type Store interface {
	Put(key Key, entry Entry)
	Get(key Key) (string, error)
	Delete(key Key) error
}

func New() KeyValueStore {
	kvs := KeyValueStore{
		store:             map[Key]Entry{},
		kvsRequestChannel: make(chan kvsRequest)}

	kvs.startRequestListener()

	return kvs
}

func (kvs KeyValueStore) startRequestListener() {
	go func() {
		for req := range kvs.kvsRequestChannel {
			switch req.requestType {
			case GetReq:
				entry, err := kvs.get(req.key)
				req.resultChan <- kvsResult{entry: entry, err: err}
			case PutReq:
				kvs.put(req.key, req.entry)
				req.resultChan <- kvsResult{}
			case DeleteReq:
				err := kvs.delete(req.key)
				req.resultChan <- kvsResult{err: err}
			default:
				logger.Error("")
				req.resultChan <- kvsResult{err: ErrUnsupportedStoreAction}
				return
			}
		}
	}()
}

func (kvs KeyValueStore) Put(key Key, entry Entry) {
	resultChan := make(chan kvsResult)
	kvs.kvsRequestChannel <- kvsRequest{
		requestType: PutReq,
		key:         key,
		entry:       entry,
		resultChan:  resultChan}
	<-resultChan
}

func (kvs KeyValueStore) put(key Key, entry Entry) {
	kvs.store[key] = entry
}

func (kvs KeyValueStore) Get(key Key) (string, error) {
	resultChan := make(chan kvsResult)
	kvs.kvsRequestChannel <- kvsRequest{requestType: GetReq, key: key, resultChan: resultChan}
	result := <-resultChan
	return string(result.entry), result.err
}

func (kvs KeyValueStore) get(key Key) (Entry, error) {
	entry, keyPresent := kvs.store[key]

	if !keyPresent {
		logger.Warning(ErrKeyNotFound)
		return "", ErrKeyNotFound
	}

	return entry, nil
}

func (kvs KeyValueStore) Delete(key Key) error {
	resultChan := make(chan kvsResult)
	kvs.kvsRequestChannel <- kvsRequest{
		requestType: DeleteReq,
		key:         key,
		entry:       "",
		resultChan:  resultChan}
	result := <-resultChan
	return result.err
}

func (kvs KeyValueStore) delete(key Key) error {
	_, keyPresent := kvs.store[key]

	if !keyPresent {
		logger.Warning(ErrKeyNotFound)
		return ErrKeyNotFound
	}

	delete(kvs.store, key)
	return nil
}
