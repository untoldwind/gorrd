package rrd

import "math"

const RraTypeAverage = "AVERAGE"

type RraAverage struct {
	RraAbstractGeneric
}

func (r *RraAverage) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpString("cf", RraTypeAverage)
	r.RraAbstractGeneric.DumpTo(rrdStore, dumper)
}

func newRraAverage(index int, store Store) (*RraAverage, error) {
	result := &RraAverage{
		newRraAbstractGeneric(index, 0),
	}
	result.InitializeCdpFunc = func(pdpTemp float64, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) {
		cumulativeVal := cpdPrep.Value
		if math.IsNaN(cumulativeVal) {
			cumulativeVal = 0
		}
		currentVal := pdpTemp
		if math.IsNaN(currentVal) {
			currentVal = 0
		}
		cpdPrep.PrimaryValue = (cumulativeVal + currentVal*float64(startPdpOffset)) / float64(pdpPerRow-cpdPrep.UnknownDatapoints)
	}

	result.CalculateCdpValueFunc = func(pdpTemp float64, elapsedPdpSt uint64, cpdPrep *RraCpdPrepGeneric) float64 {
		if math.IsNaN(cpdPrep.Value) {
			return pdpTemp * float64(elapsedPdpSt)
		}
		return cpdPrep.Value + pdpTemp*float64(elapsedPdpSt)
	}

	if err := store.ReadRraParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
