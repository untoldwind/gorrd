package cdata

import (
	"encoding/binary"
	"time"

	"github.com/go-errors/errors"
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

	rraPtrsChanged bool
	rraPtrs        []uint64
	rraStarts      []uint64
}

func OpenRrdRawFile(name string, readOnly bool) (*rrd.Rrd, error) {
	dataFile, err := OpenCDataFile(name, readOnly, binary.LittleEndian, 8, 8)
	if err != nil {
		return nil, err
	}

	rrdFile := &RrdRawFile{
		dataFile: dataFile,
	}
	reader := dataFile.Reader(0)
	if err := rrdFile.readHeaders(reader); err != nil {
		dataFile.Close()
		return nil, errors.Wrap(err, 0)
	}

	rrdFile.headerSize = reader.CurPosition()

	rrdFile.calculateRraStarts()

	rrd, err := rrd.NewRrd(rrdFile)
	if err != nil {
		return nil, errors.Wrap(err, 0)
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

func (f *RrdRawFile) readHeaders(reader *CDataReader) error {
	if err := f.readVersionHeader(reader); err != nil {
		return errors.Wrap(err, 0)
	}
	if err := f.readDatasources(reader); err != nil {
		return errors.Wrap(err, 0)
	}
	if err := f.readRras(reader); err != nil {
		return errors.Wrap(err, 0)
	}
	f.baseHeaderSize = reader.CurPosition()
	if err := f.readLiveHead(reader); err != nil {
		return errors.Wrap(err, 0)
	}
	if err := f.readPdpPreps(reader); err != nil {
		return errors.Wrap(err, 0)
	}
	if err := f.readCdpPreps(reader); err != nil {
		return errors.Wrap(err, 0)
	}
	if err := f.readRraPtrs(reader); err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (f *RrdRawFile) calculateRraStarts() {
	f.rraStarts = make([]uint64, f.header.rraCount)
	rraNextStart := f.headerSize
	for i, rraDef := range f.rraDefs {
		f.rraStarts[i] = rraNextStart
		rraNextStart += f.header.datasourceCount * rraDef.rowCount * f.dataFile.ValueSize()
	}
}
