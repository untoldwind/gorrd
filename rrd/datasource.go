package rrd

import (
	"time"

	"github.com/go-errors/errors"
)

const Undefined = "U"

type Datasource interface {
	GetName() string
	GetLastValue() string
	CalculatePdpPrep(newValue string, interval float64) (float64, error)
	UpdatePdp(pdpValue, interval float64)
	ProcessPdp(pdpValue float64, elapsed ElapsedPdpSteps, step time.Duration) float64
	DumpTo(dumper DataOutput) error
}

func newDatasource(index int, datasourceType string, store Store) (Datasource, error) {
	switch datasourceType {
	case DatasourceTypeAbsolute:
		return newDatasourceAbsolute(index, store)
	case DatasourceTypeCounter:
		return newDatasourceCounter(index, store)
	case DatasourceTypeDCounter:
		return newDatasourceDCounter(index, store)
	case DatasourceTypeDDerive:
		return newDatasourceDDerive(index, store)
	case DatasourceTypeDerive:
		return newDatasourceDerive(index, store)
	case DatasourceTypeGauge:
		return newDatasourceGauge(index, store)
	}
	return nil, errors.Errorf("Unknown datasource type: %s", datasourceType)
}
