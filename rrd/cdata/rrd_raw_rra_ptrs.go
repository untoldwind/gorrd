package cdata

type RrdRraPtr uint64

func (f *RrdRawFile) readRraPtrs(header *rrdRawHeader) error {
	f.rraPtrs = make([]RrdRraPtr, header.rraCount)
	var err error
	for i := range f.rraPtrs {
		f.rraPtrs[i], err = f.readRraPtr()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *RrdRawFile) readRraPtr() (RrdRraPtr, error) {
	ptr, err := f.dataFile.ReadUnsignedLong()
	if err != nil {
		return 0, nil
	}
	return RrdRraPtr(ptr), nil
}
