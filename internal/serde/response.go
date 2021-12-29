package serde

import (
	"github.com/beevik/etree"
	"github.com/francistm/soap-server/internal/model"
)

const namespace = "soapenv"

func BuildEnvelope(doc *etree.Document) *etree.Element {
	envelope := doc.CreateElement("Envelope")
	envelope.Space = namespace
	envelope.CreateAttr("xmlns:soapenv", model.NsSoap)

	bodyElem := envelope.CreateElement("Body")
	bodyElem.Space = namespace

	return bodyElem
}

func BuildFaultBody(bodyElem *etree.Element, err error) *etree.Element {
	faultElem := bodyElem.CreateElement("Fault")
	faultElem.Space = namespace

	faultCodeElem := faultElem.CreateElement("Faultcode")
	faultCodeElem.SetText("")

	faultStringElem := faultElem.CreateElement("Faultstring")
	faultStringElem.SetText(err.Error())

	return faultElem
}
