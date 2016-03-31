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
		r.simpleUpdate(newPdps, interval)
	} else {
		pdpTemp := r.processAllPdp(newPdps, elapsed, interval)
		rraStepCounts := r.updateAllCdpPreps(pdpTemp, elapsed)
		r.updateAberrantCdps(pdpTemp, elapsed)
		r.writeToRras(rraStepCounts)
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

func (r *Rrd) simpleUpdate(newPdps []float64, interval float64) {
	for i, newPdp := range newPdps {
		r.Datasources[i].UpdatePdp(newPdp, interval)
	}
}

func (r *Rrd) processAllPdp(newPdps []float64, elapsed ElapsedPdpSteps, interval float64) []float64 {
	pdpTemp := make([]float64, len(newPdps))
	for i, newPdp := range newPdps {
		pdpTemp[i] = r.Datasources[i].ProcessPdp(newPdp, interval, elapsed, r.Step)
	}
	return pdpTemp
}

func (r *Rrd) updateAllCdpPreps(pdpTemp []float64, elapsed ElapsedPdpSteps) []uint64 {
	rraStepCounts := make([]uint64, len(r.Rras))
	for i, rra := range r.Rras {
		rraStepCounts[i] = rra.UpdateCdpPreps(pdpTemp, elapsed)
	}
	return rraStepCounts
}

func (r *Rrd) updateAberrantCdps(pdpTemp []float64, elapsed ElapsedPdpSteps) {
	first := true
	for j := elapsed.Steps; j > 0 && j < 3; j-- {
		for _, rra := range r.Rras {
			rra.UpdateAberantCdp(pdpTemp, first)
		}
		first = false
	}
}

func (r *Rrd) writeToRras(rraStepCounts []uint64) {
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
