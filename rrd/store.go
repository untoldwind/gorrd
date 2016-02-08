package rrd

import "time"

type Store interface {
	DatasourceTypes() []string
	ReadDatasourceParams(index int, params interface{}) error
	RraTypes() []string
	ReadRraParams(index int, params interface{}) error
	LastUpdate() time.Time
	Step() time.Duration
	RowIterator(rraIndex int) (RraRowIterator, error)
	Close()
}
