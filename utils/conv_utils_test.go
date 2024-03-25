package utils

import (
	"fmt"
	"testing"
)

func TestByte2string(t *testing.T) {
	val := []byte("hello world")

	fmt.Println(Byte2string(val))
}

func TestStr2byte(t *testing.T) {
	str := "hello world"
	fmt.Printf("%s\n", string(Str2byte(str)))
}
