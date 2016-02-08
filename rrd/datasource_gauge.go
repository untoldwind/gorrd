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

func (d *DatasourceGauge) UpdatePdpPrep(newValue string, interval float64) (float64, error) {
	newval, err := strconv.ParseFloat(newValue, 64)
	if err != nil {
		return math.NaN(), errors.Wrap(err, 0)
	}

	newPdp := newval * interval
	rate := newval

	if !d.checkRateBounds(rate) {
		newPdp = math.NaN()
	}

	d.LastValue = newValue

	return newPdp, nil
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
