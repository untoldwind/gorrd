package cdata

import (
	"github.com/go-errors/errors"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type rrdRawDatasourceDef struct {
	name           string
	dataSourceType string
	parameters     []unival
}

func (f *RrdRawFile) DatasourceTypes() []string {
	result := make([]string, len(f.datasourceDefs))
	for i, datasourceDef := range f.datasourceDefs {
		result[i] = datasourceDef.dataSourceType
	}
	return result
}

func (f *RrdRawFile) ReadDatasourceParams(index int, params interface{}) error {
	rv := reflect.ValueOf(params)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.Errorf("Datasource params must be a pointer")
	}
	return f.decodeDatasourceParams(index, rv.Elem())
}

func (f *RrdRawFile) decodeDatasourceParams(index int, rv reflect.Value) error {
	if rv.Kind() != reflect.Struct {
		return errors.Errorf("Datasource params must be a pointer to a struct")
	}
	for i := 0; i < rv.Type().NumField(); i++ {
		field := rv.Type().Field(i)
		dsTag := field.Tag.Get("ds")
		pdpTag := field.Tag.Get("pdp")
		switch {
		case dsTag == "name":
			if field.Type.Kind() != reflect.String {
				return errors.Errorf("name field has to be a string")
			}
			rv.Field(i).SetString(f.datasourceDefs[index].name)
		case strings.HasPrefix(dsTag, "param"):
			paramIndex, err := strconv.ParseInt(dsTag[5:], 10, 64)
			if err != nil {
				return errors.Errorf("datasource param has invalid index: %s", err.Error())
			}
			switch field.Type.Kind() {
			case reflect.Uint64:
				rv.Field(i).SetUint(f.datasourceDefs[index].parameters[paramIndex].AsUnsignedLong())
			case reflect.Float64:
				rv.Field(i).SetFloat(f.datasourceDefs[index].parameters[paramIndex].AsDouble())
			default:
				return errors.Errorf("param must have type uint64 or float64")
			}
		case pdpTag == "lastValue":
			if field.Type.Kind() != reflect.String {
				return errors.Errorf("lastValue field has to be a string")
			}
			rv.Field(i).SetString(f.pdpPreps[index].lastDatasourceValue)
		case len(pdpTag) == 1 && unicode.IsDigit(rune(pdpTag[0])):
			scratchIndex, err := strconv.ParseInt(pdpTag, 10, 64)
			if err != nil {
				return errors.Errorf("datasource pdp has invalid index: %s", err.Error())
			}
			switch field.Type.Kind() {
			case reflect.Uint64:
				rv.Field(i).SetUint(f.pdpPreps[index].scratch[scratchIndex].AsUnsignedLong())
			case reflect.Float64:
				rv.Field(i).SetFloat(f.pdpPreps[index].scratch[scratchIndex].AsDouble())
			default:
				return errors.Errorf("datasource pdp must have type uint64 or float64")
			}
		case dsTag == "" && pdpTag == "":
			if field.Type.Kind() == reflect.Struct {
				if err := f.decodeDatasourceParams(index, rv.Field(i)); err != nil {
					return err
				}
			}
		default:
			return errors.Errorf("Unknown datasource tag: %s", field.Tag)
		}
	}
	return nil
}

func (f *RrdRawFile) readDatasources() error {
	f.datasourceDefs = make([]*rrdRawDatasourceDef, f.header.datasourceCount)

	var err error
	for i := range f.datasourceDefs {
		f.datasourceDefs[i], err = readDatasource(f.dataFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func readDatasource(dataFile *CDataFile) (*rrdRawDatasourceDef, error) {
	name, err := dataFile.ReadCString(20)
	if err != nil {
		return nil, err
	}
	datasourceType, err := dataFile.ReadCString(20)
	if err != nil {
		return nil, err
	}
	parameters, err := dataFile.ReadUnivals(10)
	if err != nil {
		return nil, err
	}
	return &rrdRawDatasourceDef{
		name:           name,
		dataSourceType: datasourceType,
		parameters:     parameters,
	}, nil
}
