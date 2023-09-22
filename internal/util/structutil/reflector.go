package structutil

import (
	"encoding/json"
	"reflect"
	"strconv"
)

func ParseTag[T any]() []string {
	var t T
	res := []string{}

	baseType := reflect.TypeOf(t)
	for i := 0; i < baseType.NumField(); i++ {
		if col := baseType.Field(i).Tag.Get("db"); col != "" && col != "-" {
			res = append(res, col)
		}
	}

	return res
}

func StringToInt64(s string) (res int64) {
	res, _ = strconv.ParseInt(s, 10, 64)
	return
}

func StringToFloat32(s string) (res float32) {
	x, _ := strconv.ParseFloat(s, 32)
	return float32(x)
}

func PointerToVal[T any](val *T) (res T) {
	if val != nil {
		return *val
	}

	return res
}

func StructToMap[T any](val any) (res map[string]T) {
	res = make(map[string]T)
	b, _ := json.Marshal(val)
	json.Unmarshal(b, &res)

	return
}

func CheckMandatoryField(val any, fields ...string) (field string) {
	structType := reflect.TypeOf(val)
	if structType.Kind() != reflect.Struct {
		return
	}

	structValue := reflect.ValueOf(val)
	for _, v := range fields {
		f := structValue.FieldByName(v)
		if f.IsZero() {
			continue
		}

		if f.Type().Name() == "string" {
			val := f.Interface().(string)
			if val == "" {
				return v
			}
		}
	}

	return
}
