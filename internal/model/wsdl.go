package model

import "encoding/xml"

type Definitions struct {
	XMLName xml.Name `xml:"wsdl:definitions"`

	XMLNsS    string `xml:"xmlns:s,attr"`
	XMLNsSOAP string `xml:"xmlns:soap,attr"`
	XMLNsWSDL string `xml:"xmlns:wsdl,attr"`
	XMLNsTns  string `xml:"xmlns:tns,attr"`

	TNS string `xml:"targetNamespace,attr"`

	Types     *Types      `xml:"wsdl:types"`
	Messages  []*Message  `xml:"wsdl:message"`
	PortTypes []*PortType `xml:"wsdl:portType"`
	Bindings  []*Binding  `xml:"wsdl:binding"`
	Service   *Service    `xml:"wsdl:service"`
}

type Types struct {
	Schema *Schema `xml:"s:schema"`
}

type Schema struct {
	TNS      string           `xml:"targetNamespace,attr"`
	Elements []*SchemaElement `xml:"s:element"`
}

type SchemaElement struct {
	Name string                    `xml:"name,attr"`
	Type *SchemaElementComplexType `xml:"s:complexType"`
}

type SchemaElementComplexType struct {
	Sequence *SchemaElementComplexTypeSequence `xml:"s:sequence"`
}

type SchemaElementComplexTypeSequence struct {
	Elements []*SchemaElementComplexTypeSeqElement `xml:"s:element"`
}

type SchemaElementComplexTypeSeqElement struct {
	Name      string `xml:"name,attr"`
	Type      string `xml:"type,attr"`
	MinOccurs int    `xml:"minOccurs,attr"`
	MaxOccurs int    `xml:"maxOccurs,attr"`
}

type Message struct {
	Name string       `xml:"name,attr"`
	Part *MessagePart `xml:"wsdl:part"`
}

type MessagePart struct {
	Name     string `xml:"name,attr"`
	ElemName string `xml:"element,attr"`
}

type PortType struct {
	Name       string               `xml:"name,attr"`
	Operations []*PortTypeOperation `xml:"wsdl:operation"`
}

type PortTypeOperation struct {
	Name   string                 `xml:"name,attr"`
	Doc    *PortTypeOperationDoc  `xml:"wsdl:documentation"`
	Input  *PortTypeOperationItem `xml:"wsdl:input,omitempty"`
	Output *PortTypeOperationItem `xml:"wsdl:output,omitempty"`
}

type PortTypeOperationDoc struct {
	NsWsdl string `xml:"xmlns:wsdl,attr"` // should be http://schemas.xmlsoap.org/wsdl/
	Body   string `xml:",chardata"`
}

type PortTypeOperationItem struct {
	Message string `xml:"message,attr"`
}

type Binding struct {
	Name       string                  `xml:"name,attr"`
	Type       string                  `xml:"type,attr"`
	Binding    *BindingSoap            `xml:"soap:binding"`
	Operations []*BindingWsdlOperation `xml:"wsdl:operation"`
}

type BindingSoap struct {
	Transport string `xml:"transport,attr"`
}

type BindingWsdlOperation struct {
	Name      string                    `xml:"name,attr"`
	Operation *BindingSoapOperation     `xml:"soap:operation"`
	Input     *BindingWsdlOperationItem `xml:"wsdl:input,omitempty"`
	Output    *BindingWsdlOperationItem `xml:"wsdl:output,omitempty"`
}

type BindingSoapOperation struct {
	Action string `xml:"soapAction,attr"`
	Style  string `xml:"style,attr"`
}

type BindingWsdlOperationItem struct {
	Body *BindingWsdlOperationItemBody `xml:"soap:body"`
}

type BindingWsdlOperationItemBody struct {
	Use string `xml:"use,attr"` // should be literal
}

type Service struct {
	Name  string         `xml:"name,attr"`
	Ports []*ServicePort `xml:"wsdl:port"`
}

type ServicePort struct {
	Name     string           `xml:"name,attr"`
	Binding  string           `xml:"binding,attr"`
	Location *ServicePortAddr `xml:"soap:address"`
}

type ServicePortAddr struct {
	Location string `xml:"location,attr"`
}
