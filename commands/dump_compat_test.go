package commands

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/codegangsta/cli"
	. "github.com/smartystreets/goconvey/convey"
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

func TestDumpCompatibility(t *testing.T) {
	rrdtool, err := exec.LookPath("rrdtool")

	if err != nil {
		t.Skipf("rrdtool not found: %s", err.Error())
		return
	}
	Convey("Given minimal rrdfile with 1s step", t, func() {
		tempDir := os.TempDir()
		rrdFileName := filepath.Join(tempDir, fmt.Sprintf("comp_update-%s.rrd", time.Now().String()))
		defer os.Remove(rrdFileName)

		cmd := exec.Command(rrdtool,
			"create",
			rrdFileName,
			"--start", "now",
			"--step", "1s",
			"DS:watts:GAUGE:5m:0:100000",
			"RRA:AVERAGE:0.5:5s:60m")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		So(cmd.Run(), ShouldBeNil)

		Convey("Then dump produces the same result", func() {
			cmd := exec.Command(rrdtool, "dump", rrdFileName)
			stdout, err := cmd.StdoutPipe()

			So(err, ShouldBeNil)

			cmd.Start()
			expectedResult, err := flattenXml(stdout)

			So(err, ShouldBeNil)

			pipeReader, pipeWriter := io.Pipe()
			go func() {
				flags := flag.NewFlagSet("gorrd", flag.ContinueOnError)
				flags.Parse([]string{rrdFileName})
				ctx := cli.NewContext(&cli.App{
					Writer: pipeWriter,
				}, flags, nil)
				dumpCommand(ctx)
				pipeWriter.Close()
			}()

			actualResult, err := flattenXml(pipeReader)

			So(err, ShouldBeNil)

			So(actualResult, ShouldResemble, expectedResult)
		})
	})
}
