package rrd

type RrdDatasource interface {
	GetName() string
	DumpTo(dumper RrdDumper) error
}

type RrdDatasourceAbstract struct {
	Name      string
	Heartbeat uint64
	Min       float64
	Max       float64
}

func (d *RrdDatasourceAbstract) GetName() string {
	return d.Name
}

func (d *RrdDatasourceAbstract) DumpTo(dumper RrdDumper) error {
	if err := dumper.DumpString("name", d.Name); err != nil {
		return err
	}
	if err := dumper.DumpUnsignedLong("minimal_heartbeat", d.Heartbeat); err != nil {
		return err
	}
	if err := dumper.DumpDouble("min", d.Min); err != nil {
		return err
	}
	return dumper.DumpDouble("max", d.Max)
}
