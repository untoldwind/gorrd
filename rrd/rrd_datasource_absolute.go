package rrd

const RrdDatasourceTypeAbsolute = "ABSOLUTE"

type RrdDatasourceAbsolute struct {
	RrdDatasourceAbstract
}

func (d *RrdDatasourceAbsolute) DumpTo(dumper RrdDumper) error {
	if err := dumper.DumpString("type", RrdDatasourceTypeAbsolute); err != nil {
		return err
	}
	return d.RrdDatasourceAbstract.DumpTo(dumper)
}
