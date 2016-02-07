package rrd

const DatasourceTypeCounter = "COUNTER"

type DatasourceCounter struct {
	DatasourceAbstract
}

func (d *DatasourceCounter) DumpTo(dumper DataDumper) error {
	if err := dumper.DumpString("type", DatasourceTypeCounter); err != nil {
		return err
	}
	return d.DatasourceAbstract.DumpTo(dumper)
}

func newDatasourceCounter(index int, store Store) (*DatasourceCounter, error) {
	result := &DatasourceCounter{}

	if err := store.ReadDatasourceParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
