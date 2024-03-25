package utils

import (
	"reflect"
	"unsafe"
)

func Byte2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Str2byte(s string) (b []byte) {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))

	sliceHeader.Data = stringHeader.Data
	sliceHeader.Len = stringHeader.Len
	sliceHeader.Cap = stringHeader.Len
	return
}
