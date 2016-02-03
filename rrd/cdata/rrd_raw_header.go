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

func readRawHeader(dataFile *CDataFile) (*rrdRawHeader, error) {
	if cookie, err := dataFile.ReadCString(4); err != nil {
		return nil, err
	} else if cookie != rrdCookie {
		return nil, errors.Errorf("Invalid cookie: %+v", cookie)
	}

	if versionStr, err := dataFile.ReadCString(5); err != nil {
		return nil, err
	} else if version, err := strconv.ParseInt(string(versionStr[:4]), 10, 8); err != nil {
		return nil, errors.Errorf("Invalid version: %+v", version)
	} else if version < 3 {
		return nil, errors.Errorf("Version %d not supported: ", version)
	}
	if floatCookie, err := dataFile.ReadDouble(); err != nil {
		return nil, err
	} else if floatCookie != rrdFloatCookie {
		return nil, errors.Errorf("Float cookie does not match: %+v != %+v", floatCookie, rrdFloatCookie)
	}

	datasourceCount, err := dataFile.ReadUnsignedLong()
	if err != nil {
		return nil, err
	}
	rraCount, err := dataFile.ReadUnsignedLong()
	if err != nil {
		return nil, err
	}
	pdpStep, err := dataFile.ReadUnsignedLong()
	if err != nil {
		return nil, err
	}
	if _, err = dataFile.ReadUnivals(10); err != nil {
		return nil, err
	}

	return &rrdRawHeader{
		datasourceCount: datasourceCount,
		rraCount:        rraCount,
		pdpStep:         pdpStep,
	}, nil
}
