package model

import "reflect"

func StructToSeqElems(i interface{}) []*SchemaElementComplexTypeSeqElement {
	typeRef := reflect.TypeOf(i)

	if typeRef.Kind() == reflect.Ptr {
		typeRef = typeRef.Elem()
	}

	out := make([]*SchemaElementComplexTypeSeqElement, 0, typeRef.NumField())

	for i := 0; i < typeRef.NumField(); i++ {
		isOptional := false
		fieldRef := typeRef.Field(i)

		if len(fieldRef.PkgPath) > 0 {
			continue
		}

		fieldTypeRef := fieldRef.Type

		if fieldTypeRef.Kind() == reflect.Ptr {
			isOptional = true
			fieldTypeRef = fieldTypeRef.Elem()
		}

		elem := &SchemaElementComplexTypeSeqElement{
			Name:      fieldRef.Name,
			MinOccurs: 1,
			MaxOccurs: 1,
		}

		if isOptional {
			elem.MinOccurs = 0
		}

		switch fieldTypeRef.Kind() {
		case reflect.String:
			elem.Type = "s:string"

		case reflect.Int, reflect.Uint:
			elem.Type = "s:int"

		case reflect.Float32, reflect.Float64:
			elem.Type = "s:float"

		case reflect.Bool:
			elem.Type = "s:boolean"
		}

		out = append(out, elem)
	}

	return out
}
