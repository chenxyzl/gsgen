package mdata

import "math"

func kToIdx[K comparable](v K) uint64 {
	var d interface{} = v
	switch x := d.(type) {
	case int:
		return uint64(x)
	case uint:
		return uint64(x)
	case int8:
		return uint64(x)
	case uint8:
		return uint64(x)
	case int16:
		return uint64(x)
	case uint16:
		return uint64(x)
	case int32:
		return uint64(x)
	case uint32:
		return uint64(x)
	case int64:
		return uint64(x)
	case uint64:
		return x
	default:
		return math.MaxUint64
	}
}

func isNum[K comparable](v K) bool {
	var d interface{} = v
	switch d.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		return true
	default:
		return false
	}
}
