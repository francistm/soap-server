package serde

import (
	"testing"

	"github.com/beevik/etree"
	"github.com/francistm/soap-server/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type definitionField struct {
	name      string
	fieldType string
	minOccurs string
	maxOccurs string
}

func TestBuildDefinitionsSingleIntArg(t *testing.T) {
	type InType struct {
		I int
	}

	got := buildDefinitionTestDoc(InType{}, nil)
	assertDefinitionFields(t, got, "Port1Action1", []definitionField{
		{name: "I", fieldType: "s:int", minOccurs: "1", maxOccurs: "1"},
	})
	assertDefinitionElementMissing(t, got, "Port1Action1Response")
}

func TestBuildDefinitionsSingleStringArg(t *testing.T) {
	type InType struct {
		S string
	}

	got := buildDefinitionTestDoc(InType{}, nil)
	assertDefinitionFields(t, got, "Port1Action1", []definitionField{
		{name: "S", fieldType: "s:string", minOccurs: "1", maxOccurs: "1"},
	})
	assertDefinitionElementMissing(t, got, "Port1Action1Response")
}

func TestBuildDefinitionsSingleBoolArg(t *testing.T) {
	type InType struct {
		B bool
	}

	got := buildDefinitionTestDoc(InType{}, nil)
	assertDefinitionFields(t, got, "Port1Action1", []definitionField{
		{name: "B", fieldType: "s:boolean", minOccurs: "1", maxOccurs: "1"},
	})
	assertDefinitionElementMissing(t, got, "Port1Action1Response")
}

func TestBuildDefinitionsSingleOptionalStringArg(t *testing.T) {
	type InType struct {
		SP *string
	}

	got := buildDefinitionTestDoc(InType{}, nil)
	assertDefinitionFields(t, got, "Port1Action1", []definitionField{
		{name: "SP", fieldType: "s:string", minOccurs: "0", maxOccurs: "1"},
	})
	assertDefinitionElementMissing(t, got, "Port1Action1Response")
}

func TestBuildDefinitionsSingleSliceArg(t *testing.T) {
	type InType struct {
		Ints []int
	}

	got := buildDefinitionTestDoc(InType{}, nil)
	assertDefinitionFields(t, got, "Port1Action1", []definitionField{
		{name: "Ints", fieldType: "s:int", minOccurs: "0", maxOccurs: "unbounded"},
	})
	assertDefinitionElementMissing(t, got, "Port1Action1Response")
}

func TestBuildDefinitionsSkipsUnsupportedArg(t *testing.T) {
	type InType struct {
		Items []struct{}
	}

	got := buildDefinitionTestDoc(InType{}, nil)
	assertDefinitionFields(t, got, "Port1Action1", nil)
}

func TestBuildDefinitionsMultipleArgs(t *testing.T) {
	type InType struct {
		S    string
		I    int
		B    bool
		SP   *string
		Ints []int
	}

	type OutType struct {
		S string
		I int
	}

	got := buildDefinitionTestDoc(InType{}, OutType{})
	assertDefinitionFields(t, got, "Port1Action1", []definitionField{
		{name: "S", fieldType: "s:string", minOccurs: "1", maxOccurs: "1"},
		{name: "I", fieldType: "s:int", minOccurs: "1", maxOccurs: "1"},
		{name: "B", fieldType: "s:boolean", minOccurs: "1", maxOccurs: "1"},
		{name: "SP", fieldType: "s:string", minOccurs: "0", maxOccurs: "1"},
		{name: "Ints", fieldType: "s:int", minOccurs: "0", maxOccurs: "unbounded"},
	})
	assertDefinitionFields(t, got, "Port1Action1Response", []definitionField{
		{name: "S", fieldType: "s:string", minOccurs: "1", maxOccurs: "1"},
		{name: "I", fieldType: "s:int", minOccurs: "1", maxOccurs: "1"},
	})
}

func buildDefinitionTestDoc(inType, outType any) *etree.Element {
	return BuildDefinitions(
		"ExampleSoap",
		model.Actions{
			"Port1": map[string]*model.Action{
				"Action1": {
					InType:  inType,
					OutType: outType,
				},
			},
		},
		WithLocation("http://example.org/?wsdl"),
		WithNamespace("http://example.org/"),
	)
}

func assertDefinitionFields(t *testing.T, defElem *etree.Element, elemName string, want []definitionField) {
	t.Helper()

	sequenceElem := findDefinitionSequence(t, defElem, elemName)
	fieldElems := sequenceElem.FindElements("./s:element")
	require.Len(t, fieldElems, len(want))

	for i, field := range want {
		assert.Equal(t, field.name, fieldElems[i].SelectAttrValue("name", ""))
		assert.Equal(t, field.fieldType, fieldElems[i].SelectAttrValue("type", ""))
		assert.Equal(t, field.minOccurs, fieldElems[i].SelectAttrValue("minOccurs", ""))
		assert.Equal(t, field.maxOccurs, fieldElems[i].SelectAttrValue("maxOccurs", ""))
	}
}

func assertDefinitionElementMissing(t *testing.T, defElem *etree.Element, elemName string) {
	t.Helper()

	assert.Nil(t, findDefinitionElement(defElem, elemName))
}

func findDefinitionSequence(t *testing.T, defElem *etree.Element, elemName string) *etree.Element {
	t.Helper()

	elem := findDefinitionElement(defElem, elemName)
	require.NotNil(t, elem)

	sequenceElem := elem.FindElement("./s:complexType/s:sequence")
	require.NotNil(t, sequenceElem)
	return sequenceElem
}

func findDefinitionElement(defElem *etree.Element, elemName string) *etree.Element {
	schemaElem := defElem.FindElement("./wsdl:types/s:schema")
	if schemaElem == nil {
		return nil
	}

	for _, elem := range schemaElem.FindElements("./s:element") {
		if elem.SelectAttrValue("name", "") == elemName {
			return elem
		}
	}

	return nil
}
