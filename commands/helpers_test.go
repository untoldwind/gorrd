package commands_test

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

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

func flattenXml(in io.Reader) (map[string]string, error) {
	decoder := xml.NewDecoder(in)
	result := make(map[string]string, 0)

	buffer := bytes.NewBufferString("")
	elementStack := make([]*elementRef, 0)
	var last *elementRef
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return result, nil
		} else if err != nil {
			return nil, err
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
			result[key] = strings.TrimSpace(buffer.String())
			elementStack, last = elementStack[0:len(elementStack)-1], elementStack[len(elementStack)-1]
		case xml.CharData:
			buffer.Write(token.(xml.CharData))
		}
	}
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