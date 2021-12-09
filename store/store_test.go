package store

import (
	"strconv"
	"testing"
)

func TestPut(t *testing.T) {
	store := New()

	key := Key("1")
	entry := Entry("EntryOne")

	store.store["1"] = entry

	store.Put("1", entry)

	if store.store[key] != entry {
		t.Errorf("Expected entry %v at key %v", entry, key)
	}
}

func TestGet(t *testing.T) {
	store := New()

	existingKey := Key("ExistingKeyToGet")
	existingEntry := Entry("ExistingEntryToGet")
	store.store[existingKey] = existingEntry

	tests := []struct {
		name          string
		key           Key
		expectedEntry Entry
		expectedError error
	}{
		{name: "Key does exist", key: existingKey, expectedEntry: existingEntry, expectedError: nil},
		{name: "Key does not exist", key: Key("NonExistingKey"), expectedError: ErrKeyNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := store.Get(tt.key)

			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got %v", err)
			}

			if tt.expectedError != nil && tt.expectedError != err {
				t.Errorf("Expected err %v but got %v", tt.expectedError, err)
			}

			if tt.expectedEntry != "" && string(tt.expectedEntry) != entry {
				t.Errorf("Expected entry %v but got %v", tt.expectedEntry, entry)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	store := New()

	existingKey := Key("ExistingKeyToDelete")
	store.store[existingKey] = "ExistingEntry"

	tests := []struct {
		name          string
		key           Key
		expectedError error
	}{
		{name: "Key does exist", key: existingKey, expectedError: nil},
		{name: "Key does not exist", key: Key("NonExistingKey"), expectedError: ErrKeyNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Delete(tt.key)

			if tt.expectedError == nil && err != nil {
				t.Errorf("Expected no error but got %v", err)
			}

			if tt.expectedError != nil && tt.expectedError != err {
				t.Errorf("Expected err %v but got %v", tt.expectedError, err)
			} else {
				entry := store.store[tt.key]
				if entry != "" {
					t.Errorf("Expected entry to be empty string but got %v", entry)
				}
			}
		})
	}
}

func BenchmarkPut(b *testing.B) {
	store := New()

	for i := 0; i < b.N; i++ {
		store.Put(Key(strconv.Itoa(i)), "entry")
	}
}

func BenchmarkGetWith1Item(b *testing.B) {
	store := New()

	store.store["0"] = "Entry"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = store.Get("0")
	}
}

func BenchmarkGetWith1000Items(b *testing.B) {
	store := New()

	for i := 0; i < 1000; i++ {
		store.store[Key(strconv.Itoa(i))] = "entry"
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = store.Get("0")
	}
}

func BenchmarkGetWith1000000Items(b *testing.B) {
	store := New()

	for i := 0; i < 1000000; i++ {
		store.store[Key(strconv.Itoa(i))] = "entry"
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = store.Get("0")
	}
}
