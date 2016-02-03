package cdata

import (
	"strconv"

	"github.com/go-errors/errors"
)

func (f *RrdRawFile) readHeader() error {
	if cookie, err := f.dataFile.ReadCString(4); err != nil {
		return err
	} else if cookie != rrdCookie {
		return errors.Errorf("Invalid cookie: %+v", cookie)
	}

	if versionStr, err := f.dataFile.ReadCString(5); err != nil {
		return err
	} else if version, err := strconv.ParseInt(string(versionStr[:4]), 10, 8); err != nil {
		return errors.Errorf("Invalid version: %+v", version)
	} else if version < 3 {
		return errors.Errorf("Version %d not supported: ", version)
	}
	if floatCookie, err := f.dataFile.ReadDouble(); err != nil {
		return err
	} else if floatCookie != rrdFloatCookie {
		return errors.Errorf("Float cookie does not match: %+v != %+v", floatCookie, rrdFloatCookie)
	}
	var err error

	f.datasourceCount, err = f.dataFile.ReadUnsignedLong()
	if err != nil {
		return err
	}
	f.rraCount, err = f.dataFile.ReadUnsignedLong()
	if err != nil {
		return err
	}
	f.pdpStep, err = f.dataFile.ReadUnsignedLong()
	if err != nil {
		return err
	}
	_, err = f.dataFile.ReadUnivals(10)

	return nil
}
