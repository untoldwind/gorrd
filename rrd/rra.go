package rrd

import "github.com/go-errors/errors"

// Rra is the generic interface to access a round robin archive.
type Rra interface {
	GetRowCount() uint64
	GetPdpPerRow() uint64
	GetPrimaryValues() []float64
	GetSecondaryValues() []float64
	UpdateCdpPreps(pdpTemp []float64, elapsed ElapsedPdpSteps) uint64
	UpdateAberantCdp(pdpTemp []float64, first bool)
	// DumpTo dumps the content of the archive to a data dumper.
	DumpTo(rrdStore Store, dumper DataOutput)
}

func newRra(index int, rraType string, store Store) (Rra, error) {
	switch rraType {
	case RraTypeAverage:
		return newRraAverage(index, store)
	case RraTypeFailures:
		return newRraFailures(index, store)
	case RraTypeMin:
		return newRraMin(index, store)
	case RraTypeMax:
		return newRraMax(index, store)
	case RraTypeLast:
		return newRraLast(index, store)
	}
	return nil, errors.Errorf("Unknown rra type: %s", rraType)
}
