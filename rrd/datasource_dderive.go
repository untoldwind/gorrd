package rrd

import "math"

const DatasourceTypeDDerive = "DDERIVE"

type DatasourceDDerive struct {
	DatasourceAbstract
}

func (d *DatasourceDDerive) UpdatePdpPrep(newValue string, interval float64) (float64, error) {
	if newValue == "U" || float64(d.Heartbeat) < interval {
		d.LastValue = "U"
		return math.NaN(), nil
	}
	d.LastValue = newValue

	return 0, nil
}

func (d *DatasourceDDerive) DumpTo(dumper DataOutput) error {
	if err := dumper.DumpString("type", DatasourceTypeDDerive); err != nil {
		return err
	}
	return d.DatasourceAbstract.DumpTo(dumper)
}
