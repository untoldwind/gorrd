package rrd

type Rra interface {
	GetIndex() int
	GetRowCount() uint64
	GetPdpPerRow() uint64
	DumpTo(rrdStore RrdStore, dumper RrdDumper) error
}

type RraAbstractGeneric struct {
	Index        int
	RowCount     uint64
	PdpPerRow    uint64
	XFilesFactor float64
}

func (r *RraAbstractGeneric) GetIndex() int {
	return r.Index
}

func (r *RraAbstractGeneric) GetRowCount() uint64 {
	return r.RowCount
}

func (r *RraAbstractGeneric) GetPdpPerRow() uint64 {
	return r.PdpPerRow
}

func (r *RraAbstractGeneric) DumpTo(rrdStore RrdStore, dumper RrdDumper) error {
	if err := dumper.DumpUnsignedLong("pdp_per_row", r.PdpPerRow); err != nil {
		return err
	}
	if err := dumper.DumpSubFields("params", func(params RrdDumper) error {
		return params.DumpDouble("xff", r.XFilesFactor)
	}); err != nil {
		return err
	}
	return dumper.DumpSubFields("database", func(database RrdDumper) error {
		rowIterator, err := rrdStore.RowIterator(r)
		if err != nil {
			return err
		}
		return ForEachRow(rowIterator, func(row *RraRow) error {
			return row.DumpTo(dumper)
		})
	})
}
