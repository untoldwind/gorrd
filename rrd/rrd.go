package rrd

type Rrd struct {
	Store RrdStore
}

func (r *Rrd) Close() {
	r.Store.Close()
}

func (r *Rrd) DumpTo(dumper RrdDumper) error {
	return dumper.Finalize()
}
