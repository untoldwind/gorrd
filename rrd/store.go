package rrd

import "time"

type Store interface {
	DatasourceTypes() []string
	ReadDatasourceParams(index int, params interface{}) error
	RraTypes() []string
	ReadRraParams(index int, params interface{}) error
	LastUpdate() time.Time
	Step() uint64
	RowIterator(rra Rra) (RraRowIterator, error)
	Close()
}
