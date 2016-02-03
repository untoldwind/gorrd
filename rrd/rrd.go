package rrd

type Rrd struct {
	Store RrdStore
}

func (r *Rrd) Close() {
	r.Store.Close()
}

func (r *Rrd) DumpTo(dumper RrdDumper) error {
	if err := dumper.DumpString("version", "0003"); err != nil {
		return err
	}
	return dumper.Finalize()
}
