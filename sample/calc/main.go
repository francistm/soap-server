package main

import (
	"context"
	"net"
	"net/http"

	"github.com/francistm/soap-server/soap"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/wsdl", wsdlHandler())

	lis, err := net.Listen("tcp4", ":8080")

	if err != nil {
		panic(err)
	}

	if err := http.Serve(lis, serveMux); err != nil {
		panic(err)
	}
}

type AddRequest struct {
	Int1 int
	Int2 int
}

type AddResponse struct {
	Acc int
}

func wsdlHandler() http.HandlerFunc {
	addAction := soap.NewAction(func(ctx context.Context, input *AddRequest) (*AddResponse, error) {
		output := &AddResponse{
			Acc: input.Int1 + input.Int2,
		}

		return output, nil
	})

	soapPort := soap.NewPort()
	soapServer := soap.NewService("Calculator", "http://example.org/")

	soapPort.AddAction("Add", addAction)
	soapServer.AddPort("Calculator", soapPort)

	return func(rw http.ResponseWriter, r *http.Request) {
		soapServer.ServeHTTP(rw, r)
	}
}
