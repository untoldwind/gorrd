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
		pdpTemp, err := r.processAllPdp(newPdps, elapsedSteps, procPdpCount, interval, preInt, postInt)
		if err != nil {
			return err
		}
		rraStepCounts, err := r.updateAllCdpPreps(pdpTemp, elapsedSteps, procPdpCount)
		if err != nil {
			return err
		}
		if err := r.updateAberrantCdps(pdpTemp, elapsedSteps); err != nil {
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

func (r *Rrd) calculateElapsedSteps(timestamp time.Time, interval float64) (uint64, float64, float64, uint64) {
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

	return uint64(occuPdpSt-procPdpSt) / uint64(r.Step/time.Second), preInt, postInt, uint64(procPdpCount)
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

func (r *Rrd) processAllPdp(newPdps []float64, elapsedSteps, procPdpCount uint64, interval, preInt, postInt float64) ([]float64, error) {
	pdpTemp := make([]float64, len(newPdps))
	for i, newPdp := range newPdps {
		pdpTemp[i] = r.Datasources[i].ProcessPdp(newPdp, interval, preInt, postInt, elapsedSteps, r.Step)
	}
	return pdpTemp, nil
}

func (r *Rrd) updateAllCdpPreps(pdpTemp []float64, elapsedSteps, procPdpCount uint64) ([]uint64, error) {
	rraStepCounts := make([]uint64, len(r.Rras))
	var err error
	for i, rra := range r.Rras {
		rraStepCounts[i], err = rra.UpdateCdpPreps(pdpTemp, elapsedSteps, procPdpCount)
		if err != nil {
			return nil, err
		}
	}
	return rraStepCounts, nil
}

func (r *Rrd) updateAberrantCdps(pdpTemp []float64, elapsedSteps uint64) error {
	first := true
	for j := elapsedSteps; j > 0 && j < 3; j-- {
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
