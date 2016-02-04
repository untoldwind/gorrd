package cdata

import (
	"github.com/go-errors/errors"
	"github.com/untoldwind/gorrd/rrd"
)

func (f *RrdRawFile) readRras() ([]rrd.Rra, error) {
	result := make([]rrd.Rra, f.header.rraCount)

	var err error
	for i := range result {
		result[i], err = readRra(f.dataFile)
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
	pdpPerRow, err := dataFile.ReadUnsignedLong()
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
			rrd.RraAbstractGeneric{
				RowCount:     rowCount,
				PdpPerRow:    pdpPerRow,
				XFilesFactor: parameters[0].AsDouble(),
			},
		}, nil
	}

	return nil, errors.Errorf("Unknown rra type: %s", rraType)
}
