# SOAP Server

## Install
`go get github.com/francistm/soap-server`

## Example
Details see [calculator example](./sample/calc/main.go)

When SOAP port & action defined, WSDL definition XML will be served in the endpoint by GET request.

## Changelog
### v1.0.6
- Fix element missing in schema when port have multiple actions
- Move `location` option as option instead mandatory parameter

### v1.0.5
- Add `location` when build soap server

### v1.0.4
- Upgrade to go 1.18
- Add context.Context to action handler
- Use generic type for request and response in NewAction

### v1.0.3
Add response envelope space, ns option

### v1.0.2
Add documetation option to NewAction