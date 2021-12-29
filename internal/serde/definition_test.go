package serde

import (
	_ "embed"
	"testing"

	"github.com/beevik/etree"
	"github.com/francistm/soap-server/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed definition_test_fixture.xml
var definitionFixture string

func TestBuildDefinitions(t *testing.T) {
	type args struct {
		serviceName string
		actions     model.Actions
	}

	type InType struct {
		S  string
		I  int
		B  bool
		SP *string
	}

	type OutType struct {
		S string
		I int
	}

	tests := []struct {
		name string
		args args
		opts []defOpt
		want string
	}{
		{
			name: "case 1",
			args: args{
				serviceName: "ExampleSoap",
				actions: model.Actions{
					"Port1": map[string]*model.Action{
						"Action1": {
							InType:  InType{},
							OutType: OutType{},
						},
					},
				},
			},
			opts: []defOpt{
				WithLocation("http://example.org/?wsdl"),
				WithNamespace("http://example.org/"),
			},
			want: definitionFixture,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := etree.NewDocument()
			got := BuildDefinitions(tt.args.serviceName, tt.args.actions, tt.opts...)
			doc.AddChild(got)
			doc.Indent(4)

			gotXML, err := doc.WriteToString()
			require.NoError(t, err)
			assert.Equal(t, tt.want, gotXML)
		})
	}
}
