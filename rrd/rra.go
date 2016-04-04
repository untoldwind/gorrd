package rrd

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
