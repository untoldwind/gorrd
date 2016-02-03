package rrd

const RrdDatasourceTypeCounter = "COUNTER"

type RrdCounterDatasource struct {
	RrdDatasourceAbstract
}

func (d *RrdCounterDatasource) DumpTo(dumper RrdDumper) error {
	if err := dumper.DumpString("type", RrdDatasourceTypeCounter); err != nil {
		return err
	}
	return d.RrdDatasourceAbstract.DumpTo(dumper)
}
