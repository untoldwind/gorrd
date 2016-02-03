package cdata

import "time"

func readLiveHead(dataFile *CDataFile) (time.Time, error) {
	timeSec, err := dataFile.ReadUnival()
	if err != nil {
		return time.Time{}, err
	}
	timeUsec, err := dataFile.ReadUnival()
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(timeSec.AsLong(), timeUsec.AsLong()*1000), nil
}
