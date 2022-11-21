package openapi

//go:generate go run -mod=mod github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.11.1-0.20220906181851-9c600dddea33 --config ../../api/config.yaml ../../api/openapi.yaml

//go:generate go run -mod=mod github.com/ogen-go/ogen/cmd/ogen@v0.54.1 --debug.ignoreNotImplemented "cookie security, complex schema merging" --target ./ogen --package ogen --clean ../../api/openapi.yaml
