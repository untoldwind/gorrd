package cdata

import "time"

const rrdRawLiveHeaderSize = 16

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

func (f *RrdRawFile) StoreLastUpdate(lastUpdate time.Time) {
	writer := f.dataFile.Writer(f.baseHeaderSize)
	writer.WriteUnival(unival(lastUpdate.Unix()))
	writer.WriteUnival(unival(lastUpdate.UnixNano() / 1000))
}
