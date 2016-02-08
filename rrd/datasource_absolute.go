package rrd

import "math"

const DatasourceTypeAbsolute = "ABSOLUTE"

type DatasourceAbsolute struct {
	DatasourceAbstract
}

func (d *DatasourceAbsolute) DumpTo(dumper DataOutput) error {
	if err := dumper.DumpString("type", DatasourceTypeAbsolute); err != nil {
		return err
	}
	return d.DatasourceAbstract.DumpTo(dumper)
}

func (d *DatasourceAbsolute) UpdatePdpPrep(newValue string, interval float64) (float64, error) {
	if newValue == "U" || float64(d.Heartbeat) < interval {
		d.LastValue = "U"
		return math.NaN(), nil
	}

	d.LastValue = newValue

	return math.NaN(), nil
}
