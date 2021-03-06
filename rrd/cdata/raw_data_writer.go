package cdata

import "github.com/go-errors/errors"

type RawDataWriter struct {
	*RawDataFile
	startPosition uint64
	position      uint64
}

func (f *RawDataWriter) WriteBytes(data []byte) error {
	if count, err := f.file.WriteAt(data, int64(f.position)); err != nil {
		return errors.Wrap(err, 0)
	} else if count != len(data) {
		return errors.Errorf("Expected %d bytes (only %d read)", len(data), count)
	}
	f.position += uint64(len(data))
	return nil
}

func (f *RawDataWriter) WriteCString(str string, maxLen int) error {
	if len([]byte(str)) >= maxLen-1 {
		return errors.Errorf("String too long len(%s) >= %d", str, maxLen+1)
	}
	data := make([]byte, maxLen)
	copy(data, []byte(str))

	return f.WriteBytes(data)
}

func (f *RawDataWriter) WriteUnival(val unival) error {
	return f.WriteUnsignedLong(val.AsUnsignedLong())
}

func (f *RawDataWriter) WriteUnivals(univals []unival) error {
	f.alignOffset()
	data := f.UnivalsToBytes(univals)
	return f.WriteBytes(data)
}

func (f *RawDataWriter) WriteDouble(val float64) error {
	return f.WriteUnival(univalForDouble(val))
}

func (f *RawDataWriter) WriteDoubles(vals []float64) error {
	f.alignOffset()
	data := make([]byte, f.valueSize*len(vals))

	offset := 0
	for _, val := range vals {
		f.univalToBytes(data[offset:], univalForDouble(val))
		offset += f.valueSize
	}
	return f.WriteBytes(data)
}

func (f *RawDataWriter) WriteUnsignedLong(val uint64) error {
	f.alignOffset()
	data := make([]byte, f.valueSize)
	f.univalToBytes(data, univalForUnsignedLong(val))

	return f.WriteBytes(data)
}

func (f *RawDataWriter) Seek(offset uint64) {
	f.position = f.startPosition + offset
}

func (f *RawDataWriter) CurPosition() uint64 {
	return f.position
}

func (f *RawDataWriter) alignOffset() {
	mod := f.position % f.byteAlignment
	if mod == 0 {
		return
	}
	f.position += f.byteAlignment - mod
}
