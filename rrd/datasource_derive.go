package rrd

import (
	"fmt"
	"math"
	"math/big"

	"github.com/go-errors/errors"
)

const DatasourceTypeDerive = "DERIVE"

type DatasourceDerive struct {
	DatasourceAbstract
}

func (d *DatasourceDerive) CalculatePdpPrep(newValue string, interval float64) (float64, error) {
	if float64(d.Heartbeat) < interval {
		d.LastValue = Undefined
	}

	rate := math.NaN()
	newPdp := math.NaN()
	if newValue != Undefined && float64(d.Heartbeat) >= interval {
		newInt := new(big.Int)
		_, err := fmt.Sscan(newValue, newInt)
		if err != nil {
			return math.NaN(), errors.Errorf("not a simple signed integer: %s", newValue)
		}
		if d.LastValue != "U" {
			prevInt := new(big.Int)
			_, err := fmt.Sscan(d.LastValue, prevInt)
			if err != nil {
				return math.NaN(), errors.Wrap(err, 0)
			}
			diff := new(big.Int)
			diff.Sub(newInt, prevInt)

			newPdp = float64(diff.Uint64())
			rate = newPdp / interval
		}
	}

	if !d.checkRateBounds(rate) {
		newPdp = math.NaN()
	}

	d.LastValue = newValue

	return newPdp, nil
}

func (d *DatasourceDerive) DumpTo(dumper DataOutput) error {
	dumper.DumpString("type", DatasourceTypeDerive)
	return d.DatasourceAbstract.DumpTo(dumper)
}

func newDatasourceDerive(index int, store Store) (*DatasourceDerive, error) {
	result := &DatasourceDerive{}

	store.ReadDatasourceParams(index, result)

	return result, nil
}
