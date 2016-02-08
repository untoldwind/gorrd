package rrd

type Rra interface {
	GetRowCount() uint64
	GetPdpPerRow() uint64
	DumpTo(rrdStore Store, dumper DataOutput) error
}

type RraCpdPrepGeneric struct {
	PrimaryValue      float64 `cdp:"8"`
	SecondaryValue    float64 `cdp:"9"`
	Value             float64 `cdp:"0"`
	UnknownDatapoints uint64  `cdp:"1"`
}

func (c *RraCpdPrepGeneric) DumpTo(dumper DataOutput) error {
	return dumper.DumpSubFields("ds", func(ds DataOutput) error {
		if err := ds.DumpDouble("primary_value", c.PrimaryValue); err != nil {
			return err
		}
		if err := ds.DumpDouble("secondary_value", c.SecondaryValue); err != nil {
			return err
		}
		if err := ds.DumpDouble("value", c.Value); err != nil {
			return err
		}
		if err := ds.DumpUnsignedLong("unknown_datapoints", c.UnknownDatapoints); err != nil {
			return err
		}
		return nil
	})
}

type RraAbstractGeneric struct {
	Index        int
	RowCount     uint64              `rra:"rowCount"`
	PdpPerRow    uint64              `rra:"pdpPerRow"`
	XFilesFactor float64             `rra:"param0"`
	CpdPreps     []RraCpdPrepGeneric `rra:"cpdPreps"`
}

func (r *RraAbstractGeneric) GetRowCount() uint64 {
	return r.RowCount
}

func (r *RraAbstractGeneric) GetPdpPerRow() uint64 {
	return r.PdpPerRow
}

func (r *RraAbstractGeneric) DumpTo(rrdStore Store, dumper DataOutput) error {
	if err := dumper.DumpUnsignedLong("pdp_per_row", r.PdpPerRow); err != nil {
		return err
	}
	if err := dumper.DumpSubFields("params", func(params DataOutput) error {
		return params.DumpDouble("xff", r.XFilesFactor)
	}); err != nil {
		return err
	}
	if err := dumper.DumpSubFields("cdp_prep", func(cdpPreps DataOutput) error {
		for _, cdpPrep := range r.CpdPreps {
			if err := cdpPrep.DumpTo(cdpPreps); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return dumper.DumpSubFields("database", func(database DataOutput) error {
		rowIterator, err := rrdStore.RowIterator(r.Index)
		if err != nil {
			return err
		}
		return ForEachRow(rowIterator, func(row *RraRow) error {
			return row.DumpTo(dumper)
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
