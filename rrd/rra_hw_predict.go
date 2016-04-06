package rrd

import "math"

const RraTypeHwPredict = "HWPREDICT"

type RraCpdPrepHwPredict struct {
	RraCpdPrepBase
	Intercept     float64 `cdp:"2"`
	LastIntercept float64 `cdp:"3"`
	Slope         float64 `cdp:"4"`
	LastSlope     float64 `cdp:"5"`
	NullCount     uint64  `cdp:"6"`
	LastNullCount uint64  `cdp:"7"`
}

type RraHwPredict struct {
	RraAbstractGeneric
}

func (r *RraHwPredict) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpString("cf", RraTypeHwPredict)
	r.RraAbstractGeneric.DumpTo(rrdStore, dumper)
}

func newRraHwPredict(index int, store Store) (*RraHwPredict, error) {
	result := &RraHwPredict{
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
