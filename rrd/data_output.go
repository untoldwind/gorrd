package rrd

import "time"

type DataOutput interface {
	DumpComment(comment string) error
	DumpString(field, value string) error
	DumpDouble(field string, value float64) error
	DumpUnsignedLong(field string, value uint64) error
	DumpTime(field string, value time.Time) error
	DumpDuration(field string, value time.Duration) error
	DumpSubFields(field string, subDump func(DataOutput) error) error
	Finalize() error
}
