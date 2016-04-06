package rrd

type RraCpdPrepBase struct {
	PrimaryValue   float64 `cdp:"8"`
	SecondaryValue float64 `cdp:"9"`
}

func (c *RraCpdPrepBase) DumpTo(dumper DataOutput) {
	dumper.DumpDouble("primary_value", c.PrimaryValue)
	dumper.DumpDouble("secondary_value", c.SecondaryValue)
}

type RraAbstract struct {
	Index     int
	RowCount  uint64 `rra:"rowCount"`
	PdpPerRow uint64 `rra:"pdpPerRow"`
}

func (r *RraAbstract) GetRowCount() uint64 {
	return r.RowCount
}

func (r *RraAbstract) GetPdpPerRow() uint64 {
	return r.PdpPerRow
}

func (r *RraAbstract) DumpDatabase(rrdStore Store, dumper DataOutput) {
	dumper.DumpSubFields("database", func(database DataOutput) error {
		rowIterator, err := rrdStore.RowIterator(r.Index)
		if err != nil {
			return err
		}
		return ForEachRow(rowIterator, func(row *RraRow) error {
			row.DumpTo(dumper)
			return nil
		})
	})

}
