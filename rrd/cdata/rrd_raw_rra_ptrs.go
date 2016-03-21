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

func (f *RrdRawFile) StoreRraPtrs() error {
	if !f.rraPtrsChanged {
		return nil
	}
	writer := f.dataFile.Writer(f.baseHeaderSize + rrdRawLiveHeaderSize +
		rrdRawPdpPrepSize*f.header.datasourceCount +
		rrdRawCdpPrepSize*f.header.datasourceCount*f.header.rraCount)

	for _, rraPtr := range f.rraPtrs {
		if err := writer.WriteUnsignedLong(rraPtr); err != nil {
			return err
		}
	}
	return nil
}
