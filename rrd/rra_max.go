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
		RraAbstractGeneric: RraAbstractGeneric{
			Index: index,
			ResetCpdFunc: func(pdpTemp float64, cpdPrep *RraCpdPrepGeneric) {
				cpdPrep.PrimaryValue = pdpTemp
				cpdPrep.SecondaryValue = pdpTemp
			},
			InitializeCdpFunc: func(pdpTemp float64, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) {
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
			},
			InitializeCarryOverFunc: func(pdpTemp float64, elapsedPdpSt, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) float64 {
				pdpIntoCdpCnt := (elapsedPdpSt - startPdpOffset) % pdpPerRow
				if pdpIntoCdpCnt == 0 || math.IsNaN(pdpTemp) {
					return math.Inf(-1)
				}
				return pdpTemp
			},
			CalculateCdpValueFunc: func(pdpTemp float64, elapsedPdpSt uint64, cpdPrep *RraCpdPrepGeneric) float64 {
				if math.IsNaN(cpdPrep.Value) || pdpTemp > cpdPrep.Value {
					return pdpTemp
				}
				return cpdPrep.Value
			},
		},
	}

	if err := store.ReadRraParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
