package cdata

type RrdCdpPrep struct {
	scratch []unival
}

func (f *RrdRawFile) readCdpPreps(header *rrdRawHeader) error {
	f.cdpPreps = make([]*RrdCdpPrep, header.rraCount*header.datasourceCount)

	var err error
	for i := range f.cdpPreps {
		f.cdpPreps[i], err = f.readCdpPrep()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *RrdRawFile) readCdpPrep() (*RrdCdpPrep, error) {
	scratch, err := f.dataFile.ReadUnivals(10)
	if err != nil {
		return nil, err
	}
	return &RrdCdpPrep{
		scratch: scratch,
	}, nil
}
