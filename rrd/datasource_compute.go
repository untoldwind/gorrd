package rrd

import "math"

const DatasourceTypeCompute = "COMPUTE"

type DatasourceCompute struct {
	Name      string
	LastValue string
}

func (d *DatasourceCompute) GetName() string {
	return d.Name
}

func (d *DatasourceCompute) GetLastValue() string {
	return d.LastValue
}

func (d *DatasourceCompute) UpdatePdpPrep(newValue string, interval float64) (float64, error) {
	// Compute datasource are never updated
	return math.NaN(), nil
}

func (d *DatasourceCompute) DumpTo(dumper DataOutput) error {
	if err := dumper.DumpString("name", d.Name); err != nil {
		return err
	}
	if err := dumper.DumpString("type", DatasourceTypeCompute); err != nil {
		return err
	}
	return dumper.Finalize()
}
