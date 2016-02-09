package cdata

type rrdCdpPrep struct {
	scratch []unival
}

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
