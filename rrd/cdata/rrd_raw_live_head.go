package cdata

import (
	"github.com/go-errors/errors"
	"time"
)

const rrdRawLiveHeaderSize = 16

func (f *RrdRawFile) readLiveHead(reader *RawDataReader) error {
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

func (f *RrdRawFile) StoreLastUpdate(lastUpdate time.Time) error {
	writer := f.dataFile.Writer(f.baseHeaderSize)
	if err := writer.WriteUnival(unival(lastUpdate.Unix())); err != nil {
		return errors.Wrap(err, 0)
	}
	return writer.WriteUnival(unival(lastUpdate.Nanosecond() / 1000))
}
