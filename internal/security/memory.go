package security

import "unsafe"

func SecureZero(b []byte) {
	if len(b) == 0 {
		return
	}
	ptr := unsafe.Pointer(&b[0])
	for i := range b {
		*(*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(i))) = 0
	}
}


