package rrd

import "time"

type RraRow struct {
	Timestamp time.Time
	Values    []float64
}

func (r *RraRow) DumpTo(dumper RrdDumper) error {
	row, err := dumper.DumpSubFields("row")
	if err != nil {
		return err
	}
	return row.Finalize()
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
