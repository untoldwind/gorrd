package rrd

const RraTypeSeasonal = "SEASONAL"

type RraCpdPrepSeasonal struct {
	RraCpdPrepBase
	Seasonal     float64 `cdp:"2"`
	LastSeasonal float64 `cdp:"3"`
	InitFlag     uint64  `cdp:"6"`
}

func (c *RraCpdPrepSeasonal) DumpTo(dumper DataOutput) {
	c.RraCpdPrepBase.DumpTo(dumper)
	dumper.DumpDouble("seasonal", c.Seasonal)
	dumper.DumpDouble("last_seasonal", c.LastSeasonal)
	dumper.DumpUnsignedLong("init_flag", c.InitFlag)
}

type RraSeasonal struct {
	RraAbstract
	Gamma           float64              `rra:"param1"`
	SmoothIdx       uint64               `rra:"param4"`
	DependentRraIdx uint64               `rra:"param3"`
	CpdPreps        []RraCpdPrepSeasonal `rra:"cpdPreps"`
}

func (r *RraSeasonal) GetPrimaryValues() []float64 {
	result := make([]float64, len(r.CpdPreps))
	for i, cpdPrep := range r.CpdPreps {
		result[i] = cpdPrep.PrimaryValue
	}
	return result
}

func (r *RraSeasonal) GetSecondaryValues() []float64 {
	result := make([]float64, len(r.CpdPreps))
	for i, cpdPrep := range r.CpdPreps {
		result[i] = cpdPrep.SecondaryValue
	}
	return result
}

func (r *RraSeasonal) UpdateCdpPreps(pdpTemp []float64, elapsed ElapsedPdpSteps) uint64 {
	return 0
}

func (r *RraSeasonal) UpdateAberantCdp(pdpTemp []float64, first bool) {
}

func (r *RraSeasonal) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpString("cf", RraTypeSeasonal)
	r.dumpParams(rrdStore, dumper)
}

func (r *RraSeasonal) dumpParams(rrdStore Store, dumper DataOutput) {
	dumper.DumpUnsignedLong("pdp_per_row", r.PdpPerRow)
	dumper.DumpSubFields("params", func(params DataOutput) error {
		params.DumpDouble("seasonal_gamma", r.Gamma)
		params.DumpUnsignedLong("seasonal_smooth_idx", r.SmoothIdx)
		params.DumpUnsignedLong("dependent_rra_idx", r.DependentRraIdx)
		return nil
	})
	dumper.DumpSubFields("cdp_prep", func(cdpPreps DataOutput) error {
		for _, cdpPrep := range r.CpdPreps {
			dumper.DumpSubFields("ds", func(ds DataOutput) error {
				cdpPrep.DumpTo(ds)
				return nil
			})
		}
		return nil
	})
	r.DumpDatabase(rrdStore, dumper)
}

func newRraSeasonal(index int, store Store) (*RraSeasonal, error) {
	result := &RraSeasonal{
		RraAbstract: RraAbstract{
			Index: index,
		},
	}

	if err := store.ReadRraParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
