package cdata

type RrdPdpPrep struct {
	lastDatasourceValue string
	scratch             []unival
}

func (f *RrdRawFile) readPdpPreps(header *rrdRawHeader) error {
	f.pdpPreps = make([]*RrdPdpPrep, header.datasourceCount)

	var err error
	for i := range f.pdpPreps {
		f.pdpPreps[i], err = f.readPdpPrep()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *RrdRawFile) readPdpPrep() (*RrdPdpPrep, error) {
	value, err := f.dataFile.ReadCString(30)
	if err != nil {
		return nil, err
	}
	scratch, err := f.dataFile.ReadUnivals(10)
	if err != nil {
		return nil, err
	}
	return &RrdPdpPrep{
		lastDatasourceValue: value,
		scratch:             scratch,
	}, nil
}
