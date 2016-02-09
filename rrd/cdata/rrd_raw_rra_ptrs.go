package cdata

func (f *RrdRawFile) readRraPtrs(reader *CDataReader) error {
	f.rraPtrs = make([]uint64, f.header.rraCount)
	var err error
	for i := range f.rraPtrs {
		f.rraPtrs[i], err = reader.ReadUnsignedLong()
		if err != nil {
			return err
		}
	}
	return nil
}
