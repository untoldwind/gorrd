package cdata

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"

	"github.com/go-errors/errors"
)

// CDataFile Helper to access files created from C code by directly mapping structs
// Honours byte order as well as byte alignment
type CDataFile struct {
	file          *os.File
	byteOrder     binary.ByteOrder
	byteAlignment uint64
	valueSize     uint64
}

type CDataReader struct {
	*CDataFile
	startPosition uint64
	position      uint64
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

func (f *CDataReader) ReadBytes(len int) ([]byte, error) {
	data := make([]byte, len)
	if count, err := f.file.ReadAt(data, int64(f.position)); err != nil {
		return nil, err
	} else if count != len {
		return nil, errors.Errorf("Expected %d bytes (only %d read)", len, count)
	}
	f.position += uint64(len)
	return data, nil
}

func (f *CDataReader) ReadCString(maxLen int) (string, error) {
	data, err := f.ReadBytes(maxLen)
	if err != nil {
		return "", nil
	}
	if idx := bytes.IndexByte(data, 0); idx >= 0 {
		return string(data[:idx]), nil
	}
	return "", errors.Errorf("Expected null terminated string")
}

func (f *CDataReader) ReadUnival() (unival, error) {
	f.alignOffset()
	data, err := f.ReadBytes(8)
	if err != nil {
		return 0, nil
	}
	return unival(f.byteOrder.Uint64(data)), nil
}

func (f *CDataReader) ReadDouble() (float64, error) {
	unival, err := f.ReadUnival()
	if err != nil {
		return 0, err
	}
	return unival.AsDouble(), nil
}

func (f *CDataReader) ReadDoubles(buffer []float64) error {
	f.alignOffset()
	data, err := f.ReadBytes(8 * len(buffer))
	if err != nil {
		return err
	}
	for i := range buffer {
		buffer[i] = math.Float64frombits(f.byteOrder.Uint64(data[i*8 : (i+1)*8]))
	}

	return nil
}

func (f *CDataReader) ReadUnsignedLong() (uint64, error) {
	unival, err := f.ReadUnival()
	if err != nil {
		return 0, err
	}
	return unival.AsUnsignedLong(), nil
}

func (f *CDataReader) ReadUnivals(count int) ([]unival, error) {
	f.alignOffset()
	data, err := f.ReadBytes(8 * count)
	if err != nil {
		return nil, nil
	}
	result := make([]unival, count)
	for i := range result {
		result[i] = unival(f.byteOrder.Uint64(data[i*8 : (i+1)*8]))
	}
	return result, nil
}

func (f *CDataReader) Seek(offset uint64) {
	f.position = f.startPosition + offset
}

func (f *CDataReader) CurPosition() uint64 {
	return f.position
}

func (f *CDataReader) alignOffset() {
	mod := f.position % f.byteAlignment
	if mod == 0 {
		return
	}
	f.position += f.byteAlignment - mod
}
