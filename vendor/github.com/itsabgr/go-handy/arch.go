package handy

import (
	"unsafe"
)

//MaxUint is max uint value
const MaxUint = ^uint(0)

//MaxInt is max int value
const MaxInt = int(MaxUint >> 1)

//UintptrSize is memory size of uintptr
const UintptrSize = unsafe.Sizeof(uintptr(0))

//Is64 return true only if system is 64 bit
func Is64() bool {
	return ^uint(0)>>63 == 1
}
