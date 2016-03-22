package cdata

func (f *RrdRawFile) StoreRow(rraIndex int, row []float64) error {
	f.rraPtrsChanged = true
	f.rraPtrs[rraIndex]++
	if f.rraPtrs[rraIndex] >= f.rraDefs[rraIndex].rowCount {
		f.rraPtrs[rraIndex] = 0
	}

	writer := f.dataFile.Writer(f.rraStarts[rraIndex] +
		f.rraPtrs[rraIndex]*f.header.datasourceCount*f.dataFile.ValueSize())

	return writer.WriteDoubles(row)
}
