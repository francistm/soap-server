package serde

import (
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

func derefType(t reflect.Type) (reflect.Type, bool) {
	isPtr := false

	for t.Kind() == reflect.Ptr {
		isPtr = true
		t = t.Elem()
	}

	return t, isPtr
}

func scalarXMLType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "s:string"

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "s:int"

	case reflect.Float32, reflect.Float64:
		return "s:float"

	case reflect.Bool:
		return "s:boolean"
	}

	return ""
}

func scalarText(v reflect.Value) (string, bool) {
	switch v.Kind() {
	case reflect.String:
		return v.String(), true

	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), true

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), true

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), true

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'E', -1, 64), true
	}

	return "", false
}

func parseScalarText(s string, t reflect.Type) (reflect.Value, error) {
	v := reflect.New(t).Elem()

	switch t.Kind() {
	case reflect.String:
		v.SetString(s)

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(s)
		if err != nil {
			return reflect.Value{}, err
		}

		v.SetBool(boolValue)

	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(s, t.Bits())
		if err != nil {
			return reflect.Value{}, err
		}

		v.SetFloat(floatValue)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(s, 10, t.Bits())
		if err != nil {
			return reflect.Value{}, err
		}

		v.SetInt(intValue)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(s, 10, t.Bits())
		if err != nil {
			return reflect.Value{}, err
		}

		v.SetUint(uintValue)

	default:
		return reflect.Value{}, errors.Errorf("unsupported scalar type %s", t)
	}

	return v, nil
}
