package rrd

const RrdDatasourceTypeDerive = "DERIVE"

type RrdDeriveDatasource struct {
	RrdDatasourceAbstractLong
}

func (d *RrdDeriveDatasource) DumpTo(dumper RrdDumper) error {
	if err := dumper.DumpString("type", RrdDatasourceTypeDerive); err != nil {
		return err
	}
	return d.RrdDatasourceAbstractLong.DumpTo(dumper)
}
