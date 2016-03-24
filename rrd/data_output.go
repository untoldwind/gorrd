package rrd

import "time"

type DataOutput interface {
	DumpComment(comment string)
	DumpString(field, value string)
	DumpDouble(field string, value float64)
	DumpUnsignedLong(field string, value uint64)
	DumpTime(field string, value time.Time)
	DumpDuration(field string, value time.Duration)
	DumpSubFields(field string, subDump func(DataOutput) error)
	Finalize() error
}
