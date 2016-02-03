package rrd

const RrdDatasourceTypeGauge = "GAUGE"

type RrdGaugeDatasource struct {
	RrdDatasourceAbstract
}

func (d *RrdGaugeDatasource) DumpTo(dumper RrdDumper) error {
	if err := dumper.DumpString("type", RrdDatasourceTypeGauge); err != nil {
		return err
	}
	return d.RrdDatasourceAbstract.DumpTo(dumper)
}
