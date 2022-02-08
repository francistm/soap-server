package serde

import (
	"reflect"
	"strconv"

	"github.com/beevik/etree"
)

func BuildEnvelope(space string, ns map[string]string, bodyChildElem *etree.Element) *etree.Element {
	envelope := etree.NewElement("Envelope")
	envelope.Space = space

	for key, value := range ns {
		envelope.CreateAttr(key, value)
	}

	headerElem := envelope.CreateElement("Header")
	headerElem.Space = space

	bodyElem := envelope.CreateElement("Body")
	bodyElem.Space = space

	bodyElem.AddChild(bodyChildElem)

	return envelope
}

func BuildFaultBodyChild(space string, err error) *etree.Element {
	faultElem := etree.NewElement("Fault")
	faultElem.Space = space

	faultCodeElem := faultElem.CreateElement("Faultcode")
	faultCodeElem.SetText("")

	faultStringElem := faultElem.CreateElement("Faultstring")
	faultStringElem.SetText(err.Error())

	return faultElem
}

func BuildResponseBodyChild(ns, soapOutTagName string, tOut, out interface{}) *etree.Element {
	outElem := etree.NewElement(soapOutTagName)
	outElem.CreateAttr("xmlns", ns)

	fieldElements := buildOutStructFieldsElements(tOut, out)

	for _, elem := range fieldElements {
		outElem.AddChild(elem)
	}

	return outElem
}

func IsActionReturnValid(t, v interface{}) bool {
	if t == nil {
		return true
	}

	if v == nil {
		return false
	}

	tTypeRef := reflect.TypeOf(t)

	if tTypeRef.Kind() == reflect.Ptr {
		tTypeRef = tTypeRef.Elem()
	}

	vTypeRef := reflect.TypeOf(v)

	if vTypeRef.Kind() == reflect.Ptr {
		vTypeRef = vTypeRef.Elem()
	}

	return tTypeRef.PkgPath() == vTypeRef.PkgPath() && tTypeRef.Name() == vTypeRef.Name()
}

func buildOutStructFieldsElements(tOut, out interface{}) []*etree.Element {
	if tOut == nil || out == nil {
		return nil
	}

	tOutTypeRef := reflect.TypeOf(tOut)
	outValueRef := reflect.ValueOf(out)

	if tOutTypeRef.Kind() == reflect.Ptr {
		tOutTypeRef = tOutTypeRef.Elem()
	}

	if outValueRef.Kind() == reflect.Ptr {
		outValueRef = outValueRef.Elem()
	}

	elems := make([]*etree.Element, 0, tOutTypeRef.NumField())

	for i := 0; i < tOutTypeRef.NumField(); i++ {
		outField := tOutTypeRef.Field(i)

		if len(outField.PkgPath) > 0 {
			continue
		}

		var elemText string
		elem := etree.NewElement(outField.Name)
		fieldValueRef := outValueRef.Field(i)

		if fieldValueRef.Kind() == reflect.Ptr {
			if fieldValueRef.IsNil() {
				elem.SetText("")
				continue
			}

			fieldValueRef = fieldValueRef.Elem()
		}

		switch fieldValueRef.Kind() {
		case reflect.String:
			elemText = fieldValueRef.String()

		case reflect.Bool:
			elemText = strconv.FormatBool(fieldValueRef.Bool())

		case reflect.Int, reflect.Uint:
			intVal := fieldValueRef.Int()
			elemText = strconv.Itoa(int(intVal))

		case reflect.Float32, reflect.Float64:
			floatVal := fieldValueRef.Float()
			elemText = strconv.FormatFloat(floatVal, 'E', -1, 64)
		}

		elem.SetText(elemText)
		elems = append(elems, elem)
	}

	return elems
}
