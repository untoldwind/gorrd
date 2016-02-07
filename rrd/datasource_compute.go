package rrd

const DatasourceTypeCompute = "COMPUTE"

type ComputeDatasource struct {
	Name      string
	LastValue string
}

func (d *ComputeDatasource) GetName() string {
	return d.Name
}

func (d *ComputeDatasource) SetLastValue(lastValue string) {
	d.LastValue = lastValue
}

func (d *ComputeDatasource) GetLastValue() string {
	return d.LastValue
}

func (d *ComputeDatasource) DumpTo(dumper DataDumper) error {
	if err := dumper.DumpString("name", d.Name); err != nil {
		return err
	}
	if err := dumper.DumpString("type", DatasourceTypeCompute); err != nil {
		return err
	}
	return dumper.Finalize()
}
