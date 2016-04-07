package cdata

import (
	"encoding/binary"
	"os"

	"github.com/go-errors/errors"
)

// RawDataFile Helper to access files created from C code by directly mapping structs
// Honours byte order as well as byte alignment
type RawDataFile struct {
	file *os.File
	//	byteOrder     binary.ByteOrder
	byteAlignment uint64
	valueSize     int
	univalToBytes func([]byte, unival)
	bytesToUnival func([]byte) unival
}

// Open a CDataFile
func OpenRawDataFile(name string, readOnly bool, byteOrder binary.ByteOrder, byteAlignment uint64, valueSize int) (*RawDataFile, error) {
	flag := os.O_RDWR
	if readOnly {
		flag = os.O_RDONLY
	}

	file, err := os.OpenFile(name, flag, 0644)
	if err != nil {
		return nil, err
	}

	var univalToBytes func([]byte, unival)
	var bytesToUnival func([]byte) unival
	switch valueSize {
	case 8:
		univalToBytes = func(dst []byte, src unival) {
			byteOrder.PutUint64(dst, src.AsUnsignedLong())
		}
		bytesToUnival = func(src []byte) unival {
			return unival(byteOrder.Uint64(src))
		}
	default:
		return nil, errors.Errorf("Invalid value size %d", valueSize)
	}
	return &RawDataFile{
		file: file,
		//		byteOrder:     byteOrder,
		byteAlignment: byteAlignment,
		valueSize:     valueSize,
		univalToBytes: univalToBytes,
		bytesToUnival: bytesToUnival,
	}, nil
}

func (f *RawDataFile) ValueSize() uint64 {
	return uint64(f.valueSize)
}

// Close the CDataFile
func (f *RawDataFile) Close() error {
	return f.file.Close()
}

func (f *RawDataFile) Reader(startPosition uint64) *RawDataReader {
	return &RawDataReader{
		RawDataFile:   f,
		startPosition: startPosition,
		position:      startPosition,
	}
}

func (f *RawDataFile) Writer(startPosition uint64) *RawDataWriter {
	return &RawDataWriter{
		RawDataFile:   f,
		startPosition: startPosition,
		position:      startPosition,
	}
}

func (f *RawDataFile) UnivalsToBytes(univals []unival) []byte {
	data := make([]byte, f.valueSize*len(univals))

	offset := 0
	for _, val := range univals {
		f.univalToBytes(data[offset:], val)
		offset += f.valueSize
	}
	return data
}

func (f *RawDataFile) BytesToUnivals(data []byte) []unival {
	offset := 0
	result := make([]unival, len(data)/f.valueSize)
	for i := range result {
		result[i] = f.bytesToUnival(data[offset:])
		offset += f.valueSize
	}
	return result
}
