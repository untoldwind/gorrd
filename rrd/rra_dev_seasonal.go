package rrd

const RraTypeDevSeasonal = "DEVSEASONAL"

type RraDevSeasonal struct {
	RraSeasonal
}

func (r *RraDevSeasonal) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpString("cf", RraTypeDevSeasonal)
	r.dumpParams(rrdStore, dumper)
}

func newRraDevSeasonal(index int, store Store) (*RraDevSeasonal, error) {
	result := &RraDevSeasonal{
		RraSeasonal{
			RraAbstract: RraAbstract{
				Index: index,
			},
		},
	}

	if err := store.ReadRraParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
