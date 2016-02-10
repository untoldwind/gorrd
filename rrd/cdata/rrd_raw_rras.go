package cdata

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
)

type rrdRawRraDef struct {
	rraType    string
	rowCount   uint64
	pdpPerRow  uint64
	parameters []unival
}

func (f *RrdRawFile) RraTypes() []string {
	result := make([]string, len(f.rraDefs))
	for i, rrdDef := range f.rraDefs {
		result[i] = rrdDef.rraType
	}
	return result
}

func (f *RrdRawFile) ReadRraParams(index int, params interface{}) error {
	rv := reflect.ValueOf(params)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.Errorf("Rra params must be a pointer")
	}
	return f.decodeRraParams(index, rv.Elem())
}

func (f *RrdRawFile) decodeRraParams(index int, rv reflect.Value) error {
	if rv.Kind() != reflect.Struct {
		return errors.Errorf("Rra params must be a pointer to a struct")
	}
	for i := 0; i < rv.Type().NumField(); i++ {
		field := rv.Type().Field(i)
		tag := field.Tag.Get("rra")
		switch {
		case tag == "rowCount":
			if field.Type.Kind() != reflect.Uint64 {
				return errors.Errorf("rra rowCount has to be an uint64")
			}
			rv.Field(i).SetUint(f.rraDefs[index].rowCount)
		case tag == "pdpPerRow":
			if field.Type.Kind() != reflect.Uint64 {
				return errors.Errorf("rra pdpPerRow has to be an uint64")
			}
			rv.Field(i).SetUint(f.rraDefs[index].pdpPerRow)
		case strings.HasPrefix(tag, "param"):
			paramIndex, err := strconv.ParseInt(tag[5:], 10, 64)
			if err != nil {
				return errors.Errorf("rra param has invalid index: %s", err.Error())
			}
			switch field.Type.Kind() {
			case reflect.Uint64:
				rv.Field(i).SetUint(f.rraDefs[index].parameters[paramIndex].AsUnsignedLong())
			case reflect.Float64:
				rv.Field(i).SetFloat(f.rraDefs[index].parameters[paramIndex].AsDouble())
			default:
				return errors.Errorf("param must have type uint64 or float64")
			}
		case tag == "cpdPreps":
			if field.Type.Kind() != reflect.Slice {
				return errors.Errorf("rra cpdPreps has to be a slice of structs")
			}
			cpdPreps := reflect.MakeSlice(field.Type, int(f.header.datasourceCount), int(f.header.datasourceCount))
			rv.Field(i).Set(cpdPreps)
			if field.Type.Elem().Kind() == reflect.Struct {
				for dsIndex := 0; dsIndex < cpdPreps.Len(); dsIndex++ {
					if err := f.decodeRraCpdPreps(index, dsIndex, cpdPreps.Index(dsIndex)); err != nil {
						return err
					}
				}
			}
		case tag == "":
			if field.Type.Kind() == reflect.Struct {
				if err := f.decodeRraParams(index, rv.Field(i)); err != nil {
					return err
				}
			}
		default:
			return errors.Errorf("Unknown rra tag: %s", tag)
		}
	}
	return nil
}

func (f *RrdRawFile) StoreRraParams(index int, params interface{}) error {
	rv := reflect.ValueOf(params)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.Errorf("Rra params must be a pointer")
	}
	return f.encodeRraParams(index, rv.Elem())
}

func (f *RrdRawFile) encodeRraParams(index int, rv reflect.Value) error {
	if rv.Kind() != reflect.Struct {
		return errors.Errorf("Rra params must be a pointer to a struct")
	}
	for i := 0; i < rv.Type().NumField(); i++ {
		field := rv.Type().Field(i)
		tag := field.Tag.Get("rra")
		switch {
		case tag == "rowCount":
			if field.Type.Kind() != reflect.Uint64 {
				return errors.Errorf("rra rowCount has to be an uint64")
			}
			f.rraDefs[index].rowCount = rv.Field(i).Uint()
		case tag == "pdpPerRow":
			if field.Type.Kind() != reflect.Uint64 {
				return errors.Errorf("rra pdpPerRow has to be an uint64")
			}
			f.rraDefs[index].pdpPerRow = rv.Field(i).Uint()
		case strings.HasPrefix(tag, "param"):
			paramIndex, err := strconv.ParseInt(tag[5:], 10, 64)
			if err != nil {
				return errors.Errorf("rra param has invalid index: %s", err.Error())
			}
			switch field.Type.Kind() {
			case reflect.Uint64:
				f.rraDefs[index].parameters[paramIndex] = univalForUnsignedLong(rv.Field(i).Uint())
			case reflect.Float64:
				f.rraDefs[index].parameters[paramIndex] = univalForDouble(rv.Field(i).Float())
			default:
				return errors.Errorf("param must have type uint64 or float64")
			}
		case tag == "cpdPreps":
			if field.Type.Kind() != reflect.Slice {
				return errors.Errorf("rra cpdPreps has to be a slice of structs")
			}
			cpdPreps := rv.Field(i)
			if field.Type.Elem().Kind() == reflect.Struct {
				for dsIndex := 0; dsIndex < cpdPreps.Len(); dsIndex++ {
					if err := f.encodeRraCpdPreps(index, dsIndex, cpdPreps.Index(dsIndex)); err != nil {
						return err
					}
				}
			}
		case tag == "":
			if field.Type.Kind() == reflect.Struct {
				if err := f.encodeRraParams(index, rv.Field(i)); err != nil {
					return err
				}
			}
		default:
			return errors.Errorf("Unknown rra tag: %s", tag)
		}
	}
	return nil
}

func (f *RrdRawFile) decodeRraCpdPreps(rraIndex, dsIndex int, rv reflect.Value) error {
	if rv.Kind() != reflect.Struct {
		return errors.Errorf("Cdp params must be a struct")
	}
	for i := 0; i < rv.Type().NumField(); i++ {
		field := rv.Type().Field(i)
		tag := field.Tag.Get("cdp")
		if tag == "" {
			continue
		}
		scratchIndex, err := strconv.ParseInt(tag, 10, 64)
		if err != nil {
			return err
		}
		switch field.Type.Kind() {
		case reflect.Uint64:
			rv.Field(i).SetUint(f.cdpPreps[rraIndex][dsIndex].scratch[scratchIndex].AsUnsignedLong())
		case reflect.Float64:
			rv.Field(i).SetFloat(f.cdpPreps[rraIndex][dsIndex].scratch[scratchIndex].AsDouble())
		default:
			return errors.Errorf("cpd field must have type uint64 or float64")
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
		tag := field.Tag.Get("cdp")
		if tag == "" {
			continue
		}
		scratchIndex, err := strconv.ParseInt(tag, 10, 64)
		if err != nil {
			return err
		}
		switch field.Type.Kind() {
		case reflect.Uint64:
			f.cdpPreps[rraIndex][dsIndex].scratch[scratchIndex] = univalForUnsignedLong(rv.Field(i).Uint())
		case reflect.Float64:
			f.cdpPreps[rraIndex][dsIndex].scratch[scratchIndex] = univalForDouble(rv.Field(i).Float())
		default:
			return errors.Errorf("cpd field must have type uint64 or float64")
		}

	}
	return nil
}

func (f *RrdRawFile) readRras(reader *CDataReader) error {
	f.rraDefs = make([]*rrdRawRraDef, f.header.rraCount)

	var err error
	for i := range f.rraDefs {
		f.rraDefs[i], err = readRra(reader, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func readRra(reader *CDataReader, index int) (*rrdRawRraDef, error) {
	rraType, err := reader.ReadCString(20)
	if err != nil {
		return nil, err
	}
	rowCount, err := reader.ReadUnsignedLong()
	if err != nil {
		return nil, err
	}
	pdpPerRow, err := reader.ReadUnsignedLong()
	if err != nil {
		return nil, err
	}
	parameters, err := reader.ReadUnivals(10)
	if err != nil {
		return nil, err
	}
	return &rrdRawRraDef{
		rraType:    rraType,
		rowCount:   rowCount,
		pdpPerRow:  pdpPerRow,
		parameters: parameters,
	}, nil
}
