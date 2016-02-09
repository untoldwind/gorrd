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

	header     *rrdRawHeader
	lastUpdate time.Time

	datasourceDefs []*rrdRawDatasourceDef
	rraDefs        []*rrdRawRraDef

	baseHeaderSize uint64
	headerSize     uint64

	pdpPreps []*rrdPdpPrep
	cdpPreps [][]*rrdCdpPrep

	rraPtrs   []uint64
	rraStarts []uint64
}

func OpenRrdRawFile(name string, readOnly bool) (*rrd.Rrd, error) {
	dataFile, err := OpenCDataFile(name, readOnly, binary.LittleEndian, 8, 8)
	if err != nil {
		return nil, err
	}

	rrdFile := &RrdRawFile{
		dataFile: dataFile,
	}
	if err := rrdFile.readHeaders(); err != nil {
		dataFile.Close()
		return nil, err
	}

	rrdFile.headerSize = dataFile.CurPosition()

	rrdFile.calculateRraStarts()

	rrd, err := rrd.NewRrd(rrdFile)
	if err != nil {
		return nil, err
	}

	return rrd, nil
}

func (f *RrdRawFile) LastUpdate() time.Time {
	return f.lastUpdate
}

func (f *RrdRawFile) Step() time.Duration {
	return time.Duration(f.header.pdpStep) * time.Second
}

func (f *RrdRawFile) Close() {
	f.dataFile.Close()
}

func (f *RrdRawFile) readHeaders() error {
	if err := f.readVersionHeader(); err != nil {
		return err
	}
	if err := f.readDatasources(); err != nil {
		return err
	}
	if err := f.readRras(); err != nil {
		return err
	}
	f.baseHeaderSize = f.dataFile.CurPosition()
	if err := f.readLiveHead(); err != nil {
		return err
	}
	if err := f.readPdpPreps(); err != nil {
		return err
	}
	if err := f.readCdpPreps(); err != nil {
		return err
	}
	if err := f.readRraPtrs(); err != nil {
		return err
	}

	return nil
}

func (f *RrdRawFile) calculateRraStarts() {
	f.rraStarts = make([]uint64, f.header.rraCount)
	rraNextStart := f.headerSize
	for i, rraDef := range f.rraDefs {
		f.rraStarts[i] = rraNextStart
		rraNextStart += f.header.datasourceCount * rraDef.pdpPerRow * f.dataFile.ValueSize()
	}
}

func (f *RrdRawFile) RowIterator(rraIndex int) (rrd.RraRowIterator, error) {
	stepPerRow := time.Duration(f.header.pdpStep*f.rraDefs[rraIndex].pdpPerRow) * time.Second
	startTime := f.lastUpdate.Truncate(stepPerRow).Add(-time.Duration(f.rraDefs[rraIndex].rowCount-1) * stepPerRow)
	iterator := &rrdRawRowIterator{
		dataFile:   f.dataFile,
		row:        0,
		rowCount:   f.rraDefs[rraIndex].rowCount,
		colCount:   f.header.datasourceCount,
		rraStart:   f.rraStarts[rraIndex],
		rraPtr:     f.rraPtrs[rraIndex],
		startTime:  startTime,
		stepPerRow: stepPerRow,
	}
	return iterator, nil
}
