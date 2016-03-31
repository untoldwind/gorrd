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
		RraAbstractGeneric: RraAbstractGeneric{
			Index: index,
			ResetCpdFunc: func(pdpTemp float64, cpdPrep *RraCpdPrepGeneric) {
				cpdPrep.PrimaryValue = pdpTemp
				cpdPrep.SecondaryValue = pdpTemp
			},
			InitializeCdpFunc: func(pdpTemp float64, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) {
				cumulativeVal := cpdPrep.Value
				if math.IsNaN(cumulativeVal) {
					cumulativeVal = 0
				}
				currentVal := pdpTemp
				if math.IsNaN(currentVal) {
					currentVal = 0
				}
				cpdPrep.PrimaryValue = (cumulativeVal + currentVal*float64(startPdpOffset)) / float64(pdpPerRow-cpdPrep.UnknownDatapoints)
			},
			InitializeCarryOverFunc: func(pdpTemp float64, elapsedPdpSt, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) float64 {
				pdpIntoCdpCnt := (elapsedPdpSt - startPdpOffset) % pdpPerRow
				if pdpIntoCdpCnt == 0 || math.IsNaN(pdpTemp) {
					return 0
				}
				return pdpTemp * float64(pdpIntoCdpCnt)
			},
			CalculateCdpValueFunc: func(pdpTemp float64, elapsedPdpSt uint64, cpdPrep *RraCpdPrepGeneric) float64 {
				if math.IsNaN(cpdPrep.Value) {
					return pdpTemp * float64(elapsedPdpSt)
				}
				return cpdPrep.Value + pdpTemp*float64(elapsedPdpSt)
			},
		},
	}

	if err := store.ReadRraParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
