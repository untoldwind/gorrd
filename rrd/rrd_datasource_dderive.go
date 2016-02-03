package rrd

const RrdDatasourceTypeDDerive = "DDERIVE"

type RrdDDeriveDatasource struct {
	RrdDatasourceAbstract
}

func (d *RrdDDeriveDatasource) DumpTo(dumper RrdDumper) error {
	if err := dumper.DumpString("type", RrdDatasourceTypeDDerive); err != nil {
		return err
	}
	return d.RrdDatasourceAbstract.DumpTo(dumper)
}
