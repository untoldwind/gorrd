package rrd

import (
	"fmt"
	"math"
	"math/big"

	"github.com/go-errors/errors"
)

const DatasourceTypeCounter = "COUNTER"

type DatasourceCounter struct {
	DatasourceAbstract
}

func (d *DatasourceCounter) CalculatePdpPrep(newValue string, interval float64) (float64, error) {
	if float64(d.Heartbeat) < interval {
		d.LastValue = Undefined
	}
	rate := math.NaN()
	newPdp := math.NaN()
	if newValue != Undefined && float64(d.Heartbeat) >= interval {
		newInt := new(big.Int)
		_, err := fmt.Sscan(newValue, newInt)
		if err != nil || newInt.Sign() < 0 {
			return math.NaN(), errors.Errorf("not a simple unsigned integer: %s", newValue)
		}
		if d.LastValue != "U" {
			prevInt := new(big.Int)
			_, err := fmt.Sscan(d.LastValue, prevInt)
			if err != nil {
				return math.NaN(), errors.Wrap(err, 0)
			}
			diff := new(big.Int)
			diff.Sub(newInt, prevInt)
			// Handle overflow
			if diff.Sign() < 0 {
				diff.Add(diff, big.NewInt(math.MaxUint32))
			}
			if diff.Sign() < 0 {
				diff.Add(diff, new(big.Int).SetUint64(math.MaxUint64-math.MaxUint32))
			}
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
