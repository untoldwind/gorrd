package cdata

import (
	"github.com/go-errors/errors"
	"github.com/untoldwind/gorrd/rrd"
)

func readRras(header *rrdRawHeader, dataFile *CDataFile) ([]rrd.Rra, error) {
	result := make([]rrd.Rra, header.rraCount)

	var err error
	for i := range result {
		result[i], err = readRra(dataFile)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func readRra(dataFile *CDataFile) (rrd.Rra, error) {
	rraType, err := dataFile.ReadCString(20)
	if err != nil {
		return nil, err
	}
	rowCount, err := dataFile.ReadUnsignedLong()
	if err != nil {
		return nil, err
	}
	pdpCount, err := dataFile.ReadUnsignedLong()
	if err != nil {
		return nil, err
	}
	parameters, err := dataFile.ReadUnivals(10)
	if err != nil {
		return nil, err
	}

	switch rraType {
	case rrd.RraTypeAverage:
		return &rrd.RraAverage{
			RowCount:     rowCount,
			PdpCount:     pdpCount,
			XFilesFactor: parameters[0].AsDouble(),
		}, nil
	}

	return nil, errors.Errorf("Unknown rra type: %s", rraType)
}
