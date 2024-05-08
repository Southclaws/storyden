// Package openapi provides a transport layer using HTTP. In this layer, most of
// the code is generated: low level HTTP handlers, request and response structs
// and object validation.
//
// This is wired up to the service layer using "Bindings" which are just glue
// code which call service APIs from the endpoint handlers and deal with errors.
package openapi
