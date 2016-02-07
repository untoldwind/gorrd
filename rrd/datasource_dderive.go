package rrd

const DatasourceTypeDDerive = "DDERIVE"

type DatasourceDDerive struct {
	DatasourceAbstract
}

func (d *DatasourceDDerive) DumpTo(dumper DataDumper) error {
	if err := dumper.DumpString("type", DatasourceTypeDDerive); err != nil {
		return err
	}
	return d.DatasourceAbstract.DumpTo(dumper)
}
