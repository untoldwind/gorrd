package cdata

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/go-errors/errors"
)

// CDataFile Helper to access files created from C code by directly mapping structs
// Honours byte order as well as byte alignment
type CDataFile struct {
	file          *os.File
	position      uint64
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
		position:      0,
		byteOrder:     byteOrder,
		byteAlignment: byteAlignment,
		valueSize:     valueSize,
	}, nil
}

// Close the CDataFile
func (f *CDataFile) Close() error {
	return f.file.Close()
}

func (f *CDataFile) ReadBytes(len int) ([]byte, error) {
	data := make([]byte, len)
	if count, err := f.file.Read(data); err != nil {
		return nil, err
	} else if count != len {
		return nil, errors.Errorf("Expected %d bytes (only %d read)", len, count)
	}
	f.position += uint64(len)
	return data, nil
}

func (f *CDataFile) ReadCString(maxLen int) (string, error) {
	data, err := f.ReadBytes(maxLen)
	if err != nil {
		return "", nil
	}
	if idx := bytes.IndexByte(data, 0); idx >= 0 {
		return string(data[:idx]), nil
	}
	return "", errors.Errorf("Expected null terminated string")
}

func (f *CDataFile) ReadUnival() (unival, error) {
	if err := f.alignOffset(); err != nil {
		return 0, err
	}
	data, err := f.ReadBytes(8)
	if err != nil {
		return 0, nil
	}
	return unival(f.byteOrder.Uint64(data)), nil
}

func (f *CDataFile) ReadDouble() (float64, error) {
	unival, err := f.ReadUnival()
	if err != nil {
		return 0, err
	}
	return unival.AsDouble(), nil
}

func (f *CDataFile) ReadUnsignedLong() (uint64, error) {
	unival, err := f.ReadUnival()
	if err != nil {
		return 0, err
	}
	return unival.AsUnsignedLong(), nil
}

func (f *CDataFile) ReadUnivals(count int) ([]unival, error) {
	if err := f.alignOffset(); err != nil {
		return nil, err
	}
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

func (f *CDataFile) Seek(offset uint64) error {
	_, err := f.file.Seek(int64(offset), 0)
	return err
}

func (f *CDataFile) ValueSize() uint64 {
	return f.valueSize
}

func (f *CDataFile) CurPosition() uint64 {
	return f.position
}

func (f *CDataFile) alignOffset() error {
	skip := f.byteAlignment - (f.position % f.byteAlignment)
	if skip >= f.byteAlignment {
		return nil
	}
	if _, err := f.file.Seek(int64(skip), 1); err != nil {
		return err
	}
	f.position += skip
	return nil
}
