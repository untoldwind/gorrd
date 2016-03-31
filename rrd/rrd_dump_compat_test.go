package rrd_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/untoldwind/gorrd/rrd/cdata"
	"github.com/untoldwind/gorrd/rrd/dump"
)

func TestDumpCompatibility(t *testing.T) {
	rrdtool, err := findRrdTool()

	if err != nil {
		t.Skipf("rrdtool not found: %s", err.Error())
		return
	}

	tempDir := os.TempDir()
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10
	properties := gopter.NewProperties(parameters)

	minGauge, maxGauge := 0, 100000
	minCounter, maxCounter := 0, 1000
	minDerive, maxDerive := 0, 10000
	rrdStart := 1455218381
	properties.Property("dump of gauge, counter, derive is compatile", prop.ForAll(
		func(gauges, counters, derives []int) (bool, error) {
			rrdFileName := filepath.Join(tempDir, fmt.Sprintf("comp_update-%d-%d.rrd", time.Now().UnixNano(), rrdStart))
			defer os.Remove(rrdFileName)

			rrdStart++
			if err := rrdtool.create(rrdFileName,
				strconv.Itoa(rrdStart),
				"1",
				fmt.Sprintf("DS:watts:GAUGE:300:%d:%d", minGauge, maxGauge),
				fmt.Sprintf("DS:counts:COUNTER:300:%d:%d", minCounter, maxCounter),
				fmt.Sprintf("DS:derive:DERIVE:300:%d:%d", minDerive, maxDerive),
				"RRA:AVERAGE:0.5:1:100",
			); err != nil {
				return false, err
			}

			numUpdates := len(gauges)
			if numUpdates < len(counters) {
				numUpdates = len(counters)
			}
			if numUpdates < len(derives) {
				numUpdates = len(derives)
			}
			updates := make([]string, numUpdates)
			for i := range updates {
				gauge := "U"
				if i < len(gauges) {
					gauge = strconv.Itoa(gauges[i])
				}
				count := "U"
				if i < len(counters) {
					count = strconv.Itoa(counters[i])
				}
				derive := "U"
				if i < len(derives) {
					derive = strconv.Itoa(derives[i])
				}
				updates[i] = fmt.Sprintf("%d:%s:%s:%s", rrdStart+i+1, gauge, count, derive)
			}
			rrdtool.update(rrdFileName, updates...)

			expectedResult, err := rrdtool.dump(rrdFileName)

			if err != nil {
				return false, err
			}

			actualResult, err := runDumpCommand(rrdFileName)

			if err != nil {
				return false, err
			}

			return reflect.DeepEqual(expectedResult, actualResult), nil
		},
		gen.SliceOf(gen.IntRange(minGauge, maxGauge)),
		gen.SliceOf(gen.IntRange(minCounter, maxCounter)).Map(integrateInts),
		gen.SliceOf(gen.IntRange(minDerive, maxDerive)).Map(integrateInts),
	))

	properties.TestingRun(t)
}

func runDumpCommand(rrdFileName string) (map[string]string, error) {
	pipeReader, pipeWriter := io.Pipe()
	go func() {
		defer pipeWriter.Close()

		rrd, err := cdata.OpenRrdRawFile(rrdFileName, true)
		if err != nil {
			return
		}
		defer rrd.Close()

		xmlDumper, err := dump.NewXmlOutput(pipeWriter, true)
		if err != nil {
			return
		}
		if err := rrd.DumpTo(xmlDumper); err != nil {
			return
		}
	}()

	return flattenXml(pipeReader)
}

func integrateInts(v interface{}) interface{} {
	values := v.([]int)
	integral := 0
	for i, v := range values {
		integral += v
		values[i] = integral
	}
	return values
}
