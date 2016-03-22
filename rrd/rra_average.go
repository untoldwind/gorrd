package rrd

const RraTypeAverage = "AVERAGE"

type RraAverage struct {
	RraAbstractGeneric
}

func (r *RraAverage) DumpTo(rrdStore Store, dumper DataOutput) error {
	if err := dumper.DumpString("cf", RraTypeAverage); err != nil {
		return err
	}
	return r.RraAbstractGeneric.DumpTo(rrdStore, dumper)
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
		},
	}

	if err := store.ReadRraParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
