package utils

import (
	"reflect"
	"strconv"
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

func ToUint64(b []byte) uint64 {
	s := Byte2string(b)
	i, _ := strconv.ParseUint(s, 10, 64)
	return i
}

func FromUint64(i uint64) string {
	s := strconv.FormatUint(i, 10)
	return s
}
