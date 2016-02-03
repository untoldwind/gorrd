package cdata

import (
	"github.com/go-errors/errors"
	"github.com/untoldwind/gorrd/rrd"
)

func (f *RrdRawFile) readRras(header *rrdRawHeader) error {
	f.rras = make([]rrd.Rra, header.rraCount)

	var err error
	for i := range f.rras {
		f.rras[i], err = f.readRra()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *RrdRawFile) readRra() (rrd.Rra, error) {
	rraType, err := f.dataFile.ReadCString(20)
	if err != nil {
		return nil, err
	}
	rowCount, err := f.dataFile.ReadUnsignedLong()
	if err != nil {
		return nil, err
	}
	pdpCount, err := f.dataFile.ReadUnsignedLong()
	if err != nil {
		return nil, err
	}
	parameters, err := f.dataFile.ReadUnivals(10)
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
