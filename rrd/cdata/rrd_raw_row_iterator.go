package cdata

import (
	"time"

	"github.com/untoldwind/gorrd/rrd"
)

type rrdRawRowIterator struct {
	dataFile   *CDataFile
	row        uint64
	colCount   uint64
	rowCount   uint64
	rraStart   uint64
	rraPtr     uint64
	startTime  time.Time
	stepPerRow time.Duration
	lastRow    *rrd.RraRow
	lastError  error
}

func (i *rrdRawRowIterator) seekStart() {
	i.lastError = i.dataFile.Seek(i.rraStart + (i.rraPtr+1)*i.colCount*i.dataFile.ValueSize())
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
			i.lastError = i.dataFile.Seek(i.rraStart)
		}
	}

	if i.lastError == nil {
		i.lastRow.Timestamp = i.startTime.Add(time.Duration(i.row) * i.stepPerRow)
		i.lastError = i.dataFile.ReadDoubles(i.lastRow.Values)
	}

	return i.row >= 0 && i.row < i.rowCount
}

func (i *rrdRawRowIterator) Value() (*rrd.RraRow, error) {
	if i.lastError != nil {
		return nil, i.lastError
	}
	return i.lastRow, nil
}