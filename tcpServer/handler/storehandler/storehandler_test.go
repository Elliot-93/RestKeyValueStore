package storehandler

import (
	"RestKeyValueStore/store"
	"bytes"
	"io"
	"testing"
)

type spyStore struct {
	putInvocations    int
	getResult         string
	getErr            error
	getInvocations    int
	deleteErr         error
	deleteInvocations int
}

func (s *spyStore) Put(key store.Key, entry store.Entry) {
	s.putInvocations++
}

func (s *spyStore) Get(key store.Key) (string, error) {
	s.getInvocations++
	return s.getResult, s.getErr
}

func (s *spyStore) Delete(key store.Key) error {
	s.deleteInvocations++
	return s.deleteErr
}

func TestHandlePut(t *testing.T) {
	type args struct {
		r io.Reader
		s spyStore
	}
	tests := []struct {
		name                   string
		args                   args
		want                   string
		expectedPutInvocations int
	}{
		{
			name:                   "no key returns err",
			args:                   args{r: bytes.NewBuffer([]byte("")), s: spyStore{}},
			want:                   "err",
			expectedPutInvocations: 0,
		},
		{
			name:                   "Key p1 not an int returns err",
			args:                   args{r: bytes.NewBuffer([]byte("a")), s: spyStore{}},
			want:                   "err",
			expectedPutInvocations: 0,
		},
		{
			name:                   "Key p2 not provided returns err",
			args:                   args{r: bytes.NewBuffer([]byte("1")), s: spyStore{}},
			want:                   "err",
			expectedPutInvocations: 0,
		},
		{
			name:                   "Key p2 not an int returns err",
			args:                   args{r: bytes.NewBuffer([]byte("1a")), s: spyStore{}},
			want:                   "err",
			expectedPutInvocations: 0,
		},
		{
			name:                   "Key p3 not provided returns err",
			args:                   args{r: bytes.NewBuffer([]byte("11")), s: spyStore{}},
			want:                   "err",
			expectedPutInvocations: 0,
		},
		{
			name:                   "Value in invalid format",
			args:                   args{r: bytes.NewBuffer([]byte("11kinvalid")), s: spyStore{}},
			want:                   "err",
			expectedPutInvocations: 0,
		},
		{
			name:                   "Valid command",
			args:                   args{r: bytes.NewBuffer([]byte("11k11v")), s: spyStore{}},
			want:                   "ack",
			expectedPutInvocations: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HandlePut(tt.args.r, &tt.args.s); got != tt.want {
				t.Errorf("HandlePut() = %v, want %v", got, tt.want)
			}

			if tt.expectedPutInvocations != tt.args.s.putInvocations {
				t.Errorf("PutInvocations = %v, want %v", tt.args.s.putInvocations, tt.expectedPutInvocations)
			}
		})
	}
}

func TestHandleGet(t *testing.T) {
	type args struct {
		r io.Reader
		s spyStore
	}
	tests := []struct {
		name                   string
		args                   args
		want                   string
		expectedGetInvocations int
	}{
		{
			name:                   "Key in invalid format",
			args:                   args{r: bytes.NewBuffer([]byte("invalid")), s: spyStore{}},
			want:                   "err",
			expectedGetInvocations: 0,
		},
		{
			name:                   "Get returns error return nil",
			args:                   args{r: bytes.NewBuffer([]byte("11k11v")), s: spyStore{getErr: store.ErrKeyNotFound}},
			want:                   "nil",
			expectedGetInvocations: 1,
		},
		{
			name:                   "Valid command",
			args:                   args{r: bytes.NewBuffer([]byte("11k11v")), s: spyStore{getResult: "value"}},
			want:                   "val15value",
			expectedGetInvocations: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HandleGet(tt.args.r, &tt.args.s); got != tt.want {
				t.Errorf("HandleGet() = %v, want %v", got, tt.want)
			}

			if tt.expectedGetInvocations != tt.args.s.getInvocations {
				t.Errorf("GetInvocations = %v, want %v", tt.args.s.putInvocations, tt.expectedGetInvocations)
			}
		})
	}
}

func TestHandleDelete(t *testing.T) {
	type args struct {
		r io.Reader
		s spyStore
	}
	tests := []struct {
		name                      string
		args                      args
		want                      string
		expectedDeleteInvocations int
	}{
		{
			name:                      "Key in invalid format",
			args:                      args{r: bytes.NewBuffer([]byte("invalid")), s: spyStore{}},
			want:                      "err",
			expectedDeleteInvocations: 0,
		},
		{
			name:                      "Delete returns err still ack",
			args:                      args{r: bytes.NewBuffer([]byte("11k11v")), s: spyStore{deleteErr: store.ErrKeyNotFound}},
			want:                      "ack",
			expectedDeleteInvocations: 1,
		},
		{
			name:                      "Valid command",
			args:                      args{r: bytes.NewBuffer([]byte("11k11v")), s: spyStore{}},
			want:                      "ack",
			expectedDeleteInvocations: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HandleDelete(tt.args.r, &tt.args.s); got != tt.want {
				t.Errorf("HandleDelete() = %v, want %v", got, tt.want)
			}

			if tt.expectedDeleteInvocations != tt.args.s.deleteInvocations {
				t.Errorf("DeleteInvocations = %v, want %v", tt.args.s.putInvocations, tt.expectedDeleteInvocations)
			}
		})
	}
}
