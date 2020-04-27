package danger

import (
	"reflect"
	"unsafe"
)

// BytesToString turns a []byte into a string with 0 MemAllocs and 0 MemBytes.
func BytesToString(bytes []byte) (s string) {
	if len(bytes) == 0 {
		return s
	}
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh.Data = uintptr(unsafe.Pointer(&bytes[0]))
	sh.Len = len(bytes)
	return s
}

// StringToBytes turns a string into a []byte with 0 MemAllocs and 0 MemBytes.
func StringToBytes(s string) (b []byte) {
	if len(s) == 0 {
		return b
	}
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{Data: sh.Data, Len: sh.Len, Cap: sh.Len}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
