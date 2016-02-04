package rrd

import "time"

type Rrd struct {
	Store       RrdStore
	Step        uint64
	LastUpdate  time.Time
	Datasources []RrdDatasource
	Rras        []Rra
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
	for _, datasource := range r.Datasources {
		subDumper, err := dumper.DumpSubFields("ds")
		if err != nil {
			return err
		}
		if err := datasource.DumpTo(subDumper); err != nil {
			return err
		}
	}
	for _, rra := range r.Rras {
		subDumper, err := dumper.DumpSubFields("rra")
		if err != nil {
			return err
		}
		if err := rra.DumpTo(r.Store, subDumper); err != nil {
			return err
		}
	}
	return dumper.Finalize()
}
