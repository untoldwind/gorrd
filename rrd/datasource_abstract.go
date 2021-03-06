package rrd

import (
	"math"
	"time"
)

type DatasourceAbstract struct {
	Name            string  `ds:"name"`
	Heartbeat       uint64  `ds:"param0"`
	Min             float64 `ds:"param1"`
	Max             float64 `ds:"param2"`
	LastValue       string  `pdp:"lastValue"`
	UnknownSecCount uint64  `pdp:"0"`
	PdpValue        float64 `pdp:"1"`
}

func (d *DatasourceAbstract) GetName() string {
	return d.Name
}

func (d *DatasourceAbstract) SetLastValue(lastValue string) {
	d.LastValue = lastValue
}

func (d *DatasourceAbstract) GetLastValue() string {
	return d.LastValue
}

func (d *DatasourceAbstract) UpdatePdp(pdpValue, interval float64) {
	if math.IsNaN(pdpValue) {
		d.UnknownSecCount += uint64(interval)
	} else if math.IsNaN(d.PdpValue) {
		d.PdpValue = pdpValue
	} else {
		d.PdpValue += pdpValue
	}
}

func (d *DatasourceAbstract) ProcessPdp(pdpValue float64, elapsed ElapsedPdpSteps, step time.Duration) float64 {
	var preUnknown float64
	if math.IsNaN(pdpValue) {
		preUnknown = elapsed.PreInt
	} else {
		if math.IsNaN(d.PdpValue) {
			d.PdpValue = 0
		}
		d.PdpValue += pdpValue / elapsed.Interval * elapsed.PreInt
	}
	var pdpTemp float64

	if elapsed.Interval > float64(d.Heartbeat) || uint64(step/time.Second/2) < d.UnknownSecCount {
		pdpTemp = math.NaN()
	} else {
		diffPdpSteps := (elapsed.Steps * uint64(step)) / uint64(time.Second)
		pdpTemp = d.PdpValue / (float64(diffPdpSteps-d.UnknownSecCount) - preUnknown)
	}

	if math.IsNaN(pdpValue) {
		d.UnknownSecCount = uint64(elapsed.PostInt)
		d.PdpValue = math.NaN()
	} else {
		d.UnknownSecCount = 0
		d.PdpValue = pdpValue / elapsed.Interval * elapsed.PostInt
	}

	return pdpTemp
}

func (d *DatasourceAbstract) DumpTo(dumper DataOutput) error {
	dumper.DumpString("name", d.Name)
	dumper.DumpUnsignedLong("minimal_heartbeat", d.Heartbeat)
	dumper.DumpDouble("min", d.Min)
	dumper.DumpDouble("max", d.Max)
	dumper.DumpString("last_ds", d.LastValue)
	dumper.DumpDouble("value", d.PdpValue)
	dumper.DumpUnsignedLong("unknown_sec", d.UnknownSecCount)
	return nil
}

func (d *DatasourceAbstract) checkRateBounds(rate float64) bool {
	return !math.IsNaN(rate) &&
		(math.IsNaN(d.Min) || rate >= d.Min) &&
		(math.IsNaN(d.Max) || rate <= d.Max)
}
