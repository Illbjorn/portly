package assert

import (
	"os"
)

func EQ[T comparable](v1, v2 T, msg string) {
	if v1 != v2 {
		terminate(msg)
	}
}

func NE[T comparable](v1, v2 T, msg string) {
	if v1 == v2 {
		terminate(msg)
	}
}

type number interface {
	uint | uint8 | uint16 | uint32 | uint64 |
		int | int8 | int16 | int32 | int64 |
		float32 | float64
}

func GT[T number](v1, v2 T, msg string) {
	if !(v1 > v2) {
		terminate(msg)
	}
}

func LT[T number](v1, v2 T, msg string) {
	if !(v1 < v2) {
		terminate(msg)
	}
}

func GE[T number](v1, v2 T, msg string) {
	if !(v1 >= v2) {
		terminate(msg)
	}
}

func LE[T number](v1, v2 T, msg string) {
	if !(v1 <= v2) {
		terminate(msg)
	}
}

func terminate(msg string) {
	println(msg)
	os.Exit(1)
}
