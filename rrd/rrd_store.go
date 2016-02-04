package rrd

type RrdStore interface {
	RowIterator(rra Rra) (RraRowIterator, error)
	Close()
}
