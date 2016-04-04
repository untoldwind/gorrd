package rrd

import "math"

const RraTypeMax = "MAX"

type RraMax struct {
	RraAbstractGeneric
}

func (r *RraMax) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpString("cf", RraTypeMax)
	r.RraAbstractGeneric.DumpTo(rrdStore, dumper)
}

func newRraMax(index int, store Store) (*RraMax, error) {
	result := &RraMax{
		newRraAbstractGeneric(index, math.Inf(-1)),
	}
	result.InitializeCdpFunc = func(pdpTemp float64, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) {
		cumulativeVal := cpdPrep.Value
		if math.IsNaN(cumulativeVal) {
			cumulativeVal = math.Inf(-1)
		}
		currentVal := pdpTemp
		if math.IsNaN(currentVal) {
			currentVal = math.Inf(-1)
		}
		if currentVal > cumulativeVal {
			cpdPrep.PrimaryValue = currentVal
		} else {
			cpdPrep.PrimaryValue = cumulativeVal
		}
	}

	result.CalculateCdpValueFunc = func(pdpTemp float64, elapsedPdpSt uint64, cpdPrep *RraCpdPrepGeneric) float64 {
		if math.IsNaN(cpdPrep.Value) || pdpTemp > cpdPrep.Value {
			return pdpTemp
		}
		return cpdPrep.Value
	}

	if err := store.ReadRraParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
