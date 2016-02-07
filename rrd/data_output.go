package rrd

import "time"

type DataDumper interface {
	DumpComment(comment string) error
	DumpString(field, value string) error
	DumpDouble(field string, value float64) error
	DumpUnsignedLong(field string, value uint64) error
	DumpTime(field string, value time.Time) error
	DumpSubFields(field string, subDump func(DataDumper) error) error
	Finalize() error
}
