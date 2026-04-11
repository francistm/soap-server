# SOAP Server

## Install

`go get github.com/francistm/soap-server`

## Example

See the [calculator example](./sample/calc/main.go) for details.

When a SOAP port and action are defined, the WSDL definition XML will be served from the endpoint by a GET request.

## Changelog

### v1.0.7
- Add slice and array parameter support
- Add a calculator `Sum` sample action that accepts a slice parameter

### v1.0.6
- Fix missing schema elements when a port has multiple actions
- Move `location` into an option instead of a mandatory parameter

### v1.0.5
- Add `location` when building a SOAP server

### v1.0.4
- Upgrade to go 1.18
- Add context.Context to action handler
- Use generic type for request and response in NewAction

### v1.0.3
- Add response envelope space and namespace options

### v1.0.2
- Add documentation option to `NewAction`
