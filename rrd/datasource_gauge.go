package rrd

import (
	"math"
	"strconv"

	"github.com/go-errors/errors"
)

const DatasourceTypeGauge = "GAUGE"

type DatasourceGauge struct {
	DatasourceAbstract
}

func (d *DatasourceGauge) CalculatePdpPrep(newValue string, interval float64) (float64, error) {
	newval, err := strconv.ParseFloat(newValue, 64)
	if err != nil {
		return math.NaN(), errors.Wrap(err, 0)
	}

	d.LastValue = newValue

	if float64(d.Heartbeat) < interval {
		return math.NaN(), nil
	}

	newPdp := newval * interval
	rate := newval

	if !d.checkRateBounds(rate) {
		newPdp = math.NaN()
	}

	return newPdp, nil
}

func (d *DatasourceGauge) DumpTo(dumper DataOutput) error {
	dumper.DumpString("type", DatasourceTypeGauge)
	return d.DatasourceAbstract.DumpTo(dumper)
}

func newDatasourceGauge(index int, store Store) (*DatasourceGauge, error) {
	result := &DatasourceGauge{}

	if err := store.ReadDatasourceParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
