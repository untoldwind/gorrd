package cdata

import (
	"strconv"

	"github.com/go-errors/errors"
)

type rrdRawHeader struct {
	datasourceCount uint64
	rraCount        uint64
	pdpStep         uint64
}

func (f *RrdRawFile) readVersionHeader(reader *CDataReader) error {
	if cookie, err := reader.ReadCString(4); err != nil {
		return err
	} else if cookie != rrdCookie {
		return errors.Errorf("Invalid cookie: %+v", cookie)
	}

	if versionStr, err := reader.ReadCString(5); err != nil {
		return err
	} else if version, err := strconv.ParseInt(string(versionStr[:4]), 10, 8); err != nil {
		return errors.Errorf("Invalid version: %+v", version)
	} else if version < 3 {
		return errors.Errorf("Version %d not supported: ", version)
	}
	if floatCookie, err := reader.ReadDouble(); err != nil {
		return err
	} else if floatCookie != rrdFloatCookie {
		return errors.Errorf("Float cookie does not match: %+v != %+v", floatCookie, rrdFloatCookie)
	}

	datasourceCount, err := reader.ReadUnsignedLong()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	rraCount, err := reader.ReadUnsignedLong()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	pdpStep, err := reader.ReadUnsignedLong()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	if _, err = reader.ReadUnivals(10); err != nil {
		return errors.Wrap(err, 0)
	}

	f.header = &rrdRawHeader{
		datasourceCount: datasourceCount,
		rraCount:        rraCount,
		pdpStep:         pdpStep,
	}
	return nil
}
