package rrd

import (
	"time"

	"github.com/go-errors/errors"
)

func (r *Rrd) Update(timestamp time.Time, values []string) error {
	if timestamp.Before(r.LastUpdate) {
		return errors.Errorf("illegal attempt to update using time %s when last update time is %s (minimum one second step)", timestamp.String(), r.LastUpdate.String())
	}

	interval := float64(timestamp.Sub(r.LastUpdate).Nanoseconds()) / 1e9
	newPdps, err := r.calculatePdpPreps(interval, values)
	if err != nil {
		return err
	}

	elapsedSteps, preInt, postInt, procPdpCount := r.calculateElapsedSteps(timestamp, interval)

	if elapsedSteps == 0 {
		if err := r.simpleUpdate(newPdps, interval); err != nil {
			return err
		}
	} else {
		if err := r.processAllPdp(newPdps, elapsedSteps, procPdpCount, interval, preInt, postInt); err != nil {
			return err
		}
	}

	return r.Store.StoreLastUpdate(timestamp)
}

func (r *Rrd) calculatePdpPreps(interval float64, values []string) ([]float64, error) {
	result := make([]float64, len(r.Datasources))

	var err error
	for i, datasource := range r.Datasources {
		result[i], err = datasource.CalculatePdpPrep(values[i], interval)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (r *Rrd) calculateElapsedSteps(timestamp time.Time, interval float64) (uint64, float64, float64, uint64) {
	proc_pdp_age := r.LastUpdate.Unix() % int64(r.Step/time.Second)
	proc_pdp_st := r.LastUpdate.Unix() - proc_pdp_age

	occu_pdp_age := timestamp.Unix() % int64(r.Step/time.Second)
	occu_pdp_st := timestamp.Unix() - occu_pdp_age

	var preInt float64
	var postInt float64
	if occu_pdp_st > proc_pdp_st {
		preInt = float64(occu_pdp_st - r.LastUpdate.Unix())
		preInt -= float64(r.LastUpdate.Nanosecond()) / 1e9
		postInt = float64(occu_pdp_age)
		postInt += float64(timestamp.Nanosecond()) / 1e9
	} else {
		preInt = interval
		postInt = 0
	}

	procPdpCount := proc_pdp_st / int64(r.Step/time.Second)

	return uint64(occu_pdp_st-proc_pdp_st) / uint64(r.Step/time.Second), preInt, postInt, uint64(procPdpCount)
}

func (r *Rrd) simpleUpdate(newPdps []float64, interval float64) error {
	for i, newPdp := range newPdps {
		r.Datasources[i].UpdatePdp(newPdp, interval)
	}

	for i, datasource := range r.Datasources {
		if err := r.Store.StoreDatasourceParams(i, datasource); err != nil {
			return err
		}
	}
	return nil
}

func (r *Rrd) processAllPdp(newPdps []float64, elapsedSteps, procPdpCount uint64, interval, preInt, postInt float64) error {
	pdpTemp := make([]float64, len(newPdps))
	for i, newPdp := range newPdps {
		pdpTemp[i] = r.Datasources[i].ProcessPdp(newPdp, interval, preInt, elapsedSteps, r.Step)
	}
	return nil
}

func (r *Rrd) processPdp(dataource Datasource) error {
	return nil
}
