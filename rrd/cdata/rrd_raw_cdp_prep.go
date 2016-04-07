package cdata

import (
	"reflect"
	"strconv"

	"github.com/go-errors/errors"
)

type rrdCdpPrep struct {
	scratch []unival
}

const rrdRawCdpPrepSize = 10 * 8

func (f *RrdRawFile) readCdpPreps(reader *CDataReader) error {
	f.cdpPreps = make([][]*rrdCdpPrep, f.header.rraCount)

	var err error
	for i := range f.cdpPreps {
		f.cdpPreps[i] = make([]*rrdCdpPrep, f.header.datasourceCount)
		for j := range f.cdpPreps[i] {
			f.cdpPreps[i][j], err = f.readCdpPrep(reader)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *RrdRawFile) readCdpPrep(reader *CDataReader) (*rrdCdpPrep, error) {
	scratch, err := reader.ReadUnivals(10)
	if err != nil {
		return nil, err
	}
	return &rrdCdpPrep{
		scratch: scratch,
	}, nil
}

func (f *RrdRawFile) storeCdpPreps() error {
	writer := f.dataFile.Writer(f.baseHeaderSize + rrdRawLiveHeaderSize +
		rrdRawPdpPrepSize*f.header.datasourceCount)

	for _, cdpPreps := range f.cdpPreps {
		for _, cdpPrep := range cdpPreps {
			if err := storeCdpPrep(writer, cdpPrep); err != nil {
				return err
			}
		}
	}
	return nil
}

func storeCdpPrep(writer *CDataWriter, cdpPrep *rrdCdpPrep) error {
	if err := writer.WriteUnivals(cdpPrep.scratch); err != nil {
		return err
	}
	return nil
}

func (f *RrdRawFile) decodeRraCpdPreps(rraIndex, dsIndex int, rv reflect.Value) error {
	if rv.Kind() != reflect.Struct {
		return errors.Errorf("Cdp params must be a struct")
	}
	for i := 0; i < rv.Type().NumField(); i++ {
		field := rv.Type().Field(i)

		if field.Type.Kind() == reflect.Struct {
			if err := f.decodeRraCpdPreps(rraIndex, dsIndex, rv.Field(i)); err != nil {
				return err
			}
			continue
		}

		scratch := f.cdpPreps[rraIndex][dsIndex].scratch
		tag := field.Tag.Get("cdp")
		if tag == "" {
			continue
		}

		if tag == "raw" && field.Type.Kind() == reflect.Slice {
			switch field.Type.Elem().Kind() {
			case reflect.Uint8:
				convered := f.dataFile.UnivalsToBytes(scratch)
				rv.Field(i).Set(reflect.ValueOf(convered))
			default:
				return errors.Errorf("cpd raw must have type []byte")
			}
		} else {
			scratchIndex, err := strconv.ParseInt(tag, 10, 64)
			if err != nil {
				return err
			}
			switch field.Type.Kind() {
			case reflect.Uint64:
				rv.Field(i).SetUint(scratch[scratchIndex].AsUnsignedLong())
			case reflect.Float64:
				rv.Field(i).SetFloat(scratch[scratchIndex].AsDouble())
			default:
				return errors.Errorf("cpd field must have type uint64 or float64")
			}
		}
	}
	return nil
}

func (f *RrdRawFile) encodeRraCpdPreps(rraIndex, dsIndex int, rv reflect.Value) error {
	if rv.Kind() != reflect.Struct {
		return errors.Errorf("Cdp params must be a struct")
	}
	for i := 0; i < rv.Type().NumField(); i++ {
		field := rv.Type().Field(i)

		if field.Type.Kind() == reflect.Struct {
			if err := f.encodeRraCpdPreps(rraIndex, dsIndex, rv.Field(i)); err != nil {
				return err
			}
			continue
		}

		scratch := f.cdpPreps[rraIndex][dsIndex].scratch
		tag := field.Tag.Get("cdp")
		if tag == "" {
			continue
		}

		if tag == "raw" && field.Type.Kind() == reflect.Slice {
			switch field.Type.Elem().Kind() {
			case reflect.Uint8:
				raw := rv.Field(i).Interface().([]uint8)
				values := f.dataFile.BytesToUnivals(raw)
				for i, v := range values {
					scratch[i] = v
				}
			default:
				return errors.Errorf("cpd scatch must have type []uint64")
			}
		} else {
			scratchIndex, err := strconv.ParseInt(tag, 10, 64)
			if err != nil {
				return err
			}
			switch field.Type.Kind() {
			case reflect.Uint64:
				scratch[scratchIndex] = univalForUnsignedLong(rv.Field(i).Uint())
			case reflect.Float64:
				scratch[scratchIndex] = univalForDouble(rv.Field(i).Float())
			default:
				return errors.Errorf("cpd field must have type uint64 or float64")
			}
		}
	}
	return nil
}
