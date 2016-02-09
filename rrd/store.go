package rrd

import "time"

type Store interface {
	DatasourceTypes() []string
	RraTypes() []string
	LastUpdate() time.Time
	Step() time.Duration

	ReadDatasourceParams(index int, params interface{}) error
	ReadRraParams(index int, params interface{}) error
	StoreLastUpdate(lastUpdate time.Time)

	RowIterator(rraIndex int) (RraRowIterator, error)
	Close()
}
