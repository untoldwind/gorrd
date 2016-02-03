package cdata

import (
	"github.com/go-errors/errors"
	"github.com/untoldwind/gorrd/rrd"
)

func readDatasources(header *rrdRawHeader, dataFile *CDataFile) ([]rrd.RrdDatasource, error) {
	result := make([]rrd.RrdDatasource, header.datasourceCount)

	var err error
	for i := range result {
		result[i], err = readDatasource(dataFile)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func readDatasource(dataFile *CDataFile) (rrd.RrdDatasource, error) {
	name, err := dataFile.ReadCString(20)
	if err != nil {
		return nil, err
	}
	datasourceType, err := dataFile.ReadCString(20)
	if err != nil {
		return nil, err
	}
	parameters, err := dataFile.ReadUnivals(10)
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
			rrd.RrdDatasourceAbstract{
				Name:      name,
				Heartbeat: parameters[0].AsUnsignedLong(),
				Min:       parameters[1].AsDouble(),
				Max:       parameters[2].AsDouble(),
			},
		}, nil
	case rrd.RrdDatasourceTypeDCounter:
		return &rrd.RrdDCounterDatasource{
			rrd.RrdDatasourceAbstract{
				Name:      name,
				Heartbeat: parameters[0].AsUnsignedLong(),
				Min:       parameters[1].AsDouble(),
				Max:       parameters[2].AsDouble(),
			},
		}, nil
	case rrd.RrdDatasourceTypeDDerive:
		return &rrd.RrdDDeriveDatasource{
			rrd.RrdDatasourceAbstract{
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
			rrd.RrdDatasourceAbstract{
				Name:      name,
				Heartbeat: parameters[0].AsUnsignedLong(),
				Min:       parameters[1].AsDouble(),
				Max:       parameters[2].AsDouble(),
			},
		}, nil
	}

	return nil, errors.Errorf("Unknown datasource type: %s", datasourceType)
}
