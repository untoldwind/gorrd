package cdata

type rrdPdpPrep struct {
	lastDatasourceValue string
	scratch             []unival
}

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
