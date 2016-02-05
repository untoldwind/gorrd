package rrd

const RrdDatasourceTypeCompute = "COMPUTE"

type RrdComputeDatasource struct {
	Name      string
	LastValue string
}

func (d *RrdComputeDatasource) GetName() string {
	return d.Name
}

func (d *RrdComputeDatasource) SetLastValue(lastValue string) {
	d.LastValue = lastValue
}

func (d *RrdComputeDatasource) GetLastValue() string {
	return d.LastValue
}

func (d *RrdComputeDatasource) DumpTo(dumper RrdDumper) error {
	if err := dumper.DumpString("name", d.Name); err != nil {
		return err
	}
	if err := dumper.DumpString("type", RrdDatasourceTypeCompute); err != nil {
		return err
	}
	return dumper.Finalize()
}
