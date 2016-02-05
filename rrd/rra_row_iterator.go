package rrd

import (
	"fmt"
	"time"
)

type RraRow struct {
	Timestamp time.Time
	Values    []float64
}

func (r *RraRow) DumpTo(dumper RrdDumper) error {
	return dumper.DumpSubFields("row", func(row RrdDumper) error {
		if err := dumper.DumpComment(fmt.Sprintf("%s / %d", r.Timestamp.String(), r.Timestamp.Unix())); err != nil {
			return err
		}
		for _, value := range r.Values {
			if err := dumper.DumpDouble("v", value); err != nil {
				return err
			}
		}
		return nil
	})
}

type RraRowIterator interface {
	Next() bool
	Value() (*RraRow, error)
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
