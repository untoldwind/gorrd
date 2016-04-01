package rrd

import "time"

type ElapsedPdpSteps struct {
	Interval     float64
	Steps        uint64
	PreInt       float64
	PostInt      float64
	ProcPdpCount uint64
}

func (r *Rrd) calculateElapsedSteps(timestamp time.Time) ElapsedPdpSteps {
	interval := float64(timestamp.Sub(r.LastUpdate).Nanoseconds()) / 1e9

	procPdpAge := r.LastUpdate.Unix() % int64(r.Step/time.Second)
	procPdpSt := r.LastUpdate.Unix() - procPdpAge

	occuPdpAge := timestamp.Unix() % int64(r.Step/time.Second)
	occuPdpSt := timestamp.Unix() - occuPdpAge

	var preInt float64
	var postInt float64
	if occuPdpSt > procPdpSt {
		preInt = float64(occuPdpSt - r.LastUpdate.Unix())
		preInt -= float64(r.LastUpdate.Nanosecond()) / 1e9
		postInt = float64(occuPdpAge)
		postInt += float64(timestamp.Nanosecond()) / 1e9
	} else {
		preInt = interval
		postInt = 0
	}

	procPdpCount := procPdpSt / int64(r.Step/time.Second)

	return ElapsedPdpSteps{
		Interval:     interval,
		Steps:        uint64(occuPdpSt-procPdpSt) / uint64(r.Step/time.Second),
		PreInt:       preInt,
		PostInt:      postInt,
		ProcPdpCount: uint64(procPdpCount),
	}
}
