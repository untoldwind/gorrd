package rrd

const DatasourceTypeDCounter = "DCOUNTER"

type DatasourceDCounter struct {
	DatasourceAbstract
}

func (d *DatasourceDCounter) DumpTo(dumper DataDumper) error {
	if err := dumper.DumpString("type", DatasourceTypeDCounter); err != nil {
		return err
	}
	return d.DatasourceAbstract.DumpTo(dumper)
}
