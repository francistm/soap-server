package serde

import (
	"github.com/beevik/etree"
	"github.com/francistm/soap-server/internal/model"
)

const namespace = "soapenv"

func BuildEnvelope(bodyChildElem *etree.Element) *etree.Element {
	envelope := etree.NewElement("Envelope")
	envelope.Space = namespace
	envelope.CreateAttr("xmlns:soapenv", model.NsSoap)

	bodyElem := envelope.CreateElement("Body")
	bodyElem.Space = namespace
	bodyElem.AddChild(bodyChildElem)

	return envelope
}

func BuildFaultBody(err error) *etree.Element {
	faultElem := etree.NewElement("Fault")
	faultElem.Space = namespace

	faultCodeElem := faultElem.CreateElement("Faultcode")
	faultCodeElem.SetText("")

	faultStringElem := faultElem.CreateElement("Faultstring")
	faultStringElem.SetText(err.Error())

	return faultElem
}
