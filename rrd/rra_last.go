package rrd

import "math"

const RraTypeLast = "LAST"

type RraLast struct {
	RraAbstractGeneric
}

func (r *RraLast) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpString("cf", RraTypeLast)
	r.RraAbstractGeneric.DumpTo(rrdStore, dumper)
}

func newRraLast(index int, store Store) (*RraLast, error) {
	result := &RraLast{
		newRraAbstractGeneric(index, math.NaN()),
	}
	result.InitializeCdpFunc = func(pdpTemp float64, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) {
		cpdPrep.PrimaryValue = pdpTemp
	}

	result.CalculateCdpValueFunc = func(pdpTemp float64, elapsedPdpSt uint64, cpdPrep *RraCpdPrepGeneric) float64 {
		return pdpTemp
	}

	if err := store.ReadRraParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
