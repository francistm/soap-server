package parser

import (
	"reflect"
	"strconv"

	"github.com/beevik/etree"
	"github.com/pkg/errors"
)

func ExtractSoapInElem(doc *etree.Document) (*etree.Element, error) {
	bodyElems := doc.FindElements("//Body/*")

	if len(bodyElems) == 0 {
		return nil, errors.New("could not found entry in body")
	}

	if len(bodyElems) > 1 {
		return nil, errors.New("found multiple entries in body")
	}

	return bodyElems[0], nil
}

func ParseSoapIn(elem *etree.Element, inType interface{}) (interface{}, error) {
	entry, err := assignSoapInStructFields(elem, inType)

	if err != nil {
		return nil, errors.Wrap(err, "unable to create soapIn for request")
	}

	return entry, nil
}

func assignSoapInStructFields(entryElem *etree.Element, inType interface{}) (interface{}, error) {
	entryRefType := reflect.TypeOf(inType)

	for entryRefType.Kind() == reflect.Ptr {
		entryRefType = entryRefType.Elem()
	}

	if entryRefType.Kind() != reflect.Struct {
		return nil, errors.Errorf("soapIn must be type of struct or *struct")
	}

	entry := reflect.New(entryRefType)

	for i := 0; i < entryRefType.NumField(); i++ {
		fieldRef := entryRefType.Field(i)
		entryFieldElem := entryElem.FindElement("./" + fieldRef.Name)

		if entryFieldElem == nil {
			continue
		}

		entryFieldElemValue := entryFieldElem.Text()

		switch fieldRef.Type.Kind() {
		case reflect.String:
			entry.Elem().Field(i).SetString(entryFieldElemValue)

		case reflect.Bool:
			elemBoolValue, err := strconv.ParseBool(entryFieldElemValue)

			if err != nil {
				return nil, err
			}

			entry.Elem().Field(i).SetBool(elemBoolValue)

		case reflect.Float32, reflect.Float64:
			elemFloatValue, err := strconv.ParseFloat(entryFieldElemValue, 64)

			if err != nil {
				return nil, err
			}

			entry.Elem().Field(i).SetFloat(elemFloatValue)

		case reflect.Int, reflect.Uint:
			elemIntValue, err := strconv.ParseInt(entryFieldElemValue, 10, 64)

			if err != nil {
				return nil, err
			}

			entry.Elem().Field(i).SetInt(elemIntValue)
		}
	}

	return entry.Interface(), nil
}
