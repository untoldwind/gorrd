package rrd

import (
	"math"
	"strconv"

	"github.com/go-errors/errors"
)

const DatasourceTypeDCounter = "DCOUNTER"

type DatasourceDCounter struct {
	DatasourceAbstract
}

func (d *DatasourceDCounter) UpdatePdpPrep(newValue string, interval float64) (float64, error) {
	if float64(d.Heartbeat) < interval {
		d.LastValue = "U"
	}
	rate := math.NaN()
	newPdp := math.NaN()
	if newValue != "U" && float64(d.Heartbeat) >= interval {
		newval, err := strconv.ParseFloat(newValue, 64)
		if err != nil {
			return math.NaN(), errors.Wrap(err, 0)
		}

		oldval, err := strconv.ParseFloat(d.LastValue, 64)
		if err != nil {
			return math.NaN(), errors.Wrap(err, 0)
		}

		if (newval > 0 && oldval < newval) ||
			(newval < 0 && newval < oldval) {
			newPdp = newval - oldval
			rate = newPdp / interval
		}
	}

	if !d.checkRateBounds(rate) {
		newPdp = math.NaN()
	}

	d.LastValue = newValue

	return newPdp, nil
}

func (d *DatasourceDCounter) DumpTo(dumper DataOutput) error {
	if err := dumper.DumpString("type", DatasourceTypeDCounter); err != nil {
		return err
	}
	return d.DatasourceAbstract.DumpTo(dumper)
}
