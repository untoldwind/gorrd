package rrd

import (
	"math"

	"github.com/go-errors/errors"
)

const DatasourceTypeDerive = "DERIVE"

type DatasourceDerive struct {
	DatasourceAbstract
}

func (d *DatasourceDerive) UpdatePdpPrep(newValue string, interval float64) (float64, error) {
	if float64(d.Heartbeat) < interval {
		d.LastValue = "U"
		return math.NaN(), nil
	}
	newPdp := math.NaN()
	if newValue != "U" && float64(d.Heartbeat) >= interval {
		if !rrdIsSignedInt(newValue) {
			return math.NaN(), errors.Errorf("not a simple signed integer: %s", newValue)
		}
		if d.LastValue != "U" {
		}
	}

	d.LastValue = newValue

	return newPdp, nil
}

func (d *DatasourceDerive) DumpTo(dumper DataOutput) error {
	if err := dumper.DumpString("type", DatasourceTypeDerive); err != nil {
		return err
	}
	return d.DatasourceAbstract.DumpTo(dumper)
}
