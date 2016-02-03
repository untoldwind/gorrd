package rrd

type RrdDumper interface {
	Finalize() error
}
