package cdata

import (
	"encoding/binary"

	"time"

	"github.com/untoldwind/gorrd/rrd"
)

const rrdCookie = "RRD"
const rrdFloatCookie = 8.642135E130

type RrdRawFile struct {
	dataFile *CDataFile

	datasources []rrd.RrdDatasource
	rras        []rrd.Rra

	lastUpdate time.Time

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
	if err := rrdFile.read(header); err != nil {
		dataFile.Close()
		return nil, err
	}

	return &rrd.Rrd{
		Store:      rrdFile,
		Step:       header.pdpStep,
		LastUpdate: rrdFile.lastUpdate,
	}, nil
}

func (f *RrdRawFile) Close() {
	f.dataFile.Close()
}

func (f *RrdRawFile) read(header *rrdRawHeader) error {
	if err := f.readDatasources(header); err != nil {
		return err
	}
	if err := f.readRras(header); err != nil {
		return err
	}
	if err := f.readLiveHead(); err != nil {
		return err
	}
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
