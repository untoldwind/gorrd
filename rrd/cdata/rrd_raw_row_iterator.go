package cdata

import (
	"time"

	"github.com/untoldwind/gorrd/rrd"
)

func (f *RrdRawFile) RowIterator(rraIndex int) (rrd.RraRowIterator, error) {
	stepPerRow := time.Duration(f.header.pdpStep*f.rraDefs[rraIndex].pdpPerRow) * time.Second
	startTime := f.lastUpdate.Truncate(stepPerRow).Add(-time.Duration(f.rraDefs[rraIndex].rowCount-1) * stepPerRow)
	iterator := &rrdRawRowIterator{
		reader:     f.dataFile.Reader(f.rraStarts[rraIndex]),
		row:        0,
		rowCount:   f.rraDefs[rraIndex].rowCount,
		colCount:   f.header.datasourceCount,
		rraPtr:     f.rraPtrs[rraIndex],
		startTime:  startTime,
		stepPerRow: stepPerRow,
	}
	return iterator, nil
}

type rrdRawRowIterator struct {
	reader     *RawDataReader
	row        uint64
	colCount   uint64
	rowCount   uint64
	rraPtr     uint64
	startTime  time.Time
	stepPerRow time.Duration
	lastRow    *rrd.RraRow
	lastError  error
}

func (i *rrdRawRowIterator) seekStart() {
	if i.rraPtr+1 < i.rowCount {
		i.reader.Seek((i.rraPtr + 1) * i.colCount * i.reader.ValueSize())
	}
}

func (i *rrdRawRowIterator) Next() bool {
	if i.lastError == nil && i.lastRow == nil {
		i.seekStart()
		i.lastRow = &rrd.RraRow{
			Values: make([]float64, i.colCount),
		}
	} else {
		i.row++
		if i.row+i.rraPtr+1 == i.rowCount {
			i.reader.Seek(0)
		}
	}

	if i.lastError == nil {
		i.lastRow.Timestamp = i.startTime.Add(time.Duration(i.row) * i.stepPerRow)
		i.lastError = i.reader.ReadDoubles(i.lastRow.Values)
	}

	return i.row >= 0 && i.row < i.rowCount
}

func (i *rrdRawRowIterator) Value() (*rrd.RraRow, error) {
	if i.lastError != nil {
		return nil, i.lastError
	}
	return i.lastRow, nil
}

func (i *rrdRawRowIterator) Close() {

}
