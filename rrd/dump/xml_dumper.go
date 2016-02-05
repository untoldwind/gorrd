package dump

import (
	"encoding/xml"
	"io"
	"strconv"
	"time"

	"github.com/untoldwind/gorrd/rrd"
)

type XmlDumber struct {
	encoder *xml.Encoder
	tag     string
}

func NewXmlDumper(output io.Writer, prettyPrint bool) (*XmlDumber, error) {
	dumper := &XmlDumber{
		encoder: xml.NewEncoder(output),
		tag:     "rrd",
	}
	if prettyPrint {
		dumper.encoder.Indent("", "  ")
	}

	if err := dumper.writeHeader(); err != nil {
		return nil, err
	}

	return dumper, nil
}

func (d *XmlDumber) writeTokens(tokens []xml.Token) error {
	for _, token := range tokens {
		if err := d.encoder.EncodeToken(token); err != nil {
			return err
		}
	}
	return nil
}
func (d *XmlDumber) writeHeader() error {
	return d.writeTokens([]xml.Token{
		xml.ProcInst{Target: "xml", Inst: []byte(`version="1.0" encoding="utf-8"`)},
		xml.Directive(`DOCTYPE rrd SYSTEM "http://oss.oetiker.ch/rrdtool/rrdtool.dtd"`),
		xml.Comment(`Round Robin Database Dump`),
		xml.StartElement{Name: xml.Name{Local: d.tag}},
	})
}

func (d *XmlDumber) DumpComment(comment string) error {
	return d.encoder.EncodeToken(xml.Comment(comment))
}

func (d *XmlDumber) DumpString(field, value string) error {
	return d.writeTokens([]xml.Token{
		xml.StartElement{Name: xml.Name{Local: field}},
		xml.CharData(value),
		xml.EndElement{Name: xml.Name{Local: field}},
	})
}

func (d *XmlDumber) DumpDouble(field string, value float64) error {
	return d.writeTokens([]xml.Token{
		xml.StartElement{Name: xml.Name{Local: field}},
		xml.CharData(strconv.FormatFloat(value, 'e', 10, 64)),
		xml.EndElement{Name: xml.Name{Local: field}},
	})
}

func (d *XmlDumber) DumpUnsignedLong(field string, value uint64) error {
	return d.writeTokens([]xml.Token{
		xml.StartElement{Name: xml.Name{Local: field}},
		xml.CharData(strconv.FormatUint(value, 10)),
		xml.EndElement{Name: xml.Name{Local: field}},
	})
}

func (d *XmlDumber) DumpTime(field string, value time.Time) error {
	return d.writeTokens([]xml.Token{
		xml.StartElement{Name: xml.Name{Local: field}},
		xml.CharData(strconv.FormatInt(value.Unix(), 10)),
		xml.EndElement{Name: xml.Name{Local: field}},
		xml.Comment(value.String()),
	})
}

func (d *XmlDumber) DumpSubFields(field string, subDump func(rrd.RrdDumper) error) error {
	dumper := &XmlDumber{
		encoder: d.encoder,
		tag:     field,
	}
	if err := d.encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: field}}); err != nil {
		return err
	}
	if err := subDump(dumper); err != nil {
		return err
	}
	return dumper.Finalize()
}

func (d *XmlDumber) Finalize() error {
	if err := d.encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: d.tag}}); err != nil {
		return err
	}
	return d.encoder.Flush()
}
