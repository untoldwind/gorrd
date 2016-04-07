package rrd_test

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/exec"
	"strings"

	"github.com/go-errors/errors"
)

type rrdTool string

func findRrdTool() (rrdTool, error) {
	rrdExec, err := exec.LookPath("rrdtool")

	if err != nil {
		return "", err
	}
	return rrdTool(rrdExec), nil
}

func (r rrdTool) create(rrdFileName, start, step string, args ...string) error {
	createArgs := make([]string, 0, len(args)+6)
	createArgs = append(createArgs, "create")
	createArgs = append(createArgs, rrdFileName)
	createArgs = append(createArgs, "--start")
	createArgs = append(createArgs, start)
	createArgs = append(createArgs, "--step")
	createArgs = append(createArgs, step)
	createArgs = append(createArgs, args...)

	cmd := exec.Command(string(r), createArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (r rrdTool) update(rrdFileName string, updates ...string) error {
	if len(updates) == 0 {
		return nil
	}
	updateArgs := make([]string, 0, len(updates)+2)
	updateArgs = append(updateArgs, "update")
	updateArgs = append(updateArgs, rrdFileName)
	updateArgs = append(updateArgs, updates...)

	cmd := exec.Command(string(r), updateArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (r rrdTool) dump(rrdFileName string) (string, error) {
	output := bytes.NewBufferString("")

	cmd := exec.Command(string(r), "dump", rrdFileName)
	cmd.Stdout = output
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return output.String(), nil
}

type elementRef struct {
	name  string
	count int
}

func (e *elementRef) Inc() {
	e.count++
}

func (e *elementRef) String() string {
	return fmt.Sprintf("%s[%d]", e.name, e.count)
}

func flattenXml(xmlStr string) (map[string]interface{}, error) {
	decoder := xml.NewDecoder(bytes.NewBufferString(xmlStr))
	result := make(map[string]interface{}, 0)

	buffer := bytes.NewBufferString("")
	elementStack := make([]*elementRef, 0)
	var last *elementRef
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return result, nil
		} else if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		switch token.(type) {
		case xml.StartElement:
			name := token.(xml.StartElement).Name.Local
			if last != nil && last.name == name {
				last.Inc()
				elementStack = append(elementStack, last)
			} else {
				last = nil
				elementStack = append(elementStack, &elementRef{name: name})
			}
			buffer.Reset()
		case xml.EndElement:
			key := ""
			for i, elementRef := range elementStack {
				if i > 0 {
					key += "/"
				}
				key += elementRef.String()
			}
			text := strings.TrimSpace(buffer.String())
			floatVal, _, err := big.ParseFloat(text, 10, 40, big.ToNearestEven)
			if err == nil {
				result[key] = floatVal
			} else {
				result[key] = text
			}
			elementStack, last = elementStack[0:len(elementStack)-1], elementStack[len(elementStack)-1]
		case xml.CharData:
			buffer.Write(token.(xml.CharData))
		}
	}
}

func compareXml(xml1Str, xml2Str string) (string, error) {
	xml1, err := flattenXml(xml1Str)
	if err != nil {
		return "", err
	}
	xml2, err := flattenXml(xml2Str)
	if err != nil {
		return "", err
	}
	diffs := make([]string, 0)
	for k, v1 := range xml1 {
		v2, ok := xml2[k]
		if !ok {
			diffs = append(diffs, fmt.Sprintf("xml2 does not have %s", k))
		} else if !almoastEqual(v1, v2) {
			diffs = append(diffs, fmt.Sprintf("diff %s: %#v != %#v", k, v1, v2))
		}
	}
	for k := range xml2 {
		_, ok := xml1[k]
		if !ok {
			diffs = append(diffs, fmt.Sprintf("xml1 does not have %s", k))
		}
	}
	return strings.Join(diffs, "\n"), nil
}

func almoastEqual(a, b interface{}) bool {
	floatA, okA := a.(*big.Float)
	floatB, okB := b.(*big.Float)

	if okA && okB {
		return floatA.Cmp(floatB) == 0
	}
	return a == b
}

func copyFile(src, dst string) error {
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !stat.Mode().IsRegular() {
		return fmt.Errorf("non-regular source file %s %s", stat.Name(), stat.Mode().String())
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func shouldHaveSameContentAs(actual interface{}, expected ...interface{}) string {
	if len(expected) != 1 {
		return "shouldHaveSameContentAs must have only one expected value"
	}
	actualIn, err := os.Open(actual.(string))
	if err != nil {
		return err.Error()
	}
	defer actualIn.Close()
	expectedIn, err := os.Open(expected[0].(string))
	if err != nil {
		return err.Error()
	}
	defer expectedIn.Close()

	actualPipeIn, actualPipeOut := io.Pipe()
	actualScanner := bufio.NewScanner(actualPipeIn)
	actualDumper := hex.Dumper(actualPipeOut)
	go func() {
		io.Copy(actualDumper, actualIn)
		actualDumper.Close()
		actualPipeOut.Close()
	}()

	expectedPipeIn, expectedPipeOut := io.Pipe()
	expectedScanner := bufio.NewScanner(expectedPipeIn)
	expectedDumper := hex.Dumper(expectedPipeOut)
	go func() {
		io.Copy(expectedDumper, expectedIn)
		expectedDumper.Close()
		expectedPipeOut.Close()
	}()

	for {
		actualNext := actualScanner.Scan()
		expectedNext := expectedScanner.Scan()
		if expectedNext && !actualNext {
			return fmt.Sprintf("%v and %v differ: extra Line %s", actual, expected[0], expectedScanner.Text())
		}
		if !expectedNext && actualNext {
			return fmt.Sprintf("%v and %v differ: extra Line %s", actual, expected[0], actualScanner.Text())
		}
		if actualScanner.Text() != expectedScanner.Text() {
			return fmt.Sprintf("%v and %v differ: %s != %s", actual, expected[0], actualScanner.Text(), expectedScanner.Text())
		}
		if !actualNext || !expectedNext {
			break
		}
	}

	return ""
}
