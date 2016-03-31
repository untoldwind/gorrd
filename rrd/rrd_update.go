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

	elapsed := r.calculateElapsedSteps(timestamp, interval)

	if elapsed.Steps == 0 {
		if err := r.simpleUpdate(newPdps, interval); err != nil {
			return err
		}
	} else {
		pdpTemp, err := r.processAllPdp(newPdps, elapsed, interval)
		if err != nil {
			return err
		}
		rraStepCounts, err := r.updateAllCdpPreps(pdpTemp, elapsed)
		if err != nil {
			return err
		}
		if err := r.updateAberrantCdps(pdpTemp, elapsed); err != nil {
			return err
		}
		if err := r.writeToRras(rraStepCounts); err != nil {
			return err
		}
		for i, rra := range r.Rras {
			if err := r.Store.StoreRraParams(i, rra); err != nil {
				return err
			}
		}
	}

	return r.writeChanges(timestamp)
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

func (r *Rrd) updatePdpPrep(newPdps []float64) error {
	return nil
}

func (r *Rrd) simpleUpdate(newPdps []float64, interval float64) error {
	for i, newPdp := range newPdps {
		r.Datasources[i].UpdatePdp(newPdp, interval)
	}

	return nil
}

func (r *Rrd) processAllPdp(newPdps []float64, elapsed ElapsedPdpSteps, interval float64) ([]float64, error) {
	pdpTemp := make([]float64, len(newPdps))
	for i, newPdp := range newPdps {
		pdpTemp[i] = r.Datasources[i].ProcessPdp(newPdp, interval, elapsed, r.Step)
	}
	return pdpTemp, nil
}

func (r *Rrd) updateAllCdpPreps(pdpTemp []float64, elapsed ElapsedPdpSteps) ([]uint64, error) {
	rraStepCounts := make([]uint64, len(r.Rras))
	var err error
	for i, rra := range r.Rras {
		rraStepCounts[i], err = rra.UpdateCdpPreps(pdpTemp, elapsed)
		if err != nil {
			return nil, err
		}
	}
	return rraStepCounts, nil
}

func (r *Rrd) updateAberrantCdps(pdpTemp []float64, elapsed ElapsedPdpSteps) error {
	first := true
	for j := elapsed.Steps; j > 0 && j < 3; j-- {
		for _, rra := range r.Rras {
			if err := rra.UpdateAberantCdp(pdpTemp, first); err != nil {
				return nil
			}
		}
		first = false
	}
	return nil
}

func (r *Rrd) writeToRras(rraStepCounts []uint64) error {
	for i, rra := range r.Rras {
		for first := true; rraStepCounts[i] > 0; rraStepCounts[i]-- {
			if first {
				r.Store.StoreRow(i, rra.GetPrimaryValues())
			} else {
				r.Store.StoreRow(i, rra.GetSecondaryValues())
			}
			first = false
		}
	}
	return nil
}

func (r *Rrd) writeChanges(timestamp time.Time) error {
	if err := r.Store.StoreLastUpdate(timestamp); err != nil {
		return err
	}

	for i, datasource := range r.Datasources {
		if err := r.Store.StoreDatasourceParams(i, datasource); err != nil {
			return err
		}
	}
	return r.Store.StoreRraPtrs()
}
