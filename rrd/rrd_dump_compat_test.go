package rrd_test

import (
	"bytes"
	"fmt"
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

const minGauge, maxGauge = 0, 100000
const minCounter, maxCounter = 0, 1000
const minDerive, maxDerive = -1000, 10000
const minAbsolute, maxAbsolute = 0, 10000
const minDCounter, maxDCounter = 0.0, 1000.0
const minDDerive, maxDDerive = -1000.0, 10000.0

func TestDumpCompatibility(t *testing.T) {
	rrdtool, err := findRrdTool()

	if err != nil {
		t.Skipf("rrdtool not found: %s", err.Error())
		return
	}

	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20
	properties := gopter.NewProperties(parameters)

	properties.Property("dump of gauge, counter, derive, absolute is compatile", prop.ForAllNoShrink(
		rrdtool.checkDumpCompatibility1,
		counterGen(1455218381),
		gen.SliceOf(gen.IntRange(minGauge, maxGauge)),
		gen.SliceOf(gen.IntRange(minCounter, maxCounter)).Map(integrateInts),
		gen.SliceOf(gen.IntRange(minDerive, maxDerive)).Map(integrateInts),
		gen.SliceOf(gen.IntRange(minAbsolute, maxAbsolute)).Map(integrateInts),
	))

	properties.Property("dump of dcounter, dderive is compatile", prop.ForAllNoShrink(
		rrdtool.checkDumpCompatibility2,
		counterGen(1455218381),
		gen.SliceOf(gen.Float64Range(minDCounter, maxDCounter)).Map(integrateFloats),
		gen.SliceOf(gen.Float64Range(minDDerive, maxDDerive)).Map(integrateFloats),
	))

	properties.TestingRun(t)
}

func (rrdtool rrdTool) checkDumpCompatibility1(rrdStart int, gauges, counters, derives, absolutes []int) (bool, error) {
	tempDir := os.TempDir()
	rrdFileName := filepath.Join(tempDir, fmt.Sprintf("comp_dump1-%d-%d.rrd", time.Now().UnixNano(), rrdStart))
	defer os.Remove(rrdFileName)
	//fmt.Println(rrdFileName)
	if err := rrdtool.create(rrdFileName,
		strconv.Itoa(rrdStart),
		"1",
		fmt.Sprintf("DS:watts:GAUGE:300:%d:%d", minGauge, maxGauge),
		fmt.Sprintf("DS:counts:COUNTER:300:%d:%d", minCounter, maxCounter),
		fmt.Sprintf("DS:derive:DERIVE:300:%d:%d", minDerive, maxDerive),
		fmt.Sprintf("DS:absolute:ABSOLUTE:300:%d:%d", minAbsolute, maxAbsolute),
		"RRA:AVERAGE:0.5:1:100",
		"RRA:MIN:0.5:1:100",
		"RRA:MAX:0.5:1:100",
		"RRA:LAST:0.5:1:100",
	//	"RRA:HWPREDICT:500:0.1:0.0035:100",
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
	if numUpdates < len(absolutes) {
		numUpdates = len(absolutes)
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
		absolute := "U"
		if i < len(absolutes) {
			absolute = strconv.Itoa(absolutes[i])
		}
		updates[i] = fmt.Sprintf("%d:%s:%s:%s:%s", rrdStart+i+1, gauge, count, derive, absolute)
	}
	if err := rrdtool.update(rrdFileName, updates...); err != nil {
		return false, err
	}

	expectedResult, err := rrdtool.dump(rrdFileName)

	if err != nil {
		return false, err
	}

	actualResult, err := runDumpCommand(rrdFileName)

	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(expectedResult, actualResult), nil
}

func (rrdtool rrdTool) checkDumpCompatibility2(rrdStart int, counters, derives []float64) (bool, error) {
	tempDir := os.TempDir()
	rrdFileName := filepath.Join(tempDir, fmt.Sprintf("comp_dump2-%d-%d.rrd", time.Now().UnixNano(), rrdStart))
	defer os.Remove(rrdFileName)

	if err := rrdtool.create(rrdFileName,
		strconv.Itoa(rrdStart),
		"1",
		fmt.Sprintf("DS:dcounter:DCOUNTER:300:%f:%f", minDCounter, maxDCounter),
		fmt.Sprintf("DS:dderive:DDERIVE:300:%f:%f", minDDerive, maxDDerive),
		"RRA:AVERAGE:0.5:1:100",
		"RRA:MIN:0.5:1:100",
		"RRA:MAX:0.5:1:100",
		"RRA:LAST:0.5:1:100",
	); err != nil {
		return false, err
	}

	numUpdates := len(counters)
	if numUpdates < len(derives) {
		numUpdates = len(derives)
	}
	updates := make([]string, numUpdates)
	for i := range updates {
		count := "U"
		if i < len(counters) {
			count = strconv.FormatFloat(counters[i], 'e', 8, 64)
		}
		derive := "U"
		if i < len(derives) {
			derive = strconv.FormatFloat(derives[i], 'e', 8, 64)
		}
		updates[i] = fmt.Sprintf("%d:%s:%s", rrdStart+i+1, count, derive)
	}
	if err := rrdtool.update(rrdFileName, updates...); err != nil {
		return false, err
	}

	expectedResult, err := rrdtool.dump(rrdFileName)

	if err != nil {
		return false, err
	}

	actualResult, err := runDumpCommand(rrdFileName)

	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(expectedResult, actualResult), nil
}

func counterGen(start int) gopter.Gen {
	return func(*gopter.GenParameters) *gopter.GenResult {
		start++
		return gopter.NewGenResult(start, gopter.NoShrinker)
	}
}

func runDumpCommand(rrdFileName string) (map[string]interface{}, error) {
	rrd, err := cdata.OpenRrdRawFile(rrdFileName, true)
	if err != nil {
		return nil, err
	}
	defer rrd.Close()

	buffer := bytes.NewBufferString("")
	xmlDumper, err := dump.NewXmlOutput(buffer, true)
	if err != nil {
		return nil, err
	}

	if err := rrd.DumpTo(xmlDumper); err != nil {
		return nil, err
	}

	return flattenXml(buffer.String())
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

func integrateFloats(v interface{}) interface{} {
	values := v.([]float64)
	integral := 0.0
	for i, v := range values {
		integral += v
		values[i] = integral
	}
	return values
}
