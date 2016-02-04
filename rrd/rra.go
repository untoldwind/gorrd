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
	if database, err := dumper.DumpSubFields("database"); err != nil {
		return err
	} else {
		rowIterator, err := rrdStore.RowIterator(r)
		if err != nil {
			return err
		}
		if err := ForEachRow(rowIterator, func(row *RraRow) error {
			return row.DumpTo(dumper)
		}); err != nil {
			return err
		}

		if err := database.Finalize(); err != nil {
			return err
		}
	}
	return dumper.Finalize()
}
