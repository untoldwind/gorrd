package rrd

const RraTypeAverage = "AVERAGE"

type RraAverage struct {
	RraAbstractGeneric
}

func (r *RraAverage) DumpTo(rrdStore RrdStore, dumper RrdDumper) error {
	if err := dumper.DumpString("cf", RraTypeAverage); err != nil {
		return err
	}
	return r.RraAbstractGeneric.DumpTo(rrdStore, dumper)
}
