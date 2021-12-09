package tcpreader_test

import (
	"RestKeyValueStore/tcpServer/tcpreader"
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
			tcpReader := tcpreader.New(tt.args.bb)
			got, err := tcpReader.ReadBytes(tt.args.length)
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

func TestParse3PartArgument(t *testing.T) {
	tests := []struct {
		name          string
		ioReader      io.Reader
		want          string
		expectedError error
	}{
		{
			name:          "no bytes returns err",
			ioReader:      bytes.NewBuffer([]byte("")),
			want:          "",
			expectedError: tcpreader.ErrReadingPartOneOfArg,
		},
		{
			name:          "p1 not an int returns err",
			ioReader:      bytes.NewBuffer([]byte("a")),
			want:          "",
			expectedError: tcpreader.ErrParsingPartOneOfArg,
		},
		{
			name:          "p2 not provided returns err",
			ioReader:      bytes.NewBuffer([]byte("1")),
			want:          "",
			expectedError: tcpreader.ErrReadingPartTwoOfArg,
		},
		{
			name:          "p2 not an int returns err",
			ioReader:      bytes.NewBuffer([]byte("1a")),
			want:          "",
			expectedError: tcpreader.ErrParsingPartTwoOfArg,
		},
		{
			name:          "p3 not provided returns err",
			ioReader:      bytes.NewBuffer([]byte("11")),
			want:          "",
			expectedError: tcpreader.ErrReadingPartThreeOfArg,
		},
		{
			name:          "Valid command",
			ioReader:      bytes.NewBuffer([]byte("11k")),
			want:          "k",
			expectedError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcpReader := tcpreader.New(tt.ioReader)
			got, err := tcpReader.Parse3PartArgument()
			if got != tt.want {
				t.Errorf("Parse3PartArgument value = %v, want %v", got, tt.want)
			}

			if err != tt.expectedError {
				t.Errorf("Parse3PartArgument error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}

func TestParseResponseLengthArgument(t *testing.T) {
	tests := []struct {
		name          string
		ioReader      io.Reader
		want          int
		expectedError error
	}{
		{
			name:          "empty byte array provided",
			ioReader:      bytes.NewBuffer([]byte("")),
			want:          0,
			expectedError: tcpreader.ErrReadingPartOneOfArg,
		},
		{
			name:          "p1 not an int returns err",
			ioReader:      bytes.NewBuffer([]byte("a")),
			want:          0,
			expectedError: tcpreader.ErrParsingPartOneOfArg,
		},
		{
			name:          "p2 not provided returns err",
			ioReader:      bytes.NewBuffer([]byte("1")),
			want:          0,
			expectedError: tcpreader.ErrReadingPartTwoOfArg,
		},
		{
			name:          "p2 not an int returns err",
			ioReader:      bytes.NewBuffer([]byte("1a")),
			want:          0,
			expectedError: tcpreader.ErrParsingPartTwoOfArg,
		},
		{
			name:          "0 provided, 0 is returned with no error",
			ioReader:      bytes.NewBuffer([]byte("0")),
			want:          0,
			expectedError: nil,
		},
		{
			name:          "Valid command",
			ioReader:      bytes.NewBuffer([]byte("12")),
			want:          2,
			expectedError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcpReader := tcpreader.New(tt.ioReader)
			got, err := tcpReader.ParseResponseLengthArg()
			if got != tt.want {
				t.Errorf("Parse3PartArgument value = %v, want %v", got, tt.want)
			}

			if err != tt.expectedError {
				t.Errorf("Parse3PartArgument error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}

func TestParseVerb(t *testing.T) {
	expectedResult := "put"
	tcpReader := tcpreader.New(bytes.NewBuffer([]byte(expectedResult)))
	got, err := tcpReader.ParseVerb()
	if err != nil {
		t.Errorf("ParseVerb() error = %v, wantErr %v", err, nil)
		return
	}
	if got != expectedResult {
		t.Errorf("ParseVerb() got = %v, want %v", got, expectedResult)
	}
}

func BenchmarkReadBytes(b *testing.B) {
	tenByteSlice := []byte("0123456789")
	tcpReader := tcpreader.New(bytes.NewBuffer(tenByteSlice))

	b.ResetTimer()

	b.Run("10 bytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = tcpReader.ReadBytes(10)
		}
	})
}
