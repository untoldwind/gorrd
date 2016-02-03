package cdata

import (
	"encoding/binary"

	"github.com/untoldwind/gorrd/rrd"
)

const rrdCookie = "RRD"
const rrdFloatCookie = 8.642135E130

type RrdRawFile struct {
	dataFile *CDataFile

	pdpPreps []*RrdPdpPrep
	cdpPreps []*RrdCdpPrep

	rraPtrs []RrdRraPtr
}

func OpenRrdRawFile(name string, readOnly bool) (*rrd.Rrd, error) {
	dataFile, err := OpenCDataFile(name, readOnly, binary.LittleEndian, 8)
	if err != nil {
		return nil, err
	}

	rrdFile := &RrdRawFile{
		dataFile: dataFile,
	}
	header, err := readRawHeader(dataFile)
	if err != nil {
		return nil, err
	}
	datasources, err := readDatasources(header, dataFile)
	if err != nil {
		return nil, err
	}
	rras, err := readRras(header, dataFile)
	if err != nil {
		return nil, err
	}
	lastUpdate, err := readLiveHead(dataFile)
	if err != nil {
		return nil, err
	}

	if err := rrdFile.read(header); err != nil {
		dataFile.Close()
		return nil, err
	}

	return &rrd.Rrd{
		Store:       rrdFile,
		Step:        header.pdpStep,
		LastUpdate:  lastUpdate,
		Datasources: datasources,
		Rras:        rras,
	}, nil
}

func (f *RrdRawFile) Close() {
	f.dataFile.Close()
}

func (f *RrdRawFile) read(header *rrdRawHeader) error {
	if err := f.readPdpPreps(header); err != nil {
		return err
	}
	if err := f.readCdpPreps(header); err != nil {
		return err
	}
	if err := f.readRraPtrs(header); err != nil {
		return err
	}

	return nil
}
