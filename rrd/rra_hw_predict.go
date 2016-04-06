package rrd

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

func (c *RraCpdPrepHwPredict) DumpTo(dumper DataOutput) {
	c.RraCpdPrepBase.DumpTo(dumper)
	dumper.DumpDouble("intercept", c.Intercept)
	dumper.DumpDouble("last_intercept", c.LastIntercept)
	dumper.DumpDouble("slope", c.Slope)
	dumper.DumpDouble("last_slope", c.LastSlope)
	dumper.DumpUnsignedLong("nan_count", c.NullCount)
	dumper.DumpUnsignedLong("last_nan_count", c.LastNullCount)
}

type RraHwPredict struct {
	RraAbstract
	Alpha           float64               `rra:"param1"`
	Beta            float64               `rra:"param2"`
	DependentRraIdx uint64                `rra:"param3"`
	CpdPreps        []RraCpdPrepHwPredict `rra:"cpdPreps"`
}

func (r *RraHwPredict) GetPrimaryValues() []float64 {
	result := make([]float64, len(r.CpdPreps))
	for i, cpdPrep := range r.CpdPreps {
		result[i] = cpdPrep.PrimaryValue
	}
	return result
}

func (r *RraHwPredict) GetSecondaryValues() []float64 {
	result := make([]float64, len(r.CpdPreps))
	for i, cpdPrep := range r.CpdPreps {
		result[i] = cpdPrep.SecondaryValue
	}
	return result
}

func (r *RraHwPredict) UpdateCdpPreps(pdpTemp []float64, elapsed ElapsedPdpSteps) uint64 {
	return 0
}

func (r *RraHwPredict) UpdateAberantCdp(pdpTemp []float64, first bool) {
}

func (r *RraHwPredict) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpString("cf", RraTypeHwPredict)
	dumper.DumpUnsignedLong("pdp_per_row", r.PdpPerRow)
	dumper.DumpSubFields("params", func(params DataOutput) error {
		params.DumpDouble("hw_alpha", r.Alpha)
		params.DumpDouble("hw_beta", r.Beta)
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

func newRraHwPredict(index int, store Store) (*RraHwPredict, error) {
	result := &RraHwPredict{
		RraAbstract: RraAbstract{
			Index: index,
		},
	}

	if err := store.ReadRraParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
