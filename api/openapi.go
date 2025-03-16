package api

import _ "embed"

//go:embed openapi.yaml
var openAPISpec []byte

func GetOpenAPISpec() []byte {
	return openAPISpec
}
