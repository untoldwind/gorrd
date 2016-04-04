package rrd

import "time"

type Store interface {
	Version() uint16
	DatasourceTypes() []string
	RraTypes() []string
	LastUpdate() time.Time
	Step() time.Duration

	ReadDatasourceParams(index int, params interface{}) error
	StoreDatasourceParams(index int, params interface{}) error
	ReadRraParams(index int, params interface{}) error
	StoreRraParams(index int, params interface{}) error
	StoreLastUpdate(lastUpdate time.Time) error
	StoreRraPtrs() error

	RowIterator(rraIndex int) (RraRowIterator, error)
	StoreRow(rraIndex int, row []float64) error
	Close()
}
