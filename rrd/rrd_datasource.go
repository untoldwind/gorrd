package rrd

type RrdDatasource interface {
	GetName() string
	SetLastValue(lastValue string)
	GetLastValue() string
	DumpTo(dumper RrdDumper) error
}

type RrdDatasourceAbstract struct {
	Name      string
	LastValue string
	Heartbeat uint64
	Min       float64
	Max       float64
}

func (d *RrdDatasourceAbstract) GetName() string {
	return d.Name
}

func (d *RrdDatasourceAbstract) SetLastValue(lastValue string) {
	d.LastValue = lastValue
}
func (d *RrdDatasourceAbstract) GetLastValue() string {
	return d.LastValue
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
