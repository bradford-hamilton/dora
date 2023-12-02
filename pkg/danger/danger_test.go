package danger

import (
	"bytes"
	"strings"
	"testing"
)

func TestBytesToString(t *testing.T) {
	tests := [...]struct {
		name  string
		bytes []byte
		want  string
	}{
		{
			name:  "When there are no bytes, it returns an empty string",
			bytes: []byte{},
			want:  "",
		},
		{
			name:  "When there are bytes, it returns the string representation",
			bytes: []byte{'b', 'r', 'a', 'd', 'f', 'o', 'r', 'd'},
			want:  "bradford",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotS := BytesToString(tt.bytes); gotS != tt.want {
				t.Errorf("BytesToString(): got: %v, want: %v", gotS, tt.want)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want []byte
	}{
		{
			name: "When the string is empty, it returns an empty byte slice",
			str:  "",
			want: []byte{},
		},
		{
			name: "When the string is != \"\", it returns the corresponding byte slice",
			str:  "bradford",
			want: []byte{'b', 'r', 'a', 'd', 'f', 'o', 'r', 'd'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotB := StringToBytes(tt.str); !bytes.Equal(gotB, tt.want) != false {
				t.Errorf("StringToBytes() = %v, want %v", gotB, tt.want)
			}
		})
	}
}

var stringSink string
var byteSink []byte
var byt = []byte(strings.Repeat("Me trysail Jack Ketch Sink me measured fer yer chains long boat hornswaggle quarter brig black spot careen Admiral of the Black!", 10))
var str = strings.Repeat("Me trysail Jack Ketch Sink me measured fer yer chains long boat hornswaggle quarter brig black spot careen Admiral of the Black!", 10)

func BenchmarkStringToBytesSafe(b *testing.B) {
	b.SetBytes(int64(len(str)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := string(byt)
		stringSink = v
	}
}

func BenchmarkBytesToStringSafe(b *testing.B) {
	b.SetBytes(int64(len(byt)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := []byte(str)
		byteSink = v
	}
}

func BenchmarkStringToBytesDanger(b *testing.B) {
	b.SetBytes(int64(len(str)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := StringToBytes(str)
		byteSink = v
	}
}

func BenchmarkBytesToStringDanger(b *testing.B) {
	b.SetBytes(int64(len(byt)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := BytesToString(byt)
		stringSink = v
	}
}

// Results:
// BenchmarkStringToBytesSafe-16      	 6326526	       190 ns/op	6720.33 MB/s	    1280 B/op	       1 allocs/op
// BenchmarkBytesToStringSafe-16      	 6035268	       197 ns/op	6481.43 MB/s	    1280 B/op	       1 allocs/op
// BenchmarkStringToBytesDanger-16    	1000000000	       1.20 ns/op	1068924.12 MB/s	       0 B/op	       0 allocs/op
// BenchmarkBytesToStringDanger-16    	991642945	       1.21 ns/op	1057479.20 MB/s	       0 B/op	       0 allocs/op
