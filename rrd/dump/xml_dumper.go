package dump

import (
	"encoding/xml"
	"io"
)

type XmlDumber struct {
	encoder *xml.Encoder
}

func NewXmlDumper(output io.Writer) (*XmlDumber, error) {
	dumper := &XmlDumber{
		encoder: xml.NewEncoder(output),
	}

	if err := dumper.writeHeader(); err != nil {
		return nil, err
	}

	return dumper, nil
}

func (d *XmlDumber) writeHeader() error {
	tokens := []xml.Token{
		xml.ProcInst{Target: "xml", Inst: []byte(`version="1.0" encoding="utf-8"`)},
		xml.StartElement{Name: xml.Name{Local: "rrd"}},
	}
	for _, token := range tokens {
		if err := d.encoder.EncodeToken(token); err != nil {
			return err
		}
	}
	return nil
}

func (d *XmlDumber) Finalize() error {
	if err := d.encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "rrd"}}); err != nil {
		return err
	}
	return d.encoder.Flush()
}
