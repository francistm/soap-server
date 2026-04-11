package serde

import (
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSoapInArrayParameters(t *testing.T) {
	type InType struct {
		IDs []int
	}

	elem := etree.NewElement("Port1Action1")
	elem.CreateElement("IDs").SetText("1")
	elem.CreateElement("IDs").SetText("2")

	got, err := ParseSoapIn(elem, InType{})
	require.NoError(t, err)
	assert.Equal(t, InType{IDs: []int{1, 2}}, got)
}

func TestParseSoapInOptionalArrayParameters(t *testing.T) {
	type InType struct {
		IDs *[]int
	}

	elem := etree.NewElement("Port1Action1")
	elem.CreateElement("IDs").SetText("1")
	elem.CreateElement("IDs").SetText("2")

	got, err := ParseSoapIn(elem, InType{})
	require.NoError(t, err)
	assert.Equal(t, InType{IDs: &[]int{1, 2}}, got)
}

func TestParseSoapInArrayParameterLengthMismatch(t *testing.T) {
	type InType struct {
		IDs [2]int
	}

	elem := etree.NewElement("Port1Action1")
	elem.CreateElement("IDs").SetText("1")

	_, err := ParseSoapIn(elem, InType{})
	require.Error(t, err)
}
