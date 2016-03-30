package rrd

import "math"

type Rra interface {
	GetRowCount() uint64
	GetPdpPerRow() uint64
	GetPrimaryValues() []float64
	GetSecondaryValues() []float64
	UpdateCdpPreps(pdpTemp []float64, elapsedSteps, procPdpCount uint64) (uint64, error)
	UpdateAberantCdp(pdpTemp []float64, first bool) error
	DumpTo(rrdStore Store, dumper DataOutput)
}

type RraCpdPrepGeneric struct {
	PrimaryValue      float64 `cdp:"8"`
	SecondaryValue    float64 `cdp:"9"`
	Value             float64 `cdp:"0"`
	UnknownDatapoints uint64  `cdp:"1"`
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
	Index                   int
	RowCount                uint64              `rra:"rowCount"`
	PdpPerRow               uint64              `rra:"pdpPerRow"`
	XFilesFactor            float64             `rra:"param0"`
	CpdPreps                []RraCpdPrepGeneric `rra:"cpdPreps"`
	ResetCpdFunc            func(pdpTemp float64, cpdPrep *RraCpdPrepGeneric) error
	UpdateAberantCdpFunc    func(pdpTemp float64, cpdPrep *RraCpdPrepGeneric) error
	InitializeCdpFunc       func(pdpTemp float64, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) error
	InitializeCarryOverFunc func(pdpTemp float64, elapsedPdpSt, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) (float64, error)
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

func (r *RraAbstractGeneric) UpdateCdpPreps(pdpTemp []float64, elapsedSteps, procPdpCount uint64) (uint64, error) {
	startPdpOffset := r.PdpPerRow - procPdpCount%r.PdpPerRow
	var rraStepCount uint64
	if startPdpOffset <= elapsedSteps {
		rraStepCount = minUInt64((elapsedSteps-startPdpOffset)/r.PdpPerRow+1, r.RowCount)
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

				var err error
				r.CpdPreps[i].Value, err = r.InitializeCarryOverFunc(pdp, elapsedSteps, r.PdpPerRow, startPdpOffset, &r.CpdPreps[i])
				if err != nil {
					return 0, err
				}

				if math.IsNaN(pdp) {
					r.CpdPreps[i].UnknownDatapoints = (elapsedSteps - startPdpOffset) % r.PdpPerRow
				} else {
					r.CpdPreps[i].UnknownDatapoints = 0
				}
			}
		}
	} else {
		// There is just one PDP pre CDP
		if elapsedSteps > 2 {
			for i, pdp := range pdpTemp {
				if err := r.ResetCpdFunc(pdp, &r.CpdPreps[i]); err != nil {
					return 0, err
				}
			}
			return rraStepCount, nil
		}
	}
	return rraStepCount, nil
}

func (r *RraAbstractGeneric) UpdateAberantCdp(pdpTemp []float64, first bool) error {
	for i, pdp := range pdpTemp {
		if first {
			r.CpdPreps[i].PrimaryValue = pdp
		} else {
			r.CpdPreps[i].SecondaryValue = pdp
		}
		if r.UpdateAberantCdpFunc != nil {
			if err := r.UpdateAberantCdpFunc(pdp, &r.CpdPreps[i]); err != nil {
				return err
			}
		}
	}

	return nil
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

func newRra(index int, rraType string, store Store) (Rra, error) {
	switch rraType {
	case RraTypeAverage:
		return newRraAverage(index, store)
	}
	return nil, nil
}

func minUInt64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
