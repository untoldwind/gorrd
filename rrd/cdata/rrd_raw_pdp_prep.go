package cdata

type rrdPdpPrep struct {
	lastDatasourceValue string
	scratch             []unival
}

const rrdRawPdpPrepSize = 30 + 8*10

func (f *RrdRawFile) readPdpPreps(reader *CDataReader) error {
	f.pdpPreps = make([]*rrdPdpPrep, f.header.datasourceCount)

	var err error
	for i := range f.pdpPreps {
		f.pdpPreps[i], err = f.readPdpPrep(reader)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *RrdRawFile) readPdpPrep(reader *CDataReader) (*rrdPdpPrep, error) {
	value, err := reader.ReadCString(30)
	if err != nil {
		return nil, err
	}
	scratch, err := reader.ReadUnivals(10)
	if err != nil {
		return nil, err
	}
	return &rrdPdpPrep{
		lastDatasourceValue: value,
		scratch:             scratch,
	}, nil
}

func (f *RrdRawFile) storePdpPreps() error {
	writer := f.dataFile.Writer(f.baseHeaderSize + rrdRawLiveHeaderSize)

	for _, pdpPrep := range f.pdpPreps {
		if err := storePdpPrep(writer, pdpPrep); err != nil {
			return err
		}
	}
	return nil
}

func storePdpPrep(writer *CDataWriter, pdpPrep *rrdPdpPrep) error {
	if err := writer.WriteCString(pdpPrep.lastDatasourceValue, 30); err != nil {
		return err
	}
	if err := writer.WriteUnivals(pdpPrep.scratch); err != nil {
		return err
	}
	return nil
}
