package rrd

const RrdDatasourceTypeDCounter = "DCOUNTER"

type RrdDCounterDatasource struct {
	RrdDatasourceAbstract
}

func (d *RrdDCounterDatasource) DumpTo(dumper RrdDumper) error {
	if err := dumper.DumpString("type", RrdDatasourceTypeDCounter); err != nil {
		return err
	}
	return d.RrdDatasourceAbstract.DumpTo(dumper)
}
