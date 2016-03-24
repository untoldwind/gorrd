package rrd

import (
	"math"
	"strconv"

	"github.com/go-errors/errors"
)

const DatasourceTypeAbsolute = "ABSOLUTE"

type DatasourceAbsolute struct {
	DatasourceAbstract
}

func (d *DatasourceAbsolute) CalculatePdpPrep(newValue string, interval float64) (float64, error) {
	newval, err := strconv.ParseFloat(newValue, 64)
	if err != nil {
		return math.NaN(), errors.Wrap(err, 0)
	}

	d.LastValue = newValue

	if interval > float64(d.Heartbeat) {
		return math.NaN(), nil
	}

	newPdp := newval
	rate := newval / interval

	if !d.checkRateBounds(rate) {
		newPdp = math.NaN()
	}

	return newPdp, nil
}

func (d *DatasourceAbsolute) DumpTo(dumper DataOutput) error {
	dumper.DumpString("type", DatasourceTypeAbsolute)
	return d.DatasourceAbstract.DumpTo(dumper)
}
