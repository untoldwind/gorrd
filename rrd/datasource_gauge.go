package rrd

const DatasourceTypeGauge = "GAUGE"

type DatasourceGauge struct {
	DatasourceAbstract
}

func (d *DatasourceGauge) DumpTo(dumper DataDumper) error {
	if err := dumper.DumpString("type", DatasourceTypeGauge); err != nil {
		return err
	}
	return d.DatasourceAbstract.DumpTo(dumper)
}

func newDatasourceGauge(index int, store Store) (*DatasourceGauge, error) {
	result := &DatasourceGauge{}

	if err := store.ReadDatasourceParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
