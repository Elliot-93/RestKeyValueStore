// Package store a concurrent key entries store in which each key has an owner.
// Only the owner of a key can modify that entries.
package store

//todo: multiple queues with different request types funneled into 1?

import (
	"RestKeyValueStore/logger"
	"errors"
)

const (
	GetReq             = "Get"
	GetAllSummariesReq = "GetAllSummaries"
	GetSummaryReq      = "GetSummary"
	PutReq             = "Put"
	DeleteReq          = "Delete"
)

type kvsResult struct {
	entry          Entry
	entrySummaries []EntrySummary
	err            error
}

type kvsRequest struct {
	requestType   string
	key           Key
	entry         Entry
	adminOverride bool
	resultChan    chan kvsResult
}

type Key string

type keyValueStore struct {
	store             map[Key]Entry
	kvsRequestChannel chan kvsRequest
	initialized       bool
}

var (
	kvsStore                  keyValueStore
	ErrKeyNotFound            = errors.New("entries not present for key")
	ErrKeyBelongsToOtherUser  = errors.New("key added by a different owner")
	ErrUnsupportedStoreAction = errors.New("unsupported store action requested")
)

func lazyInit() {
	if !kvsStore.initialized {
		kvsStore.store = map[Key]Entry{}
		kvsStore.kvsRequestChannel = make(chan kvsRequest)
		kvsStore.initialized = true
		requestListener()
	}
}

func requestListener() {
	go func() {
		for req := range kvsStore.kvsRequestChannel {
			switch req.requestType {
			case GetReq:
				entry, err := get(req.key)
				req.resultChan <- kvsResult{entry: entry, err: err}
			case PutReq:
				err := put(req.key, req.entry, req.adminOverride)
				req.resultChan <- kvsResult{err: err}
			case DeleteReq:
				err := deleteFromStore(req.key, req.entry.Owner, req.adminOverride)
				req.resultChan <- kvsResult{err: err}
			case GetAllSummariesReq:
				entrySummaries := getAllEntrySummaries()
				req.resultChan <- kvsResult{entrySummaries: entrySummaries}
			case GetSummaryReq:
				entrySummary, err := getEntrySummary(req.key)
				req.resultChan <- kvsResult{entrySummaries: []EntrySummary{entrySummary}, err: err}
			default:
				logger.Error("")
				req.resultChan <- kvsResult{err: ErrUnsupportedStoreAction}
				return
			}
		}
	}()
}

func Put(key Key, entry Entry, adminOverride bool) error {
	lazyInit()
	resultChan := make(chan kvsResult)
	kvsStore.kvsRequestChannel <- kvsRequest{
		requestType:   PutReq,
		key:           key,
		entry:         entry,
		adminOverride: adminOverride,
		resultChan:    resultChan}
	result := <-resultChan
	return result.err
}

func put(key Key, entry Entry, adminOverride bool) error {
	existingEntry, keyPresent := kvsStore.store[key]

	if keyPresent && adminOverride {
		entry.Owner = existingEntry.Owner
		kvsStore.store[key] = entry
		return nil
	}

	if keyPresent && entry.Owner != existingEntry.Owner {
		return ErrKeyBelongsToOtherUser
	}

	kvsStore.store[key] = entry
	return nil
}

func Get(key Key) (string, error) {
	lazyInit()
	resultChan := make(chan kvsResult)
	kvsStore.kvsRequestChannel <- kvsRequest{requestType: GetReq, key: key, resultChan: resultChan}
	result := <-resultChan
	return result.entry.Value, result.err
}

func get(key Key) (Entry, error) {
	entry, keyPresent := kvsStore.store[key]

	if !keyPresent {
		return Entry{}, ErrKeyNotFound
	}

	return entry, nil
}

func Delete(key Key, owner string, adminOverride bool) error {
	lazyInit()
	resultChan := make(chan kvsResult)
	kvsStore.kvsRequestChannel <- kvsRequest{
		requestType:   DeleteReq,
		key:           key,
		entry:         Entry{Owner: owner},
		adminOverride: adminOverride,
		resultChan:    resultChan}
	result := <-resultChan
	return result.err
}

func deleteFromStore(key Key, owner string, adminOverride bool) error {
	entry, keyPresent := kvsStore.store[key]

	if !keyPresent {
		return ErrKeyNotFound
	}

	if entry.Owner != owner && !adminOverride {
		return ErrKeyBelongsToOtherUser
	}

	delete(kvsStore.store, key)
	return nil
}

func GetAllEntrySummaries() []EntrySummary {
	lazyInit()
	resultChan := make(chan kvsResult)
	kvsStore.kvsRequestChannel <- kvsRequest{requestType: GetAllSummariesReq, resultChan: resultChan}
	result := <-resultChan
	return result.entrySummaries
}

func getAllEntrySummaries() []EntrySummary {
	lazyInit()
	var entrySummaries []EntrySummary

	for k, v := range kvsStore.store {
		entrySummaries = append(entrySummaries, EntrySummary{Key: string(k), Owner: v.Owner})
	}

	return entrySummaries
}

func GetEntrySummary(key Key) (EntrySummary, error) {
	lazyInit()
	resultChan := make(chan kvsResult)
	kvsStore.kvsRequestChannel <- kvsRequest{requestType: GetSummaryReq, key: key, resultChan: resultChan}
	result := <-resultChan

	if result.entrySummaries != nil {
		return result.entrySummaries[0], result.err
	} else {
		return EntrySummary{}, result.err
	}
}

func getEntrySummary(key Key) (EntrySummary, error) {
	lazyInit()
	entry, keyPresent := kvsStore.store[key]

	if !keyPresent {
		return EntrySummary{}, ErrKeyNotFound
	}

	return EntrySummary{Key: string(key), Owner: entry.Owner}, nil
}
