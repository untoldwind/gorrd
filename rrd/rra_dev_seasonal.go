package rrd

import "math"

const RraTypeDevSeasonal = "DEVSEASONAL"

type RraDevSeasonal struct {
	RraAbstractGeneric
}

func (r *RraDevSeasonal) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpString("cf", RraTypeDevSeasonal)
	r.RraAbstractGeneric.DumpTo(rrdStore, dumper)
}

func newRraDevSeasonal(index int, store Store) (*RraDevSeasonal, error) {
	result := &RraDevSeasonal{
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
