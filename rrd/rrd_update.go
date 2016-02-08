package rrd

import (
	"time"

	"github.com/go-errors/errors"
)

func (r *Rrd) Update(timestamp time.Time, values []string) error {
	if timestamp.Before(r.LastUpdate) {
		return errors.Errorf("illegal attempt to update using time %s when last update time is %s (minimum one second step)", timestamp.String(), r.LastUpdate.String())
	}
	elapsedSteps := int64(timestamp.Sub(r.LastUpdate) / r.Step)

	if elapsedSteps == 0 {

	}

	return nil
}

func (r *Rrd) updatePdpPrep(timestamp time.Time, values []string) ([]float64, error) {
	result := make([]float64, len(r.Datasources))
	interval := float64(timestamp.Sub(r.LastUpdate).Nanoseconds()) / 1e9

	var err error
	for i, datasource := range r.Datasources {
		result[i], err = datasource.UpdatePdpPrep(values[i], interval)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (r *Rrd) simpleUpdate(values []string) error {
	return nil
}
