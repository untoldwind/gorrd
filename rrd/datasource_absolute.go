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

func (d *DatasourceAbsolute) UpdatePdpPrep(newValue string, interval float64) (float64, error) {
	newval, err := strconv.ParseFloat(newValue, 64)
	if err != nil {
		return math.NaN(), errors.Wrap(err, 0)
	}

	newPdp := newval
	rate := newval / interval

	if !d.checkRateBounds(rate) {
		newPdp = math.NaN()
	}

	d.LastValue = newValue

	return newPdp, nil
}

func (d *DatasourceAbsolute) DumpTo(dumper DataOutput) error {
	if err := dumper.DumpString("type", DatasourceTypeAbsolute); err != nil {
		return err
	}
	return d.DatasourceAbstract.DumpTo(dumper)
}
