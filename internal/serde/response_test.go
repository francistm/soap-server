package serde

import (
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildResponseBodyChild(t *testing.T) {
	type args struct {
		ns             string
		soapOutTagName string
		tOut           interface{}
		out            interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				ns:             "http://example.org/",
				soapOutTagName: "FooResponse",
				tOut:           struct{ Foo string }{},
				out:            &struct{ Foo string }{Foo: "blah blah"},
			},
			want: `<FooResponse xmlns="http://example.org/"><Foo>blah blah</Foo></FooResponse>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := etree.NewDocument()

			got := BuildResponseBodyChild(tt.args.ns, tt.args.soapOutTagName, tt.args.tOut, tt.args.out)
			doc.AddChild(got)

			out, err := doc.WriteToString()
			require.NoError(t, err)
			assert.Equal(t, out, tt.want)
		})
	}
}
