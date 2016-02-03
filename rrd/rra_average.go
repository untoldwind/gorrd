package rrd

const RraTypeAverage = "AVERAGE"

type RraAverage struct {
	RowCount     uint64
	PdpCount     uint64
	XFilesFactor float64
}
