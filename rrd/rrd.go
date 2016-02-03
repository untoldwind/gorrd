package rrd

import "time"

type Rrd struct {
	Store      RrdStore
	Step       uint64
	LastUpdate time.Time
}

func (r *Rrd) Close() {
	r.Store.Close()
}

func (r *Rrd) DumpTo(dumper RrdDumper) error {
	if err := dumper.DumpString("version", "0003"); err != nil {
		return err
	}
	if err := dumper.DumpUnsignedLong("step", r.Step); err != nil {
		return err
	}
	if err := dumper.DumpTime("lastupdate", r.LastUpdate); err != nil {
		return err
	}
	return dumper.Finalize()
}
