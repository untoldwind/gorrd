package cdata

type rrdCdpPrep struct {
	scratch []unival
}

func (f *RrdRawFile) readCdpPreps() error {
	f.cdpPreps = make([][]*rrdCdpPrep, f.header.rraCount)

	var err error
	for i := range f.cdpPreps {
		f.cdpPreps[i] = make([]*rrdCdpPrep, f.header.datasourceCount)
		for j := range f.cdpPreps[i] {
			f.cdpPreps[i][j], err = f.readCdpPrep()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *RrdRawFile) readCdpPrep() (*rrdCdpPrep, error) {
	scratch, err := f.dataFile.ReadUnivals(10)
	if err != nil {
		return nil, err
	}
	return &rrdCdpPrep{
		scratch: scratch,
	}, nil
}
