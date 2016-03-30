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
			ResetCpdFunc: func(pdpTemp float64, cpdPrep *RraCpdPrepGeneric) error {
				cpdPrep.PrimaryValue = pdpTemp
				cpdPrep.SecondaryValue = pdpTemp
				return nil
			},
			InitializeCdpFunc: func(pdpTemp float64, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) error {
				cumulativeVal := cpdPrep.Value
				if math.IsNaN(cumulativeVal) {
					cumulativeVal = 0
				}
				currentVal := pdpTemp
				if math.IsNaN(currentVal) {
					currentVal = 0
				}
				cpdPrep.PrimaryValue = cumulativeVal + currentVal*float64(startPdpOffset)/float64(pdpPerRow-cpdPrep.UnknownDatapoints)
				return nil
			},
			InitializeCarryOverFunc: func(pdpTemp float64, elapsedPdpSt, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) (float64, error) {
				pdpIntoCdpCnt := (elapsedPdpSt - startPdpOffset) % pdpPerRow
				if pdpIntoCdpCnt == 0 || math.IsNaN(pdpTemp) {
					return 0, nil
				}
				return pdpTemp * float64(pdpIntoCdpCnt), nil
			},
			CalculateCdpValueFunc: func(pdpTemp float64, elapsedPdpSt uint64, cpdPrep *RraCpdPrepGeneric) (float64, error) {
				if math.IsNaN(cpdPrep.Value) {
					return pdpTemp * float64(elapsedPdpSt), nil
				}
				return cpdPrep.Value + pdpTemp*float64(elapsedPdpSt), nil
			},
		},
	}

	if err := store.ReadRraParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
