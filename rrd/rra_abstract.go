package rrd

type RraCpdPrepBase struct {
	PrimaryValue   float64 `cdp:"8"`
	SecondaryValue float64 `cdp:"9"`
}

func (c *RraCpdPrepBase) DumpTo(dumper DataOutput) {
	dumper.DumpSubFields("ds", func(ds DataOutput) error {
		ds.DumpDouble("primary_value", c.PrimaryValue)
		ds.DumpDouble("secondary_value", c.SecondaryValue)
		return nil
	})
}

type RraAbstract struct {
	Index     int
	RowCount  uint64 `rra:"rowCount"`
	PdpPerRow uint64 `rra:"pdpPerRow"`
}
