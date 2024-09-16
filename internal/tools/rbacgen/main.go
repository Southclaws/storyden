package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dave/jennifer/jen"
	"github.com/pb33f/libopenapi"
)

func main() {
	schemaFlag := flag.String("schema", "api/openapi.yaml", "path to openapi schema")
	outputFlag := flag.String("output", "app/transports/http/bindings/openapi_rbac/openapi_rbac_gen.go", "path to output file")

	flag.Parse()

	filename := *schemaFlag
	outfile := *outputFlag

	if err := run(filename, outfile); err != nil {
		fmt.Printf("Error: %e\n", err)
		os.Exit(1)
	}
}

type Operation struct {
	Name string
}

func run(filename, outfile string) error {
	spec, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	document, err := libopenapi.NewDocument(spec)
	if err != nil {
		return fmt.Errorf("cannot create new document: %w", err)
	}

	docModel, errors := document.BuildV3Model()
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("error: %e\n", errors[i])
		}

		return fmt.Errorf("cannot create v3 model from document: %d errors reported", len(errors))
	}

	ops := []Operation{}

	for _, path := range docModel.Model.Paths.PathItems.FromOldest() {
		for _, op := range path.GetOperations().FromOldest() {
			ops = append(ops, Operation{
				Name: op.OperationId,
			})
		}
	}

	f := jen.NewFile("openapi_rbac")

	f.ImportName("github.com/Southclaws/storyden/app/resources/rbac", "rbac")

	funcs := []jen.Code{}
	cases := []jen.Code{}

	for _, op := range ops {
		funcs = append(funcs, jen.Id(op.Name).Params().Params(
			jen.Bool(),
			jen.Op("*").Qual("github.com/Southclaws/storyden/app/resources/rbac", "Permission"),
		))

		cases = append(cases, jen.Case(jen.Lit(op.Name)).Block(
			jen.Return(
				jen.Id("optable").Dot(op.Name).Params(),
			),
		))
	}

	cases = append(cases, jen.Default().Block(
		jen.Panic(jen.Lit("unknown operation, must re-run rbacgen")),
	))

	f.Type().Id("OperationPermissions").Interface(
		funcs...,
	)

	f.Func().
		Id("GetOperationPermission").
		Params(
			jen.Id("optable").Id("OperationPermissions"),
			jen.Id("op").String(),
		).
		Params(
			jen.Bool(),
			jen.Op("*").Qual("github.com/Southclaws/storyden/app/resources/rbac", "Permission"),
		).Block(
		jen.Switch(jen.Id("op")).Block(cases...),
	)

	return f.Save(outfile)
}
