package serde

import (
	"reflect"

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

func ParseSoapIn(elem *etree.Element, inType any) (any, error) {
	entry, err := assignSoapInStructFields(elem, inType)

	if err != nil {
		return nil, errors.Wrap(err, "unable to create soapIn for request")
	}

	return entry, nil
}

func assignSoapInStructFields(entryElem *etree.Element, inType any) (any, error) {
	var (
		isPtr        = false
		entry        reflect.Value
		entryRefType = reflect.TypeOf(inType)
	)

	for entryRefType.Kind() == reflect.Ptr {
		isPtr = true
		entryRefType = entryRefType.Elem()
	}

	if entryRefType.Kind() != reflect.Struct {
		return nil, errors.Errorf("soapIn must be type of struct or *struct")
	}

	entry = reflect.New(entryRefType)

	for i := 0; i < entryRefType.NumField(); i++ {
		fieldRef := entryRefType.Field(i)
		entryFieldElems := entryElem.FindElements("./" + fieldRef.Name)

		if len(entryFieldElems) == 0 {
			continue
		}

		entryField := entry.Elem().Field(i)
		fieldTypeRef := fieldRef.Type

		if err := setFieldFromElements(entryField, fieldTypeRef, entryFieldElems); err != nil {
			return nil, err
		}
	}

	if !isPtr {
		entry = entry.Elem()
	}

	return entry.Interface(), nil
}

func setFieldFromElements(fieldValue reflect.Value, fieldType reflect.Type, elems []*etree.Element) error {
	switch fieldType.Kind() {
	case reflect.Ptr:
		elemTypeRef := fieldType.Elem()
		ptrValue := reflect.New(elemTypeRef)

		if err := setFieldFromElements(ptrValue.Elem(), elemTypeRef, elems); err != nil {
			return err
		}

		fieldValue.Set(ptrValue)

	case reflect.Slice:
		elemTypeRef, _ := derefType(fieldType.Elem())
		sliceValue := reflect.MakeSlice(fieldType, 0, len(elems))

		for _, fieldElem := range elems {
			parsedValue, err := parseScalarText(fieldElem.Text(), elemTypeRef)
			if err != nil {
				return err
			}

			if fieldType.Elem().Kind() == reflect.Ptr {
				ptrValue := reflect.New(elemTypeRef)
				ptrValue.Elem().Set(parsedValue)
				parsedValue = ptrValue
			}

			sliceValue = reflect.Append(sliceValue, parsedValue)
		}

		fieldValue.Set(sliceValue)

	case reflect.Array:
		if len(elems) != fieldValue.Len() {
			return errors.Errorf("field %s expects %d elements, got %d", fieldType, fieldValue.Len(), len(elems))
		}

		elemTypeRef, _ := derefType(fieldType.Elem())

		for i, fieldElem := range elems {
			parsedValue, err := parseScalarText(fieldElem.Text(), elemTypeRef)
			if err != nil {
				return err
			}

			if fieldType.Elem().Kind() == reflect.Ptr {
				ptrValue := reflect.New(elemTypeRef)
				ptrValue.Elem().Set(parsedValue)
				parsedValue = ptrValue
			}

			fieldValue.Index(i).Set(parsedValue)
		}

	default:
		parsedValue, err := parseScalarText(elems[0].Text(), fieldType)
		if err != nil {
			return err
		}

		fieldValue.Set(parsedValue)
	}

	return nil
}
