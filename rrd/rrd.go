package rrd

import "time"

type Rrd struct {
	Version     uint16
	Store       Store
	Step        time.Duration
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

	version := store.Version()
	rraTypes := store.RraTypes()
	rras := make([]Rra, len(rraTypes))
	for i, rraType := range rraTypes {
		rras[i], err = newRra(i, rraType, store)
		if err != nil {
			return nil, err
		}
	}

	return &Rrd{
		Version:     version,
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
