package commands_test

import (
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

	count := 0
	actualBytes := make([]byte, 16)
	expectedBytes := make([]byte, 16)
	for {
		actualCount, err := actualIn.Read(actualBytes)
		if err != nil && err != io.EOF {
			return err.Error()
		}
		expectedCount, err := expectedIn.Read(expectedBytes)
		if err != nil && err != io.EOF {
			return err.Error()
		}
		actualHex := hex.Dump(actualBytes[0:actualCount])
		expectedHex := hex.Dump(expectedBytes[0:expectedCount])
		if actualHex != expectedHex {
			return fmt.Sprintf("%v and %v differ at %d: %s != %s", actual, expected[0], count, actualHex, expectedHex)
		}
		if actualCount == 0 || expectedCount == 0 {
			break
		}
		count += actualCount
	}
	return ""
}
