package main

//go:generate go tool opencli generate go-cobra opencli.yaml --package cligen --output internal/cligen

//go:generate go run ./internal/apigen
