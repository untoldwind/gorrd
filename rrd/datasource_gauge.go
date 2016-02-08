package rrd

import "math"

const DatasourceTypeGauge = "GAUGE"

type DatasourceGauge struct {
	DatasourceAbstract
}

func (d *DatasourceGauge) UpdatePdpPrep(newValue string, interval float64) (float64, error) {
	if newValue == "U" || float64(d.Heartbeat) < interval {
		d.LastValue = "U"
		return math.NaN(), nil
	}

	d.LastValue = newValue

	return 0, nil
}

func (d *DatasourceGauge) DumpTo(dumper DataOutput) error {
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
