package serde

import (
	"fmt"
	"reflect"

	"github.com/beevik/etree"
	"github.com/francistm/soap-server/internal"
	"github.com/francistm/soap-server/internal/model"
)

type defOpt func(o *definitionOption)

type definitionOption struct {
	location  string
	namespace string
}

func WithLocation(s string) defOpt {
	return func(o *definitionOption) {
		o.location = s
	}
}

func WithNamespace(s string) defOpt {
	return func(o *definitionOption) {
		o.namespace = s
	}
}

func BuildDefinitions(serviceName string, actions model.Actions, opts ...defOpt) *etree.Element {
	defOption := &definitionOption{
		location:  "",
		namespace: "http://example.org/",
	}

	for _, opt := range opts {
		opt(defOption)
	}

	defElem := etree.NewElement("definitions")
	defElem.CreateAttr("xmlns:s", internal.NsXml)
	defElem.CreateAttr("xmlns:soap", internal.NsSoap)
	defElem.CreateAttr("xmlns:wsdl", internal.NsWsdl)
	defElem.CreateAttr("xmlns:tns", defOption.namespace)

	typesElem := defElem.CreateElement("wsdl:types")
	schemaElem := typesElem.CreateElement("s:schema")
	schemaElem.CreateAttr("targetNamespace", defOption.namespace)

	serviceElem := etree.NewElement("wsdl:service")
	serviceElem.CreateAttr("name", serviceName)

	var sequenceElements []*etree.Element
	var messageElems []*etree.Element
	var portTypeElems []*etree.Element
	var bindingElems []*etree.Element

	for portName, port := range actions {
		wsdlPortName := fmt.Sprintf("%sSoap", portName)

		portTypeElem := etree.NewElement("wsdl:portType")
		portTypeElem.CreateAttr("name", wsdlPortName)

		bindingElem := etree.NewElement("wsdl:binding")
		bindingElem.CreateAttr("name", wsdlPortName)
		bindingElem.CreateAttr("type", fmt.Sprintf("tns:%s", wsdlPortName))
		soapBindingElem := etree.NewElement("soap:binding")
		soapBindingElem.CreateAttr("transport", internal.TransportHttp)

		portElem := serviceElem.CreateElement("wsdl:port")
		portElem.CreateAttr("name", portName)
		portElem.CreateAttr("binding", fmt.Sprintf("tns:%s", wsdlPortName))
		addressElem := portElem.CreateElement("soap:address")
		addressElem.CreateAttr("location", defOption.location)

		for actionName, action := range port {
			elemInName := portName + actionName
			elemOutName := portName + actionName + internal.ElemOutSuffix
			soapInTypeName := portName + actionName + internal.SoapInSuffix
			soapOutTypeName := portName + actionName + internal.SoapOutSuffix

			sequenceElements = buildDefSeqElems(elemInName, elemOutName, action.InType, action.OutType)
			messageElems = buildDefElems(elemInName, elemOutName, soapInTypeName, soapOutTypeName, action.InType, action.OutType)
			portOperationElem := buildPortOperation(actionName, soapInTypeName, soapOutTypeName, action.InType, action.OutType)
			bindingOperationElem := buildBindingOperation(defOption.namespace, portName, actionName, soapInTypeName, soapOutTypeName, action.InType, action.OutType)

			portTypeElem.AddChild(portOperationElem)
			bindingElem.AddChild(bindingOperationElem)
		}

		portTypeElems = append(portTypeElems, portTypeElem)
		bindingElems = append(bindingElems, bindingElem)
	}

	appendChildren(schemaElem, sequenceElements)
	appendChildren(defElem, messageElems)
	appendChildren(defElem, portTypeElems)
	appendChildren(defElem, bindingElems)
	defElem.AddChild(serviceElem)

	return defElem
}

func buildDefSeqElems(inName, outName string, tIn, tOut interface{}) []*etree.Element {
	elems := make([]*etree.Element, 0, 2)

	if tIn != nil {
		subElement := etree.NewElement("s:element")
		subElement.CreateAttr("name", inName)
		complexTypeElem := subElement.CreateElement("s:complexType")
		sequenceElem := complexTypeElem.CreateElement("s:sequence")

		elements := buildSeqElemsFromStruct(tIn)

		for _, el := range elements {
			sequenceElem.AddChild(el)
		}

		elems = append(elems, subElement)
	}

	if tOut != nil {
		subElement := etree.NewElement("s:element")
		subElement.CreateAttr("name", outName)
		complexTypeElem := subElement.CreateElement("s:complexType")
		sequenceElem := complexTypeElem.CreateElement("s:sequence")

		elements := buildSeqElemsFromStruct(tOut)

		for _, el := range elements {
			sequenceElem.AddChild(el)
		}

		elems = append(elems, subElement)
	}

	return elems
}

func buildDefElems(inName, outName, soapInName, soapOutName string, tIn, tOut interface{}) []*etree.Element {
	out := make([]*etree.Element, 0, 2)

	if tIn != nil {
		messageElem := etree.NewElement("wsdl:message")
		messageElem.CreateAttr("name", soapInName)
		partElem := messageElem.CreateElement("wsdl:part")
		partElem.CreateAttr("name", "parameters")
		partElem.CreateAttr("element", fmt.Sprintf("tns:%s", inName))

		out = append(out, messageElem)
	}

	if tOut != nil {
		messageElem := etree.NewElement("wsdl:message")
		messageElem.CreateAttr("name", soapOutName)
		partElem := messageElem.CreateElement("wsdl:part")
		partElem.CreateAttr("name", "parameters")
		partElem.CreateAttr("element", fmt.Sprintf("tns:%s", outName))

		out = append(out, messageElem)
	}

	return out
}

func buildPortOperation(operationName, soapInName, soapOutName string, tIn, tOut interface{}) *etree.Element {
	opElem := etree.NewElement("wsdl:operation")
	opElem.CreateAttr("name", operationName)

	if tIn != nil {
		el := opElem.CreateElement("wsdl:input")
		el.CreateAttr("message", soapInName)
	}

	if tOut != nil {
		el := opElem.CreateElement("wsdl:output")
		el.CreateAttr("message", soapOutName)
	}

	return opElem
}

func buildBindingOperation(tns, portName, actionName, soapInName, soapOutName string, tIn, tOut interface{}) *etree.Element {
	wsdlOperation := etree.NewElement("wsdl:operation")
	wsdlOperation.CreateAttr("name", actionName)

	soapOperation := wsdlOperation.CreateElement("soap:operation")
	soapOperation.CreateAttr("soapAction", fmt.Sprintf("%s%s%s", tns, portName, actionName))

	if tIn != nil {
		el := wsdlOperation.CreateElement("wsdl:input")
		body := el.CreateElement("soap:body")
		body.CreateAttr("use", "literal")

		wsdlOperation.AddChild(el)
	}

	if tOut != nil {
		el := wsdlOperation.CreateElement("wsdl:output")
		body := el.CreateElement("soap:body")
		body.CreateAttr("use", "literal")

		wsdlOperation.AddChild(el)
	}

	return wsdlOperation
}

func buildSeqElemsFromStruct(t interface{}) []*etree.Element {
	typeRef := reflect.TypeOf(t)

	if typeRef.Kind() == reflect.Ptr {
		typeRef = typeRef.Elem()
	}

	out := make([]*etree.Element, 0, typeRef.NumField())

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

		elem := etree.NewElement("s:element")

		var (
			attrType      string
			attrMinOccurs string
			attrMaxOccurs string
		)

		if isOptional {
			attrMinOccurs = "0"
		} else {
			attrMinOccurs = "1"
		}

		attrMaxOccurs = "1"

		switch fieldTypeRef.Kind() {
		case reflect.String:
			attrType = "s:string"

		case reflect.Int, reflect.Uint:
			attrType = "s:int"

		case reflect.Float32, reflect.Float64:
			attrType = "s:float"

		case reflect.Bool:
			attrType = "s:boolean"
		}

		elem.CreateAttr("name", fieldRef.Name)
		elem.CreateAttr("type", attrType)
		elem.CreateAttr("minOccurs", attrMinOccurs)
		elem.CreateAttr("maxOccurs", attrMaxOccurs)

		out = append(out, elem)
	}

	return out
}

func appendChildren(parent *etree.Element, children []*etree.Element) {
	for _, el := range children {
		parent.AddChild(el)
	}
}
