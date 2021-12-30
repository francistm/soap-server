package serde

import (
	"github.com/beevik/etree"
	"github.com/francistm/soap-server/internal"
)

func BuildEnvelope(bodyChildElem *etree.Element) *etree.Element {
	envelope := etree.NewElement("Envelope")
	envelope.Space = internal.XmlSoapNs
	envelope.CreateAttr("xmlns:soapenv", internal.NsSoap)

	bodyElem := envelope.CreateElement("Body")
	bodyElem.Space = internal.XmlSoapNs
	bodyElem.AddChild(bodyChildElem)

	return envelope
}

func BuildFaultBody(err error) *etree.Element {
	faultElem := etree.NewElement("Fault")
	faultElem.Space = internal.XmlSoapNs

	faultCodeElem := faultElem.CreateElement("Faultcode")
	faultCodeElem.SetText("")

	faultStringElem := faultElem.CreateElement("Faultstring")
	faultStringElem.SetText(err.Error())

	return faultElem
}
