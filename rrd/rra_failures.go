package rrd

import "math"

const RraTypeFailures = "FAILURES"

type RraCpdPrepFailures struct {
	RraCpdPrepBase
}

type RraHwFailures struct {
	RraAbstractGeneric
}

func (r *RraHwFailures) DumpTo(rrdStore Store, dumper DataOutput) {
	dumper.DumpString("cf", RraTypeFailures)
	r.RraAbstractGeneric.DumpTo(rrdStore, dumper)
}

func newRraFailures(index int, store Store) (*RraHwFailures, error) {
	result := &RraHwFailures{
		newRraAbstractGeneric(index, math.NaN()),
	}

	result.InitializeCdpFunc = func(pdpTemp float64, pdpPerRow, startPdpOffset uint64, cpdPrep *RraCpdPrepGeneric) {
		cpdPrep.PrimaryValue = pdpTemp
	}

	if err := store.ReadRraParams(index, result); err != nil {
		return nil, err
	}

	return result, nil
}
