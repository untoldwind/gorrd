package rrd

const RrdDatasourceTypeCompute = "COMPUTE"

type RrdComputeDatasource struct {
	Name string
}

func (d *RrdComputeDatasource) GetName() string {
	return d.Name
}
