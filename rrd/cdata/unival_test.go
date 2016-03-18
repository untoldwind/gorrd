package cdata

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestUnival(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("Unival from uint64", prop.ForAll(
		func(i uint64) bool {
			val := unival(i)

			return val.AsUnsignedLong() == i
		},
		gen.UInt64()))

	properties.Property("Unival from int64", prop.ForAll(
		func(i int64) bool {
			val := unival(i)

			return val.AsLong() == i
		},
		gen.Int64()))

	properties.Property("Unival from float64", prop.ForAll(
		func(f float64) bool {
			val := univalForDouble(f)

			return val.AsDouble() == f
		},
		gen.Float64()))

	properties.TestingRun(t)
}
