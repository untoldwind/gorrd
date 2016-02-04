package cdata

import (
	"time"

	"github.com/go-errors/errors"
	"github.com/untoldwind/gorrd/rrd"
)

type rrdRawRowIterator struct {
	dataFile   *CDataFile
	row        int64
	rowCount   uint64
	rraStart   uint64
	rraPtr     uint64
	lastUpdate time.Time
	pdpStep    int64
	pdpPerRow  int64
}

func (i *rrdRawRowIterator) Next() bool {
	i.row++
	return i.row >= 0 && uint64(i.row) < i.rowCount
}

func (i *rrdRawRowIterator) Value() (*rrd.RraRow, error) {
	if i.row < 0 {
		return nil, errors.Errorf("RowIterator not initinalized")
	} else if uint64(i.row) >= i.rowCount {
		return nil, errors.Errorf("RowIterator exhausted")
	}
	if i.row == 0 {
		if err := i.dataFile.Seek(i.rraStart + i.rraPtr); err != nil {
			return nil, err
		}
	} else if uint64(i.row)+i.rraPtr >= i.rowCount {
		if err := i.dataFile.Seek(i.rraStart); err != nil {
			return nil, err
		}
	}
	now := i.lastUpdate.Add(time.Duration(-i.row*i.pdpStep*i.pdpPerRow) * time.Second)

	return &rrd.RraRow{
		Timestamp: now,
	}, nil
}
