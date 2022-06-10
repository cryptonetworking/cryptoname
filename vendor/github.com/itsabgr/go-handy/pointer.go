package handy

import "unsafe"

//Ref returns a uintptr points to the value
func Ref[T any](value *T) uintptr {
	return uintptr(unsafe.Pointer(value))
}

//DeRef returns the value pointer points to
func DeRef[T any](ptr uintptr) T {
	return *((*T)(unsafe.Pointer(ptr)))
}
