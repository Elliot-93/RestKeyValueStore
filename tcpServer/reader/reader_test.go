package reader_test

import (
	"RestKeyValueStore/tcpServer/reader"
	"bytes"
	"io"
	"testing"
)

func TestReadBytes(t *testing.T) {
	type args struct {
		bb     *bytes.Buffer
		length int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name:    "Read returns error",
			args:    args{bb: bytes.NewBuffer(make([]byte, 0)), length: 3},
			want:    string(make([]byte, 3)),
			wantErr: io.EOF,
		},
		{
			name:    "Read returns bytes as string",
			args:    args{bb: bytes.NewBuffer([]byte("bytes")), length: 3},
			want:    "byt",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reader.ReadBytes(tt.args.bb, tt.args.length)
			if err != tt.wantErr {
				t.Errorf("ReadBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkReadBytes(b *testing.B) {
	tenByteSlice := []byte("0123456789")

	b.Run("10 bytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = reader.ReadBytes(bytes.NewBuffer(tenByteSlice), 10)
		}
	})
}
