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

	elapsedSteps := int64(timestamp.Sub(r.LastUpdate.Truncate(r.Step)) / r.Step)

	if elapsedSteps == 0 {
		if err := r.simpleUpdate(newPdps, interval); err != nil {
			return err
		}
	} else {
		if err := r.processAllPdp(elapsedSteps); err != nil {
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

func (r *Rrd) calculateElapsedSteps(timestamp time.Time, interval float64) (int64, float64, float64, int64) {
	proc_pdp_age := r.LastUpdate.Unix() % int64(r.Step/time.Second)
	proc_pdp_st := r.LastUpdate.Unix() - proc_pdp_age

	occu_pdp_age := timestamp.Unix() % int64(r.Step/time.Second)
	occu_pdp_st := timestamp.Unix() - occu_pdp_age

	var pre_int float64
	var post_int float64
	if occu_pdp_st > proc_pdp_st {
		pre_int = float64(occu_pdp_st - r.LastUpdate.Unix())
		pre_int -= float64(r.LastUpdate.Nanosecond()) / 1e9
		post_int = float64(occu_pdp_age)
		post_int += float64(timestamp.Nanosecond()) / 1e9
	} else {
		pre_int = interval
		post_int = 0
	}

	proc_pdp_cnt := proc_pdp_st / int64(r.Step/time.Second)

	return (occu_pdp_st - proc_pdp_st) / int64(r.Step/time.Second), pre_int, post_int, proc_pdp_cnt
}

func (r *Rrd) simpleUpdate(pdpPreps []float64, interval float64) error {
	for i, pdpPrep := range pdpPreps {
		r.Datasources[i].UpdatePdp(pdpPrep, interval)
	}

	for i, datasource := range r.Datasources {
		if err := r.Store.StoreDatasourceParams(i, datasource); err != nil {
			return err
		}
	}
	return nil
}

func (r *Rrd) processAllPdp(elapsedSteps int64) error {
	return nil
}

func (r *Rrd) processPdp(dataource Datasource) error {
	return nil
}
