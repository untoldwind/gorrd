package rrd

import "fmt"

const RraTypeFailures = "FAILURES"

type RraCpdPrepFailures struct {
	RraCpdPrepBase
	History []byte `cdp:"raw"`
}

func (c *RraCpdPrepFailures) DumpToWitHistory(windowLen uint64, dumper DataOutput) {
	c.RraCpdPrepBase.DumpTo(dumper)
	historyString := ""
	for i := uint64(0); i < windowLen; i++ {
		historyString += fmt.Sprintf("%d", c.History[i])
	}
	dumper.DumpString("history", historyString)
}

type RraHwFailures struct {
	RraAbstract
	DeltaPos         float64              `rra:"param1"`
	DeltaNeg         float64              `rra:"param2"`
	WindowLen        uint64               `rra:"param4"`
	FailureThreshold uint64               `rra:"param5"`
	DependentRraIdx  uint64               `rra:"param3"`
	CpdPreps         []RraCpdPrepFailures `rra:"cpdPreps"`
}

func (r *RraHwFailures) GetPrimaryValues() []float64 {
	result := make([]float64, len(r.CpdPreps))
	for i, cpdPrep := range r.CpdPreps {
		result[i] = cpdPrep.PrimaryValue
	}
	return result
}

func (r *RraHwFailures) GetSecondaryValues() []float64 {
	result := make([]float64, len(r.CpdPreps))
	for i, cpdPrep := range r.CpdPreps {
		result[i] = cpdPrep.SecondaryValue
	}
	return result
}

func (r *RraHwFailures) UpdateCdpPreps(pdpTemp []float64, elapsed ElapsedPdpSteps) uint64 {
	return 0
}

func (r *RraHwFailures) UpdateAberantCdp(pdpTemp []float64, first bool) {
}

func (r *RraHwFailures) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpString("cf", RraTypeFailures)
	dumper.DumpUnsignedLong("pdp_per_row", r.PdpPerRow)
	dumper.DumpSubFields("params", func(params DataOutput) error {
		params.DumpDouble("delta_pos", r.DeltaPos)
		params.DumpDouble("delta_neg", r.DeltaNeg)
		params.DumpUnsignedLong("window_len", r.WindowLen)
		params.DumpUnsignedLong("failure_threshold", r.FailureThreshold)
		params.DumpUnsignedLong("dependent_rra_idx", r.DependentRraIdx)
		return nil
	})
	dumper.DumpSubFields("cdp_prep", func(cdpPreps DataOutput) error {
		for _, cdpPrep := range r.CpdPreps {
			dumper.DumpSubFields("ds", func(ds DataOutput) error {
				cdpPrep.DumpToWitHistory(r.WindowLen, ds)
				return nil
			})
		}
		return nil
	})
	r.DumpDatabase(rrdStore, dumper)
}

func newRraFailures(index int, store Store) (*RraHwFailures, error) {
	result := &RraHwFailures{
		RraAbstract: RraAbstract{
			Index: index,
		},
	}

	if err := store.ReadRraParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
