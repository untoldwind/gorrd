package cdata

import "time"

func (f *RrdRawFile) readLiveHead(reader *CDataReader) error {
	timeSec, err := reader.ReadUnival()
	if err != nil {
		return err
	}
	timeUsec, err := reader.ReadUnival()
	if err != nil {
		return err
	}
	f.lastUpdate = time.Unix(timeSec.AsLong(), timeUsec.AsLong()*1000)
	return nil
}
