package cdata

func (f *RrdRawFile) readRraPtrs() error {
	f.rraPtrs = make([]uint64, f.header.rraCount)
	var err error
	for i := range f.rraPtrs {
		f.rraPtrs[i], err = f.dataFile.ReadUnsignedLong()
		if err != nil {
			return err
		}
	}
	return nil
}
