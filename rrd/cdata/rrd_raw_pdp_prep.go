package cdata

type rrdPdpPrep struct {
	lastDatasourceValue string
	scratch             []unival
}

func (f *RrdRawFile) readPdpPreps() error {
	f.pdpPreps = make([]*rrdPdpPrep, f.header.datasourceCount)

	var err error
	for i := range f.pdpPreps {
		f.pdpPreps[i], err = f.readPdpPrep()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *RrdRawFile) readPdpPrep() (*rrdPdpPrep, error) {
	value, err := f.dataFile.ReadCString(30)
	if err != nil {
		return nil, err
	}
	scratch, err := f.dataFile.ReadUnivals(10)
	if err != nil {
		return nil, err
	}
	return &rrdPdpPrep{
		lastDatasourceValue: value,
		scratch:             scratch,
	}, nil
}
