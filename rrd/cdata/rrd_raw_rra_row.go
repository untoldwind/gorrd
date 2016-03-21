package cdata

func (f *RrdRawFile) StoreRow(rraIndex int, row []float64) error {
	writer := f.dataFile.Writer(f.rraStarts[rraIndex] +
		f.rraPtrs[rraIndex]*f.header.datasourceCount*f.dataFile.ValueSize())

	for _, col := range row {
		if err := writer.WriteDouble(col); err != nil {
			return err
		}
	}

	f.rraPtrsChanged = true
	f.rraPtrs[rraIndex]++
	if f.rraPtrs[rraIndex] >= f.rraDefs[rraIndex].rowCount {
		f.rraPtrs[rraIndex] = 0
	}

	return nil
}
