package rrd

import "math"

type RraCpdPrepGeneric struct {
	RraCpdPrepBase
	Value             float64 `cdp:"0"`
	UnknownDatapoints uint64  `cdp:"1"`
}

func (r *RraCpdPrepGeneric) Reset(pdp float64) {
	r.PrimaryValue = pdp
	r.SecondaryValue = pdp
}

func (c *RraCpdPrepGeneric) DumpTo(dumper DataOutput) {
	dumper.DumpSubFields("ds", func(ds DataOutput) error {
		ds.DumpDouble("primary_value", c.PrimaryValue)
		ds.DumpDouble("secondary_value", c.SecondaryValue)
		ds.DumpDouble("value", c.Value)
		ds.DumpUnsignedLong("unknown_datapoints", c.UnknownDatapoints)
		return nil
	})
}

type RraAbstractGeneric struct {
	RraAbstract
	XFilesFactor            float64             `rra:"param0"`
	CpdPreps                []RraCpdPrepGeneric `rra:"cpdPreps"`
	UpdateAberantCdpFunc    func(pdpTemp float64, cpdPrep *RraCpdPrepGeneric)
	InitializeCdpFunc       func(pdpTemp float64, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric)
	InitializeCarryOverFunc func(pdpTemp float64, elapsedPdpSt, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) float64
	CalculateCdpValueFunc   func(pdpTemp float64, elapsedPdpSt uint64, cpdPrep *RraCpdPrepGeneric) float64
}

func newRraAbstractGeneric(index int, initialCarryOver float64) RraAbstractGeneric {
	return RraAbstractGeneric{
		RraAbstract: RraAbstract{
			Index: index,
		},
		InitializeCarryOverFunc: func(pdpTemp float64, elapsedPdpSt, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) float64 {
			pdpIntoCdpCnt := (elapsedPdpSt - startPdpOffset) % pdpPerRow
			if pdpIntoCdpCnt == 0 || math.IsNaN(pdpTemp) {
				return initialCarryOver
			}
			return pdpTemp
		},
		CalculateCdpValueFunc: func(pdpTemp float64, elapsedPdpSt uint64, cpdPrep *RraCpdPrepGeneric) float64 {
			return pdpTemp
		},
	}
}

func (r *RraAbstractGeneric) GetRowCount() uint64 {
	return r.RowCount
}

func (r *RraAbstractGeneric) GetPdpPerRow() uint64 {
	return r.PdpPerRow
}

func (r *RraAbstractGeneric) GetPrimaryValues() []float64 {
	result := make([]float64, len(r.CpdPreps))
	for i, cpdPrep := range r.CpdPreps {
		result[i] = cpdPrep.PrimaryValue
	}
	return result
}

func (r *RraAbstractGeneric) GetSecondaryValues() []float64 {
	result := make([]float64, len(r.CpdPreps))
	for i, cpdPrep := range r.CpdPreps {
		result[i] = cpdPrep.SecondaryValue
	}
	return result
}

func (r *RraAbstractGeneric) UpdateCdpPreps(pdpTemp []float64, elapsed ElapsedPdpSteps) uint64 {
	startPdpOffset := r.PdpPerRow - elapsed.ProcPdpCount%r.PdpPerRow
	var rraStepCount uint64
	if startPdpOffset <= elapsed.Steps {
		rraStepCount = minUInt64((elapsed.Steps-startPdpOffset)/r.PdpPerRow+1, r.RowCount)
	}
	if r.PdpPerRow > 1 {
		if rraStepCount > 0 {
			for i, pdp := range pdpTemp {
				if math.IsNaN(pdp) {
					r.CpdPreps[i].UnknownDatapoints += startPdpOffset
					r.CpdPreps[i].SecondaryValue = math.NaN()
				} else {
					r.CpdPreps[i].SecondaryValue = pdp
				}

				if float64(r.CpdPreps[i].UnknownDatapoints) > float64(r.PdpPerRow)*r.XFilesFactor {
					r.CpdPreps[i].PrimaryValue = math.NaN()
				} else {
					r.InitializeCdpFunc(pdp, r.PdpPerRow, startPdpOffset, &r.CpdPreps[i])
				}

				r.CpdPreps[i].Value = r.InitializeCarryOverFunc(pdp, elapsed.Steps, r.PdpPerRow, startPdpOffset, &r.CpdPreps[i])

				if math.IsNaN(pdp) {
					r.CpdPreps[i].UnknownDatapoints = (elapsed.Steps - startPdpOffset) % r.PdpPerRow
				} else {
					r.CpdPreps[i].UnknownDatapoints = 0
				}
			}
		} else {
			for i, pdp := range pdpTemp {
				if math.IsNaN(pdp) {
					r.CpdPreps[i].UnknownDatapoints += elapsed.Steps
				} else {
					r.CpdPreps[i].Value = r.CalculateCdpValueFunc(pdp, elapsed.Steps, &r.CpdPreps[i])
				}
			}
		}
	} else {
		// There is just one PDP pre CDP
		if elapsed.Steps > 2 {
			for i, pdp := range pdpTemp {
				r.CpdPreps[i].Reset(pdp)
			}
		}
	}
	return rraStepCount
}

func (r *RraAbstractGeneric) UpdateAberantCdp(pdpTemp []float64, first bool) {
	if r.PdpPerRow != 1 {
		return
	}
	for i, pdp := range pdpTemp {
		if first {
			r.CpdPreps[i].PrimaryValue = pdp
		} else {
			r.CpdPreps[i].SecondaryValue = pdp
		}
		if r.UpdateAberantCdpFunc != nil {
			r.UpdateAberantCdpFunc(pdp, &r.CpdPreps[i])
		}
	}
}

func (r *RraAbstractGeneric) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpUnsignedLong("pdp_per_row", r.PdpPerRow)
	dumper.DumpSubFields("params", func(params DataOutput) error {
		params.DumpDouble("xff", r.XFilesFactor)
		return nil
	})
	dumper.DumpSubFields("cdp_prep", func(cdpPreps DataOutput) error {
		for _, cdpPrep := range r.CpdPreps {
			cdpPrep.DumpTo(cdpPreps)
		}
		return nil
	})
	dumper.DumpSubFields("database", func(database DataOutput) error {
		rowIterator, err := rrdStore.RowIterator(r.Index)
		if err != nil {
			return err
		}
		return ForEachRow(rowIterator, func(row *RraRow) error {
			row.DumpTo(dumper)
			return nil
		})
	})
}

func minUInt64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
