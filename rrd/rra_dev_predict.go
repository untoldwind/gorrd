package rrd

const RraTypeDevPredict = "DEVPREDICT"

type RraDevPredict struct {
	RraAbstract
	DependentRraIdx uint64           `rra:"param3"`
	CpdPreps        []RraCpdPrepBase `rra:"cpdPreps"`
}

func (r *RraDevPredict) GetPrimaryValues() []float64 {
	result := make([]float64, len(r.CpdPreps))
	for i, cpdPrep := range r.CpdPreps {
		result[i] = cpdPrep.PrimaryValue
	}
	return result
}

func (r *RraDevPredict) GetSecondaryValues() []float64 {
	result := make([]float64, len(r.CpdPreps))
	for i, cpdPrep := range r.CpdPreps {
		result[i] = cpdPrep.SecondaryValue
	}
	return result
}

func (r *RraDevPredict) UpdateCdpPreps(pdpTemp []float64, elapsed ElapsedPdpSteps) uint64 {
	return 0
}

func (r *RraDevPredict) UpdateAberantCdp(pdpTemp []float64, first bool) {
}

func (r *RraDevPredict) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpString("cf", RraTypeDevPredict)
	dumper.DumpUnsignedLong("pdp_per_row", r.PdpPerRow)
	dumper.DumpSubFields("params", func(params DataOutput) error {
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

func newRraDevPredict(index int, store Store) (*RraDevPredict, error) {
	result := &RraDevPredict{
		RraAbstract: RraAbstract{
			Index: index,
		},
	}

	if err := store.ReadRraParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
