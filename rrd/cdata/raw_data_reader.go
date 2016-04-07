package cdata

import (
	"bytes"

	"github.com/go-errors/errors"
)

type RawDataReader struct {
	*RawDataFile
	startPosition uint64
	position      uint64
}

func (f *RawDataReader) ReadBytes(len int) ([]byte, error) {
	data := make([]byte, len)
	if count, err := f.file.ReadAt(data, int64(f.position)); err != nil {
		return nil, err
	} else if count != len {
		return nil, errors.Errorf("Expected %d bytes (only %d read)", len, count)
	}
	f.position += uint64(len)
	return data, nil
}

func (f *RawDataReader) ReadCString(maxLen int) (string, error) {
	data, err := f.ReadBytes(maxLen)
	if err != nil {
		return "", nil
	}
	if idx := bytes.IndexByte(data, 0); idx >= 0 {
		return string(data[:idx]), nil
	}
	return "", errors.Errorf("Expected null terminated string")
}

func (f *RawDataReader) ReadUnival() (unival, error) {
	f.alignOffset()
	data, err := f.ReadBytes(f.valueSize)
	if err != nil {
		return 0, errors.Wrap(err, 0)
	}
	return f.bytesToUnival(data), nil
}

func (f *RawDataReader) ReadDouble() (float64, error) {
	unival, err := f.ReadUnival()
	if err != nil {
		return 0, err
	}
	return unival.AsDouble(), nil
}

func (f *RawDataReader) ReadDoubles(buffer []float64) error {
	f.alignOffset()
	data, err := f.ReadBytes(f.valueSize * len(buffer))
	if err != nil {
		return err
	}
	offset := 0
	for i := range buffer {
		buffer[i] = f.bytesToUnival(data[offset:]).AsDouble()
		offset += f.valueSize
	}

	return nil
}

func (f *RawDataReader) ReadUnsignedLong() (uint64, error) {
	unival, err := f.ReadUnival()
	if err != nil {
		return 0, err
	}
	return unival.AsUnsignedLong(), nil
}

func (f *RawDataReader) ReadUnivals(count int) ([]unival, error) {
	f.alignOffset()
	data, err := f.ReadBytes(f.valueSize * count)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	return f.BytesToUnivals(data), nil
}

func (f *RawDataReader) Seek(offset uint64) {
	f.position = f.startPosition + offset
}

func (f *RawDataReader) CurPosition() uint64 {
	return f.position
}

func (f *RawDataReader) alignOffset() {
	mod := f.position % f.byteAlignment
	if mod == 0 {
		return
	}
	f.position += f.byteAlignment - mod
}
