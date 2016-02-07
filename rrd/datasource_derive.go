package rrd

const DatasourceTypeDerive = "DERIVE"

type DatasourceDerive struct {
	DatasourceAbstract
}

func (d *DatasourceDerive) DumpTo(dumper DataDumper) error {
	if err := dumper.DumpString("type", DatasourceTypeDerive); err != nil {
		return err
	}
	return d.DatasourceAbstract.DumpTo(dumper)
}
