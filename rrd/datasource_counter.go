package rrd

import (
	"math"

	"github.com/go-errors/errors"
)

const DatasourceTypeCounter = "COUNTER"

type DatasourceCounter struct {
	DatasourceAbstract
}

func (d *DatasourceCounter) UpdatePdpPrep(newValue string, interval float64) (float64, error) {
	if float64(d.Heartbeat) < interval {
		d.LastValue = "U"
		return math.NaN(), nil
	}
	newPdp := math.NaN()
	if newValue != "U" && float64(d.Heartbeat) >= interval {
		if !rrdIsUnsignedInt(newValue) {
			return math.NaN(), errors.Errorf("not a simple unsigned integer: %s", newValue)
		}
		if d.LastValue != "U" {

		}
	}

	d.LastValue = newValue

	return newPdp, nil
}

func (d *DatasourceCounter) DumpTo(dumper DataOutput) error {
	if err := dumper.DumpString("type", DatasourceTypeCounter); err != nil {
		return err
	}
	return d.DatasourceAbstract.DumpTo(dumper)
}

func newDatasourceCounter(index int, store Store) (*DatasourceCounter, error) {
	result := &DatasourceCounter{}

	if err := store.ReadDatasourceParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
