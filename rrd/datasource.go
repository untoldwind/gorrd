package rrd

import (
	"math"

	"github.com/go-errors/errors"
)

type Datasource interface {
	GetName() string
	GetLastValue() string
	UpdatePdpPrep(newValue string, interval float64) (float64, error)
	DumpTo(dumper DataOutput) error
}

type DatasourceAbstract struct {
	Name            string  `ds:"name"`
	Heartbeat       uint64  `ds:"param0"`
	Min             float64 `ds:"param1"`
	Max             float64 `ds:"param2"`
	LastValue       string  `pdp:"lastValue"`
	UnknownSecCount uint64  `pdp:"0"`
	Value           float64 `pdp:"1"`
}

func (d *DatasourceAbstract) GetName() string {
	return d.Name
}

func (d *DatasourceAbstract) SetLastValue(lastValue string) {
	d.LastValue = lastValue
}

func (d *DatasourceAbstract) GetLastValue() string {
	return d.LastValue
}

func (d *DatasourceAbstract) DumpTo(dumper DataOutput) error {
	if err := dumper.DumpString("name", d.Name); err != nil {
		return err
	}
	if err := dumper.DumpUnsignedLong("minimal_heartbeat", d.Heartbeat); err != nil {
		return err
	}
	if err := dumper.DumpDouble("min", d.Min); err != nil {
		return err
	}
	if err := dumper.DumpDouble("max", d.Max); err != nil {
		return err
	}
	if err := dumper.DumpString("last_ds", d.LastValue); err != nil {
		return err
	}
	if err := dumper.DumpDouble("value", d.Value); err != nil {
		return err
	}
	if err := dumper.DumpUnsignedLong("unknown_sec", d.UnknownSecCount); err != nil {
		return err
	}
	return nil
}

func (d *DatasourceAbstract) checkRateBounds(rate float64) bool {
	return !math.IsNaN(rate) &&
		(math.IsNaN(d.Min) || rate >= d.Min) &&
		(math.IsNaN(d.Max) || rate <= d.Max)
}

func newDatasource(index int, datasourceType string, store Store) (Datasource, error) {
	switch datasourceType {
	case DatasourceTypeCounter:
		return newDatasourceCounter(index, store)
	case DatasourceTypeGauge:
		return newDatasourceGauge(index, store)
	}
	return nil, errors.Errorf("Unknown datasource type: %s", datasourceType)
}