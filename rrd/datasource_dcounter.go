package rrd

import "math"

const DatasourceTypeDCounter = "DCOUNTER"

type DatasourceDCounter struct {
	DatasourceAbstract
}

func (d *DatasourceDCounter) UpdatePdpPrep(newValue string, interval float64) (float64, error) {
	if newValue == "U" || float64(d.Heartbeat) < interval {
		d.LastValue = "U"
		return math.NaN(), nil
	}
	d.LastValue = newValue

	return 0, nil
}

func (d *DatasourceDCounter) DumpTo(dumper DataOutput) error {
	if err := dumper.DumpString("type", DatasourceTypeDCounter); err != nil {
		return err
	}
	return d.DatasourceAbstract.DumpTo(dumper)
}
