package cdata

import "math"

type unival uint64

func (u unival) AsDouble() float64 {
	return math.Float64frombits(uint64(u))
}

func (u unival) AsLong() int64 {
	return int64(u)
}

func (u unival) AsUnsignedLong() uint64 {
	return uint64(u)
}
