package rrd

import "time"

type Rrd struct {
	Store       Store
	Step        uint64
	LastUpdate  time.Time
	Datasources []Datasource
	Rras        []Rra
}

func NewRrd(store Store) (*Rrd, error) {
	datasourceTypes := store.DatasourceTypes()
	datasources := make([]Datasource, len(datasourceTypes))
	var err error
	for i, datasourceType := range datasourceTypes {
		datasources[i], err = newDatasource(i, datasourceType, store)
		if err != nil {
			return nil, err
		}
	}

	rraTypes := store.RraTypes()
	rras := make([]Rra, len(rraTypes))
	for i, rraType := range rraTypes {
		rras[i], err = newRra(i, rraType, store)
		if err != nil {
			return nil, err
		}
	}

	return &Rrd{
		Store:       store,
		Step:        store.Step(),
		LastUpdate:  store.LastUpdate(),
		Datasources: datasources,
		Rras:        rras,
	}, nil
}

func (r *Rrd) Close() {
	r.Store.Close()
}

func (r *Rrd) Update(timestamp time.Time, values []string) error {
	return nil
}

func (r *Rrd) DumpTo(dumper DataDumper) error {
	if err := dumper.DumpString("version", "0003"); err != nil {
		return err
	}
	if err := dumper.DumpUnsignedLong("step", r.Step); err != nil {
		return err
	}
	if err := dumper.DumpComment("Seconds"); err != nil {
		return err
	}
	if err := dumper.DumpTime("lastupdate", r.LastUpdate); err != nil {
		return err
	}
	for _, datasource := range r.Datasources {
		if err := dumper.DumpSubFields("ds", func(sub DataDumper) error {
			return datasource.DumpTo(sub)
		}); err != nil {
			return err
		}
	}
	for _, rra := range r.Rras {
		if err := dumper.DumpSubFields("rra", func(sub DataDumper) error {
			return rra.DumpTo(r.Store, sub)
		}); err != nil {
			return err
		}
	}
	return dumper.Finalize()
}
