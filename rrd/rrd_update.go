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
	newPdps := make([]float64, len(r.Datasources))
	var err error
	for i, datasource := range r.Datasources {
		newPdps[i], err = datasource.CalculatePdpPrep(values[i], interval)
		if err != nil {
			return err
		}
	}

	elapsedSteps := int64(timestamp.Sub(r.LastUpdate.Truncate(r.Step)) / r.Step)

	if elapsedSteps == 0 {
		if err := r.simpleUpdate(newPdps, interval); err != nil {
			return err
		}
	}

	return r.Store.StoreLastUpdate(timestamp)
}

func (r *Rrd) updatePdpPrep(timestamp time.Time, values []string) ([]float64, error) {
	result := make([]float64, len(r.Datasources))
	interval := float64(timestamp.Sub(r.LastUpdate).Nanoseconds()) / 1e9

	var err error
	for i, datasource := range r.Datasources {
		result[i], err = datasource.CalculatePdpPrep(values[i], interval)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
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
