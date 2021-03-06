package dump

import (
	"encoding/xml"
	"io"
	"strconv"
	"time"

	"github.com/go-errors/errors"
	"github.com/untoldwind/gorrd/rrd"
)

type XmlDataOutput struct {
	encoder   *xml.Encoder
	tag       string
	lastError error
}

func NewXmlOutput(output io.Writer, prettyPrint bool) (*XmlDataOutput, error) {
	dumper := &XmlDataOutput{
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

func (d *XmlDataOutput) writeTokens(tokens []xml.Token) error {
	for _, token := range tokens {
		if err := d.encoder.EncodeToken(token); err != nil {
			return errors.Wrap(err, 0)
		}
	}
	return nil
}
func (d *XmlDataOutput) writeHeader() error {
	return d.writeTokens([]xml.Token{
		xml.ProcInst{Target: "xml", Inst: []byte(`version="1.0" encoding="utf-8"`)},
		xml.Directive(`DOCTYPE rrd SYSTEM "http://oss.oetiker.ch/rrdtool/rrdtool.dtd"`),
		xml.Comment(`Round Robin Database Dump`),
		xml.StartElement{Name: xml.Name{Local: d.tag}},
	})
}

func (d *XmlDataOutput) DumpComment(comment string) {
	if d.lastError != nil {
		return
	}
	if err := d.encoder.EncodeToken(xml.Comment(comment)); err != nil {
		d.lastError = errors.Wrap(err, 0)
	}
}

func (d *XmlDataOutput) DumpString(field, value string) {
	if d.lastError != nil {
		return
	}
	d.lastError = d.writeTokens([]xml.Token{
		xml.StartElement{Name: xml.Name{Local: field}},
		xml.CharData(value),
		xml.EndElement{Name: xml.Name{Local: field}},
	})
}

func (d *XmlDataOutput) DumpDouble(field string, value float64) {
	if d.lastError != nil {
		return
	}
	if err := d.writeTokens([]xml.Token{
		xml.StartElement{Name: xml.Name{Local: field}},
		xml.CharData(strconv.FormatFloat(value, 'e', 10, 64)),
		xml.EndElement{Name: xml.Name{Local: field}},
	}); err != nil {
		d.lastError = errors.Wrap(err, 0)
	}
}

func (d *XmlDataOutput) DumpUnsignedLong(field string, value uint64) {
	if d.lastError != nil {
		return
	}
	if err := d.writeTokens([]xml.Token{
		xml.StartElement{Name: xml.Name{Local: field}},
		xml.CharData(strconv.FormatUint(value, 10)),
		xml.EndElement{Name: xml.Name{Local: field}},
	}); err != nil {
		d.lastError = errors.Wrap(err, 0)
	}
}

func (d *XmlDataOutput) DumpTime(field string, value time.Time) {
	if d.lastError != nil {
		return
	}
	if err := d.writeTokens([]xml.Token{
		xml.StartElement{Name: xml.Name{Local: field}},
		xml.CharData(strconv.FormatInt(value.Unix(), 10)),
		xml.EndElement{Name: xml.Name{Local: field}},
		xml.Comment(value.String()),
	}); err != nil {
		d.lastError = errors.Wrap(err, 0)
	}
}

func (d *XmlDataOutput) DumpDuration(field string, value time.Duration) {
	if d.lastError != nil {
		return
	}
	if err := d.writeTokens([]xml.Token{
		xml.StartElement{Name: xml.Name{Local: field}},
		xml.CharData(strconv.FormatInt(int64(value.Seconds()), 10)),
		xml.EndElement{Name: xml.Name{Local: field}},
		xml.Comment(value.String()),
	}); err != nil {
		d.lastError = errors.Wrap(err, 0)
	}
}

func (d *XmlDataOutput) DumpSubFields(field string, subDump func(rrd.DataOutput) error) {
	if d.lastError != nil {
		return
	}
	dumper := &XmlDataOutput{
		encoder: d.encoder,
		tag:     field,
	}
	if err := d.encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: field}}); err != nil {
		d.lastError = err
		return
	}
	if err := subDump(dumper); err != nil {
		d.lastError = err
		return
	}
	d.lastError = dumper.Finalize()
}

func (d *XmlDataOutput) Finalize() error {
	if d.lastError != nil {
		return d.lastError
	}
	if err := d.encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: d.tag}}); err != nil {
		return errors.Wrap(err, 0)
	}
	return d.encoder.Flush()
}
