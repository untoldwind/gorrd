package cdata

import (
	"encoding/binary"
	"os"
)

// CDataFile Helper to access files created from C code by directly mapping structs
// Honours byte order as well as byte alignment
type CDataFile struct {
	file          *os.File
	byteOrder     binary.ByteOrder
	byteAlignment uint64
	valueSize     uint64
}

// Open a CDataFile
func OpenCDataFile(name string, readOnly bool, byteOrder binary.ByteOrder, byteAlignment, valueSize uint64) (*CDataFile, error) {
	flag := os.O_RDWR
	if readOnly {
		flag = os.O_RDONLY
	}

	file, err := os.OpenFile(name, flag, 0644)
	if err != nil {
		return nil, err
	}

	return &CDataFile{
		file:          file,
		byteOrder:     byteOrder,
		byteAlignment: byteAlignment,
		valueSize:     valueSize,
	}, nil
}

func (f *CDataFile) ValueSize() uint64 {
	return f.valueSize
}

// Close the CDataFile
func (f *CDataFile) Close() error {
	return f.file.Close()
}

func (f *CDataFile) Reader(startPosition uint64) *CDataReader {
	return &CDataReader{
		CDataFile:     f,
		startPosition: startPosition,
		position:      startPosition,
	}
}

func (f *CDataFile) Writer(startPosition uint64) *CDataWriter {
	return &CDataWriter{
		CDataFile:     f,
		startPosition: startPosition,
		position:      startPosition,
	}
}
