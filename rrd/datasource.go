package rrd

import (
	"math"
	"time"

	"github.com/go-errors/errors"
)

type Datasource interface {
	GetName() string
	GetLastValue() string
	CalculatePdpPrep(newValue string, interval float64) (float64, error)
	UpdatePdp(pdpValue, interval float64)
	ProcessPdp(pdpValue, interval, preInt float64, elapsedSteps uint64, step time.Duration) float64
	DumpTo(dumper DataOutput) error
}

type DatasourceAbstract struct {
	Name            string  `ds:"name"`
	Heartbeat       uint64  `ds:"param0"`
	Min             float64 `ds:"param1"`
	Max             float64 `ds:"param2"`
	LastValue       string  `pdp:"lastValue"`
	UnknownSecCount uint64  `pdp:"0"`
	PdpValue        float64 `pdp:"1"`
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

func (d *DatasourceAbstract) UpdatePdp(pdpValue, interval float64) {
	if math.IsNaN(pdpValue) {
		d.UnknownSecCount += uint64(interval)
	} else if math.IsNaN(d.PdpValue) {
		d.PdpValue = pdpValue
	} else {
		d.PdpValue += pdpValue
	}
}

func (d *DatasourceAbstract) ProcessPdp(pdpValue, interval, preInt float64, elapsedSteps uint64, step time.Duration) float64 {
	var preUnknown float64
	if math.IsNaN(pdpValue) {
		preUnknown = preInt
	} else {
		if math.IsNaN(d.PdpValue) {
			d.PdpValue = 0
		}
		d.PdpValue += pdpValue / interval * preInt
	}
	var pdpTemp float64

	if interval > float64(d.Heartbeat) {
		pdpTemp = math.NaN()
	} else {
		diffPdpSteps := (elapsedSteps * uint64(step)) / uint64(time.Second)
		pdpTemp = d.PdpValue / (float64(diffPdpSteps-d.UnknownSecCount) - preUnknown)
	}
	return pdpTemp
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
	if err := dumper.DumpDouble("value", d.PdpValue); err != nil {
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
