package cdata

import (
	"encoding/binary"

	"fmt"
	"github.com/untoldwind/gorrd/rrd"
	"time"
)

const rrdCookie = "RRD"
const rrdFloatCookie = 8.642135E130

type RrdRawFile struct {
	dataFile *CDataFile

	datasourceCount uint64
	rraCount        uint64
	pdpStep         uint64

	datasources []rrd.RrdDatasource
	rras        []rrd.Rra

	lastUpdate time.Time

	pdpPreps []*RrdPdpPrep
	cdpPreps []*RrdCdpPrep
}

func OpenRrdRawFile(name string, readOnly bool) (*RrdRawFile, error) {
	dataFile, err := OpenDataFile(name, readOnly, binary.LittleEndian, 8)
	if err != nil {
		return nil, err
	}

	rrdFile := &RrdRawFile{
		dataFile: dataFile,
	}
	if err := rrdFile.read(); err != nil {
		dataFile.Close()
		return nil, err
	}

	fmt.Printf("%#v\n", rrdFile)
	for _, rra := range rrdFile.rras {
		fmt.Printf("%#v\n", rra)
	}
	fmt.Printf("%v\n", rrdFile.lastUpdate)
	for _, pdpPrep := range rrdFile.pdpPreps {
		fmt.Printf("%#v\n", pdpPrep)
	}
	for _, cdpPrep := range rrdFile.cdpPreps {
		fmt.Printf("%#v\n", cdpPrep)
	}
	return rrdFile, nil
}

func (f *RrdRawFile) Close() {
	f.dataFile.Close()
}

func (f *RrdRawFile) read() error {
	if err := f.readHeader(); err != nil {
		return err
	}
	if err := f.readDatasources(); err != nil {
		return err
	}
	if err := f.readRras(); err != nil {
		return err
	}
	if err := f.readLiveHead(); err != nil {
		return err
	}
	if err := f.readPdpPreps(); err != nil {
		return err
	}
	if err := f.readCdpPreps(); err != nil {
		return err
	}

	return nil
}
