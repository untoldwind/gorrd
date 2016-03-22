package cdata

import (
	"math"
	"unsafe"
)

const rrdNaNBits = 0x7FF8000000000000

type unival uint64

func (u unival) AsDouble() float64 {
	return *(*float64)(unsafe.Pointer(&u))
}

func (u unival) AsLong() int64 {
	return int64(u)
}

func (u unival) AsUnsignedLong() uint64 {
	return uint64(u)
}

func univalForUnsignedLong(val uint64) unival {
	return unival(val)
}

func univalForDouble(val float64) unival {
	if math.IsNaN(val) {
		return unival(rrdNaNBits)
	}
	return *(*unival)(unsafe.Pointer(&val))
}
