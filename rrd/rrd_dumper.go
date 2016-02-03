package rrd

type RrdDumper interface {
	DumpString(field, value string) error
	Finalize() error
}
