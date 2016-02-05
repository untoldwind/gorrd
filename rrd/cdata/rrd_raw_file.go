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

	headerSize uint64

	pdpPreps []*RrdPdpPrep
	cdpPreps []*RrdCdpPrep

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
	if err := rrdFile.readRawHeader(); err != nil {
		return nil, err
	}
	datasources, err := rrdFile.readDatasources()
	if err != nil {
		return nil, err
	}
	rras, err := rrdFile.readRras()
	if err != nil {
		return nil, err
	}
	if err := rrdFile.read(datasources); err != nil {
		dataFile.Close()
		return nil, err
	}

	rrdFile.headerSize = dataFile.CurPosition()

	rrdFile.calculateRraStarts(rras)

	return &rrd.Rrd{
		Store:       rrdFile,
		Step:        rrdFile.header.pdpStep,
		LastUpdate:  rrdFile.lastUpdate,
		Datasources: datasources,
		Rras:        rras,
	}, nil
}

func (f *RrdRawFile) Close() {
	f.dataFile.Close()
}

func (f *RrdRawFile) read(datasources []rrd.RrdDatasource) error {
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
	for i, datasource := range datasources {
		datasource.SetLastValue(f.pdpPreps[i].lastDatasourceValue)
	}

	return nil
}

func (f *RrdRawFile) calculateRraStarts(rras []rrd.Rra) {
	f.rraStarts = make([]uint64, f.header.rraCount)
	rraNextStart := f.headerSize
	for i, rra := range rras {
		f.rraStarts[i] = rraNextStart
		rraNextStart += f.header.datasourceCount * rra.GetRowCount() * f.dataFile.ValueSize()
	}
}

func (f *RrdRawFile) RowIterator(rra rrd.Rra) (rrd.RraRowIterator, error) {
	stepPerRow := time.Duration(f.header.pdpStep*rra.GetPdpPerRow()) * time.Second
	startTime := f.lastUpdate.Truncate(stepPerRow).Add(-time.Duration(rra.GetRowCount()-1) * stepPerRow)
	iterator := &rrdRawRowIterator{
		dataFile:   f.dataFile,
		row:        0,
		rowCount:   rra.GetRowCount(),
		colCount:   f.header.datasourceCount,
		rraStart:   f.rraStarts[rra.GetIndex()],
		rraPtr:     f.rraPtrs[rra.GetIndex()],
		startTime:  startTime,
		stepPerRow: stepPerRow,
	}
	return iterator, nil
}
