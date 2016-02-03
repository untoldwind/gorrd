package rrd

import "time"

type RrdDumper interface {
	DumpString(field, value string) error
	DumpUnsignedLong(field string, value uint64) error
	DumpTime(field string, value time.Time) error
	Finalize() error
}
