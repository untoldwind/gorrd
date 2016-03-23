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

	runDumpCommand := func(rrdFileName string) (map[string]string, error) {
		pipeReader, pipeWriter := io.Pipe()
		go func() {
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
			pipeWriter.Close()
		}()

		return flattenXml(pipeReader)
	}

	tempDir := os.TempDir()
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10
	properties := gopter.NewProperties(parameters)

	minGauge, maxGauge := 0, 100000
	rrdStart := 1455218381
	properties.Property("dump of single gauge is compatile", prop.ForAll(
		func(values []int) (bool, error) {
			rrdFileName := filepath.Join(tempDir, fmt.Sprintf("comp_update-%s-%d.rrd", time.Now().String(), rrdStart))
			defer os.Remove(rrdFileName)

			rrdStart++
			if err := rrdtool.create(rrdFileName,
				strconv.Itoa(rrdStart),
				"1",
				fmt.Sprintf("DS:watts:GAUGE:300:%d:%d", minGauge, maxGauge),
				"RRA:AVERAGE:0.5:1:100",
			); err != nil {
				return false, err
			}

			updates := make([]string, len(values))
			for i, value := range values {
				updates[i] = fmt.Sprintf("%d:%d", rrdStart+i+1, value)
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
		gen.SliceOf(gen.IntRange(minGauge, maxGauge))))

	properties.TestingRun(t)
}
