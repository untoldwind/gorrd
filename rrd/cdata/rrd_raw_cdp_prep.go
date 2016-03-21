package cdata

type rrdCdpPrep struct {
	scratch []unival
}

const rrdRawCdpPrepSize = 10 * 8

func (f *RrdRawFile) readCdpPreps(reader *CDataReader) error {
	f.cdpPreps = make([][]*rrdCdpPrep, f.header.rraCount)

	var err error
	for i := range f.cdpPreps {
		f.cdpPreps[i] = make([]*rrdCdpPrep, f.header.datasourceCount)
		for j := range f.cdpPreps[i] {
			f.cdpPreps[i][j], err = f.readCdpPrep(reader)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *RrdRawFile) readCdpPrep(reader *CDataReader) (*rrdCdpPrep, error) {
	scratch, err := reader.ReadUnivals(10)
	if err != nil {
		return nil, err
	}
	return &rrdCdpPrep{
		scratch: scratch,
	}, nil
}

func (f *RrdRawFile) storeCdpPreps() error {
	writer := f.dataFile.Writer(f.baseHeaderSize + rrdRawLiveHeaderSize +
		rrdRawPdpPrepSize*f.header.datasourceCount)

	for _, cdpPreps := range f.cdpPreps {
		for _, cdpPrep := range cdpPreps {
			if err := storeCdpPrep(writer, cdpPrep); err != nil {
				return err
			}
		}
	}
	return nil
}

func storeCdpPrep(writer *CDataWriter, cdpPrep *rrdCdpPrep) error {
	if err := writer.WriteUnivals(cdpPrep.scratch); err != nil {
		return err
	}
	return nil
}
