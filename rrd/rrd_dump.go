package rrd

import "fmt"

func (r *Rrd) DumpTo(dumper DataOutput) error {
	dumper.DumpString("version", fmt.Sprintf("%04d", r.Version))
	dumper.DumpDuration("step", r.Step)
	dumper.DumpTime("lastupdate", r.LastUpdate)
	for _, datasource := range r.Datasources {
		dumper.DumpSubFields("ds", func(sub DataOutput) error {
			return datasource.DumpTo(sub)
		})
	}
	for _, rra := range r.Rras {
		dumper.DumpSubFields("rra", func(sub DataOutput) error {
			rra.DumpTo(r.Store, sub)
			return nil
		})
	}
	return dumper.Finalize()
}
