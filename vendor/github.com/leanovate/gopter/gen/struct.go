package gen

import (
	"reflect"

	"github.com/leanovate/gopter"
)

func StructPtr(rt reflect.Type, gens map[string]gopter.Gen) gopter.Gen {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return Fail(rt)
	}
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		result := reflect.New(rt)

		for name, gen := range gens {
			field, ok := rt.FieldByName(name)
			if !ok {
				continue
			}
			value, ok := gen(genParams).Retrieve()
			if !ok {
				return gopter.NewEmptyResult(rt)
			}
			result.Elem().FieldByIndex(field.Index).Set(reflect.ValueOf(value))
		}

		return gopter.NewGenResult(result.Interface(), gopter.NoShrinker)
	}
}
