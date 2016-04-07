package cdata

import "github.com/go-errors/errors"

type CDataWriter struct {
	*CDataFile
	startPosition uint64
	position      uint64
}

func (f *CDataWriter) WriteBytes(data []byte) error {
	if count, err := f.file.WriteAt(data, int64(f.position)); err != nil {
		return errors.Wrap(err, 0)
	} else if count != len(data) {
		return errors.Errorf("Expected %d bytes (only %d read)", len(data), count)
	}
	f.position += uint64(len(data))
	return nil
}

func (f *CDataWriter) WriteCString(str string, maxLen int) error {
	if len([]byte(str)) >= maxLen-1 {
		return errors.Errorf("String too long len(%s) >= %d", str, maxLen+1)
	}
	data := make([]byte, maxLen)
	copy(data, []byte(str))

	return f.WriteBytes(data)
}

func (f *CDataWriter) WriteUnival(val unival) error {
	return f.WriteUnsignedLong(val.AsUnsignedLong())
}

func (f *CDataWriter) WriteUnivals(univals []unival) error {
	f.alignOffset()
	data := f.UnivalsToBytes(univals)
	return f.WriteBytes(data)
}

func (f *CDataWriter) WriteDouble(val float64) error {
	return f.WriteUnival(univalForDouble(val))
}

func (f *CDataWriter) WriteDoubles(vals []float64) error {
	f.alignOffset()
	data := make([]byte, f.valueSize*len(vals))

	offset := 0
	for _, val := range vals {
		f.univalToBytes(data[offset:], univalForDouble(val))
		offset += f.valueSize
	}
	return f.WriteBytes(data)
}

func (f *CDataWriter) WriteUnsignedLong(val uint64) error {
	f.alignOffset()
	data := make([]byte, f.valueSize)
	f.univalToBytes(data, univalForUnsignedLong(val))

	return f.WriteBytes(data)
}

func (f *CDataWriter) Seek(offset uint64) {
	f.position = f.startPosition + offset
}

func (f *CDataWriter) CurPosition() uint64 {
	return f.position
}

func (f *CDataWriter) alignOffset() {
	mod := f.position % f.byteAlignment
	if mod == 0 {
		return
	}
	f.position += f.byteAlignment - mod
}
