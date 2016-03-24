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

func (d *DatasourceCompute) CalculatePdpPrep(newValue string, interval float64) (float64, error) {
	// Compute datasource are never updated
	return math.NaN(), nil
}

func (d *DatasourceCompute) UpdatePdp(pdpValue, interval float64) {

}

func (d *DatasourceCompute) DumpTo(dumper DataOutput) error {
	dumper.DumpString("name", d.Name)
	dumper.DumpString("type", DatasourceTypeCompute)
	return dumper.Finalize()
}
