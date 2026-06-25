package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Southclaws/enumerator/generate"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
)

func main() {
	// meant to be run from ./api
	schemaFlag := flag.String("schema", "openapi.yaml", "path to openapi schema")
	outputFlag := flag.String("output", "../app/transports/http/bindings/openapi_rbac/openapi_rbac_gen.go", "path to output file")
	enumFlag := flag.String("enum", "../app/resources/rbac/rbac_enum_gen.go", "path to enum output file")
	operationEnumFlag := flag.String("operation-enum", "../app/transports/http/openapi/operation/operation_enum_gen.go", "path to operation enum output file")
	operationCostFlag := flag.String("operation-cost", "../app/transports/http/openapi/operation/operation_cost_gen.go", "path to operation cost map output file")

	flag.Parse()

	filename := *schemaFlag
	outfile := *outputFlag
	enumOutfile := *enumFlag
	operationEnumOutfile := *operationEnumFlag
	operationCostOutfile := *operationCostFlag

	if err := run(filename, outfile, enumOutfile, operationEnumOutfile, operationCostOutfile); err != nil {
		fmt.Printf("Error: %e\n", err)
		os.Exit(1)
	}
}

type Operation struct {
	Name          string
	RateLimitCost int
}

func run(filename, outfile, enumOutfilePath, operationEnumOutfile, operationCostOutfile string) error {
	spec, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	config := datamodel.NewDocumentConfiguration()
	config.AllowFileReferences = true
	config.BasePath = filepath.Dir(filename)

	document, err := libopenapi.NewDocumentWithConfiguration(spec, config)
	if err != nil {
		return fmt.Errorf("cannot create new document: %w", err)
	}

	docModel, err := document.BuildV3Model()
	if err != nil {
		return fmt.Errorf("cannot create v3 model from document: %w", err)
	}

	ops := []Operation{}

	for _, path := range docModel.Model.Paths.PathItems.FromOldest() {
		for _, op := range path.GetOperations().FromOldest() {
			rateLimitCost := 0

			if op.Extensions != nil {
				xStorydenNode := op.Extensions.GetOrZero("x-storyden")
				if xStorydenNode != nil && xStorydenNode.Content != nil {
					for i := 0; i < len(xStorydenNode.Content); i += 2 {
						if i+1 < len(xStorydenNode.Content) {
							keyNode := xStorydenNode.Content[i]
							valueNode := xStorydenNode.Content[i+1]
							if keyNode.Value == "rate-limit-cost" {
								if costInt, err := strconv.Atoi(valueNode.Value); err == nil {
									rateLimitCost = costInt
								}
							}
						}
					}
				}
			}

			ops = append(ops, Operation{
				Name:          op.OperationId,
				RateLimitCost: rateLimitCost,
			})
		}
	}

	enum := generate.Enum{
		Name:       "Permission",
		Values:     []generate.Value{},
		Sourcename: "string",
	}
	for name, schema := range docModel.Model.Components.Schemas.FromOldest() {
		if name == "Permission" {
			for _, v := range schema.Schema().Enum {
				enum.Values = append(enum.Values, generate.Value{
					Symbol: fmt.Sprintf("Permission%s", strcase.ToCamel(v.Value)),
					Value:  fmt.Sprintf("`%s`", v.Value),
				})
			}
		}
	}

	enumOutfile, err := os.Create(enumOutfilePath)
	if err != nil {
		return err
	}

	err = generate.Generate("rbac", []generate.Enum{enum}, enumOutfile)
	if err != nil {
		return err
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

	if err := f.Save(outfile); err != nil {
		return err
	}

	operationEnum := generate.Enum{
		Name:       "OperationID",
		Values:     []generate.Value{},
		Sourcename: "string",
	}

	for _, op := range ops {
		operationEnum.Values = append(operationEnum.Values, generate.Value{
			Symbol: fmt.Sprintf("OperationID%s", op.Name),
			Value:  fmt.Sprintf("`%s`", op.Name),
		})
	}

	operationEnumFile, err := os.Create(operationEnumOutfile)
	if err != nil {
		return err
	}

	if err := generate.Generate("operation", []generate.Enum{operationEnum}, operationEnumFile); err != nil {
		return err
	}

	costFile := jen.NewFile("operation")

	costMapValues := []jen.Code{}
	for _, op := range ops {
		if op.RateLimitCost > 0 {
			costMapValues = append(costMapValues,
				jen.Id(fmt.Sprintf("OperationID%s", op.Name)).Op(":").Lit(op.RateLimitCost),
			)
		}
	}

	costFile.Var().Id("RateLimitCostOverrides").Op("=").Map(jen.Id("OperationID")).Int().Values(
		costMapValues...,
	)

	allOpsValues := []jen.Code{}
	for _, op := range ops {
		allOpsValues = append(allOpsValues, jen.Id(fmt.Sprintf("OperationID%s", op.Name)))
	}

	costFile.Var().Id("AllOperations").Op("=").Index().Id("OperationID").Values(
		allOpsValues...,
	)

	return costFile.Save(operationCostOutfile)
}
