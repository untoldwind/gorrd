package rrd

type Rra interface {
	DumpTo(dumper RrdDumper) error
}
