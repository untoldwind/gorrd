package rrd

import "math"

const RraTypeDevPredict = "DEVPREDICT"

type RraDevPredict struct {
	RraAbstractGeneric
}

func (r *RraDevPredict) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpString("cf", RraTypeDevPredict)
	r.RraAbstractGeneric.DumpTo(rrdStore, dumper)
}

func newRraDevPredict(index int, store Store) (*RraDevPredict, error) {
	result := &RraDevPredict{
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
