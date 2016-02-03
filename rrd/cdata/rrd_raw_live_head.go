package cdata

import "time"

func (f *RrdRawFile) readLiveHead() error {
	timeSec, err := f.dataFile.ReadUnival()
	if err != nil {
		return err
	}
	timeUsec, err := f.dataFile.ReadUnival()
	if err != nil {
		return err
	}
	f.lastUpdate = time.Unix(timeSec.AsLong(), timeUsec.AsLong()*1000)

	return nil
}
