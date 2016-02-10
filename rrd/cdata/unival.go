package cdata

import "unsafe"

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
	return *(*unival)(unsafe.Pointer(&val))
}
