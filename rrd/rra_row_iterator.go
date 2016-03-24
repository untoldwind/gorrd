package rrd

import (
	"fmt"
	"time"
)

type RraRow struct {
	Timestamp time.Time
	Values    []float64
}

func (r *RraRow) DumpTo(dumper DataOutput) {
	dumper.DumpSubFields("row", func(row DataOutput) error {
		dumper.DumpComment(fmt.Sprintf("%s / %d", r.Timestamp.String(), r.Timestamp.Unix()))
		for _, value := range r.Values {
			dumper.DumpDouble("v", value)
		}
		return nil
	})
}

type RraRowIterator interface {
	Next() bool
	Value() (*RraRow, error)
	Close()
}

func ForEachRow(iterator RraRowIterator, collector func(row *RraRow) error) error {
	for iterator.Next() {
		row, err := iterator.Value()
		if err != nil {
			return err
		}
		if err := collector(row); err != nil {
			return err
		}
	}
	return nil
}
