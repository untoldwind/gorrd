package rrd

import "math"

const RraTypeSeasonal = "SEASONAL"

type RraCpdPrepSeasonal struct {
	RraCpdPrepBase
	Seasonal     float64 `cdp:"2"`
	LastSeasonal float64 `cdp:"3"`
	InitFlat     uint64  `cdp:"6"`
}
type RraSeasonal struct {
	RraAbstractGeneric
}

func (r *RraSeasonal) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpString("cf", RraTypeSeasonal)
	r.RraAbstractGeneric.DumpTo(rrdStore, dumper)
}

func newRraSeasonal(index int, store Store) (*RraSeasonal, error) {
	result := &RraSeasonal{
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
