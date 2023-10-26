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

func CheckMandatoryField(val any) (field string) {
	structType := reflect.TypeOf(val)
	if structType.Kind() != reflect.Ptr || structType.Elem().Kind() != reflect.Struct {
		return
	}

	structValue := reflect.ValueOf(val).Elem()
	for i := 0; i < structValue.NumField(); i++ {
		f := structValue.Field(i)

		if !f.IsValid() {
			continue
		}

		fieldName := structType.Elem().Field(i).Name
		fieldTag := structType.Elem().Field(i).Tag
		if val, ok := fieldTag.Lookup("validate"); val != "required" || !ok {
			continue
		}

		if f.Kind() == reflect.Ptr && f.IsNil() {
			return fieldName
		}

		switch f.Type().Kind() {
		case reflect.String:
			val := f.Interface().(string)
			if val == "" {
				return fieldName
			}
		case reflect.Float32:
			val := f.Interface().(float64)
			if val == 0 {
				return fieldName
			}
		case reflect.Float64:
			val := f.Interface().(float64)
			if val == 0 {
				return fieldName
			}
		case reflect.Int64:
			val := f.Interface().(int64)
			if val == 0 {
				return fieldName
			}
		case reflect.Uint64:
			val := f.Interface().(uint64)
			if val == 0 {
				return fieldName
			}
		}
	}

	return
}
