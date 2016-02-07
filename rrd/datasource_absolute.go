package rrd

const DatasourceTypeAbsolute = "ABSOLUTE"

type DatasourceAbsolute struct {
	DatasourceAbstract
}

func (d *DatasourceAbsolute) DumpTo(dumper DataDumper) error {
	if err := dumper.DumpString("type", DatasourceTypeAbsolute); err != nil {
		return err
	}
	return d.DatasourceAbstract.DumpTo(dumper)
}
