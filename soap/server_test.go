package soap

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	type args struct {
		name   string
		domain string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case 1",
			args: args{
				name:   "Test1",
				domain: "http://example.org",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewService(tt.args.name, tt.args.domain)

			require.NotNil(t, got)
			assert.True(t, strings.HasSuffix(got.namespace, "/"))
		})
	}
}

func TestService_executeAction(t *testing.T) {
	tests := []struct {
		name           string
		s              *Service
		requestBody    string
		statusCodeWant int
	}{
		{
			name: "one-way request",
			s: &Service{
				name:      "test",
				namespace: "http://example.org/",
				actions: map[string]IAction{
					"A1P1": &Action[struct{ P1 string }, *NilOut]{
						handler: func(ctx context.Context, in struct{ P1 string }) (*NilOut, error) {
							return nil, nil
						},
					},
				},
			},
			requestBody: `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:test="http://example.org/">
<soapenv:Header/>
<soapenv:Body>
	<test:A1P1>
		<P1>Blah</P1>
	</test:A1P1>
</soapenv:Body>
</soapenv:Envelope>`,
			statusCodeWant: http.StatusAccepted,
		},
		{
			name: "request & response",
			s: &Service{
				name:      "test",
				namespace: "http://example.org/",
				actions: map[string]IAction{
					"A1P1": &Action[*struct{ P1 string }, *struct{ P2 string }]{
						handler: func(ctx context.Context, in *struct{ P1 string }) (*struct{ P2 string }, error) {
							return &struct{ P2 string }{"foo"}, nil
						},
					},
				},
			},
			requestBody: `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:test="http://example.org/">
<soapenv:Header/>
<soapenv:Body>
	<test:A1P1>
		<P1>Blah</P1>
	</test:A1P1>
</soapenv:Body>
</soapenv:Envelope>`,
			statusCodeWant: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := strings.NewReader(tt.requestBody)
			request := httptest.NewRequest("POST", "http://example.org/wsdl", requestBody)
			responseWriter := httptest.NewRecorder()

			tt.s.executeAction(responseWriter, request)

			assert.Equal(t, tt.statusCodeWant, responseWriter.Result().StatusCode)
		})
	}
}

func TestService_handleSoapOutError(t *testing.T) {
	tests := []struct {
		name         string
		s            *Service
		errorMessage string
	}{
		{
			name: "case 1",
			s: &Service{
				name:      "TestServ",
				namespace: "http://example.org/",
				actions: map[string]IAction{
					"Port1Action1": &Action[any, any]{},
				},
			},
			errorMessage: "test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "http://example.org/wsdl", nil)
			responseWriter := httptest.NewRecorder()
			tt.s.handleSoapOutError(responseWriter, request, errors.New(tt.errorMessage))

			responseResult := responseWriter.Result()

			statusCode := responseResult.StatusCode
			responseBody, _ := io.ReadAll(responseResult.Body)

			assert.Equal(t, http.StatusInternalServerError, statusCode)
			assert.Contains(t, string(responseBody), tt.errorMessage)
		})
	}
}
