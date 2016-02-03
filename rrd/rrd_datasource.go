package rrd

type RrdDatasource interface {
	GetName() string
}

type RrdDatasourceAbstractLong struct {
	Name      string
	Heartbeat uint64
	Min       uint64
	Max       uint64
}

func (d *RrdDatasourceAbstractLong) GetName() string {
	return d.Name
}

type RrdDatasourceAbstractDouble struct {
	Name      string
	Heartbeat uint64
	Min       float64
	Max       float64
}

func (d *RrdDatasourceAbstractDouble) GetName() string {
	return d.Name
}
