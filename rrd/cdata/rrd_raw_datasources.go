package cdata

import (
	"github.com/go-errors/errors"
	"github.com/untoldwind/gorrd/rrd"
)

func (f *RrdRawFile) readDatasources() error {
	f.datasources = make([]rrd.RrdDatasource, f.datasourceCount)

	var err error
	for i := range f.datasources {
		f.datasources[i], err = f.readDatasource()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *RrdRawFile) readDatasource() (rrd.RrdDatasource, error) {
	name, err := f.dataFile.ReadCString(20)
	if err != nil {
		return nil, err
	}
	datasourceType, err := f.dataFile.ReadCString(20)
	if err != nil {
		return nil, err
	}
	parameters, err := f.dataFile.ReadUnivals(10)
	if err != nil {
		return nil, err
	}

	switch datasourceType {
	case rrd.RrdDatasourceTypeAbsolute:
		return &rrd.RrdDatasourceAbsolute{
			rrd.RrdDatasourceAbstractLong{
				Name:      name,
				Heartbeat: parameters[0].AsUnsignedLong(),
				Min:       parameters[1].AsUnsignedLong(),
				Max:       parameters[2].AsUnsignedLong(),
			},
		}, nil
	case rrd.RrdDatasourceTypeCounter:
		return &rrd.RrdCounterDatasource{
			rrd.RrdDatasourceAbstractLong{
				Name:      name,
				Heartbeat: parameters[0].AsUnsignedLong(),
				Min:       parameters[1].AsUnsignedLong(),
				Max:       parameters[2].AsUnsignedLong(),
			},
		}, nil
	case rrd.RrdDatasourceTypeDCounter:
		return &rrd.RrdDCounterDatasource{
			rrd.RrdDatasourceAbstractDouble{
				Name:      name,
				Heartbeat: parameters[0].AsUnsignedLong(),
				Min:       parameters[1].AsDouble(),
				Max:       parameters[2].AsDouble(),
			},
		}, nil
	case rrd.RrdDatasourceTypeDDerive:
		return &rrd.RrdDDeriveDatasource{
			rrd.RrdDatasourceAbstractDouble{
				Name:      name,
				Heartbeat: parameters[0].AsUnsignedLong(),
				Min:       parameters[1].AsDouble(),
				Max:       parameters[2].AsDouble(),
			},
		}, nil
	case rrd.RrdDatasourceTypeDerive:
		return &rrd.RrdDeriveDatasource{
			rrd.RrdDatasourceAbstractLong{
				Name:      name,
				Heartbeat: parameters[0].AsUnsignedLong(),
				Min:       parameters[1].AsUnsignedLong(),
				Max:       parameters[2].AsUnsignedLong(),
			},
		}, nil
	case rrd.RrdDatasourceTypeGauge:
		return &rrd.RrdGaugeDatasource{
			rrd.RrdDatasourceAbstractLong{
				Name:      name,
				Heartbeat: parameters[0].AsUnsignedLong(),
				Min:       parameters[1].AsUnsignedLong(),
				Max:       parameters[2].AsUnsignedLong(),
			},
		}, nil
	}

	return nil, errors.Errorf("Unknown datasource type: %s", datasourceType)
}
