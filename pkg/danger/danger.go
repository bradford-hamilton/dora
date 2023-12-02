package danger

import (
	"unsafe"
)

// BytesToString turns a []byte into a string with 0 MemAllocs and 0 MemBytes.
// This is an unsafe operation and may lead to problems if the bytes passed in
// are changed while the string is used. No checking whether bytes are valid
// UTF-8 data is performed.
func BytesToString(bytes []byte) (s string) {
	if len(bytes) == 0 {
		return s
	}
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}

// StringToBytes turns a string into a []byte with 0 MemAllocs and 0 MemBytes.
// This is an unsafe operation and will lead to problems if the underlying bytes
// are changed.
func StringToBytes(s string) (b []byte) {
	if len(s) == 0 {
		return b
	}
	const max = 0x7fff0000 // 2147418112
	if len(s) > max {
		panic("string too large")
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
