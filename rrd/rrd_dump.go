package rrd

func (r *Rrd) DumpTo(dumper DataOutput) error {
	if err := dumper.DumpString("version", "0003"); err != nil {
		return err
	}
	if err := dumper.DumpDuration("step", r.Step); err != nil {
		return err
	}
	if err := dumper.DumpTime("lastupdate", r.LastUpdate); err != nil {
		return err
	}
	for _, datasource := range r.Datasources {
		if err := dumper.DumpSubFields("ds", func(sub DataOutput) error {
			return datasource.DumpTo(sub)
		}); err != nil {
			return err
		}
	}
	for _, rra := range r.Rras {
		if err := dumper.DumpSubFields("rra", func(sub DataOutput) error {
			return rra.DumpTo(r.Store, sub)
		}); err != nil {
			return err
		}
	}
	return dumper.Finalize()
}
