package rrd

const RraTypeAverage = "AVERAGE"

type RraAverage struct {
	RowCount     uint64
	PdpCount     uint64
	XFilesFactor float64
}

func (r *RraAverage) DumpTo(dumper RrdDumper) error {
	if err := dumper.DumpString("cf", RraTypeAverage); err != nil {
		return err
	}
	if err := dumper.DumpUnsignedLong("pdp_per_row", r.PdpCount); err != nil {
		return err
	}
	if params, err := dumper.DumpSubFields("params"); err != nil {
		return err
	} else {
		if err := params.DumpDouble("xff", r.XFilesFactor); err != nil {
			return err
		}
		if err := params.Finalize(); err != nil {
			return err
		}
	}
	return dumper.Finalize()
}
