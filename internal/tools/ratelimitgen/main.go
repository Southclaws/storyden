package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dave/jennifer/jen"
	"github.com/pb33f/libopenapi"
)

func main() {
	schemaFlag := flag.String("schema", "api/openapi.yaml", "path to openapi schema")
	outputFlag := flag.String("output", "app/transports/http/middleware/limiter/ratelimit_config_gen.go", "path to output file")

	flag.Parse()

	filename := *schemaFlag
	outfile := *outputFlag

	if err := run(filename, outfile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

type RateLimitConfig struct {
	OperationID string
	Cost        int
	Limit       int
	Period      string
}

type RouteMapping struct {
	Path        string
	Method      string
	OperationID string
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

	docModel, err := document.BuildV3Model()
	if err != nil {
		return fmt.Errorf("cannot create v3 model from document: %w", err)
	}

	configs := []RateLimitConfig{}
	routeMappings := []RouteMapping{}

	// Iterate through paths to get path string and operations
	for pathPair := docModel.Model.Paths.PathItems.First(); pathPair != nil; pathPair = pathPair.Next() {
		path := pathPair.Key()
		pathItem := pathPair.Value()

		for opPair := pathItem.GetOperations().First(); opPair != nil; opPair = opPair.Next() {
			method := opPair.Key()
			op := opPair.Value()
			
			// Add route mapping for all operations
			routeMappings = append(routeMappings, RouteMapping{
				Path:        path,
				Method:      method,
				OperationID: op.OperationId,
			})

			// Check for x-storyden extension
			if op.Extensions == nil {
				continue
			}

			storydenNode, ok := op.Extensions.Get("x-storyden")
			if !ok || storydenNode == nil {
				continue
			}

			// We need to manually parse the YAML node
			// The structure is a map with "rateLimit" key
			var storydenData map[string]interface{}
			if err := storydenNode.Decode(&storydenData); err != nil {
				continue
			}

			rateLimitData, ok := storydenData["rateLimit"].(map[string]interface{})
			if !ok {
				continue
			}

			config := RateLimitConfig{
				OperationID: op.OperationId,
				Cost:        1,  // default
				Limit:       0,  // 0 means use global default
				Period:      "", // empty means use global default
			}

			if cost, ok := rateLimitData["cost"]; ok {
				switch v := cost.(type) {
				case int:
					config.Cost = v
				case float64:
					config.Cost = int(v)
				}
			}

			if limit, ok := rateLimitData["limit"]; ok {
				switch v := limit.(type) {
				case int:
					config.Limit = v
				case float64:
					config.Limit = int(v)
				}
			}

			if period, ok := rateLimitData["period"]; ok {
				if s, ok := period.(string); ok {
					config.Period = s
				}
			}

			configs = append(configs, config)
		}
	}

	return generateCode(configs, routeMappings, outfile)
}

func generateCode(configs []RateLimitConfig, routeMappings []RouteMapping, outfile string) error {
	f := jen.NewFile("limiter")

	f.PackageComment("Package limiter contains rate limiting middleware.")
	f.PackageComment("THIS FILE IS GENERATED. DO NOT EDIT MANUALLY.")
	f.PackageComment("To edit rate limit configuration, edit the OpenAPI spec x-storyden extensions and run codegen.")

	f.ImportName("time", "time")

	// Generate OperationRateLimitConfig struct
	f.Comment("OperationRateLimitConfig defines per-operation rate limiting configuration")
	f.Type().Id("OperationRateLimitConfig").Struct(
		jen.Comment("Cost is the number of requests this operation counts as"),
		jen.Id("Cost").Int(),
		jen.Comment("Limit is the maximum number of requests allowed in the period (0 means use global default)"),
		jen.Id("Limit").Int(),
		jen.Comment("Period is the time window for the limit (empty means use global default)"),
		jen.Id("Period").Qual("time", "Duration"),
	)

	// Generate the map of operation configs
	f.Comment("OperationRateLimits contains per-operation rate limit configurations extracted from OpenAPI spec")
	f.Var().Id("OperationRateLimits").Op("=").Map(jen.String()).Id("OperationRateLimitConfig").Values(
		jen.DictFunc(func(d jen.Dict) {
			for _, config := range configs {
				var period jen.Code
				if config.Period == "" {
					period = jen.Lit(0)
				} else {
					// Parse duration string
					_, err := time.ParseDuration(config.Period)
					if err != nil {
						// If parsing fails, use 0 (global default)
						period = jen.Lit(0)
					} else {
						period = jen.Qual("time", "Duration").Call(
							jen.Lit(parseDurationToNanoseconds(config.Period)),
						)
					}
				}

				d[jen.Lit(config.OperationID)] = jen.Values(jen.Dict{
					jen.Id("Cost"):   jen.Lit(config.Cost),
					jen.Id("Limit"):  jen.Lit(config.Limit),
					jen.Id("Period"): period,
				})
			}
		}),
	)

	// Generate route to operation mapping
	f.Comment("RouteToOperation maps HTTP method and path to operation ID")
	f.Var().Id("RouteToOperation").Op("=").Map(jen.String()).String().Values(
		jen.DictFunc(func(d jen.Dict) {
			for _, mapping := range routeMappings {
				// Convert OpenAPI path format {param} to Echo format :param
				echoPath := convertPathFormat(mapping.Path)
				key := mapping.Method + ":/api" + echoPath
				d[jen.Lit(key)] = jen.Lit(mapping.OperationID)
			}
		}),
	)

	// Generate helper function
	f.Comment("GetOperationConfig returns the rate limit config for an operation, or nil if not configured")
	f.Func().Id("GetOperationConfig").Params(
		jen.Id("operationID").String(),
	).Params(
		jen.Op("*").Id("OperationRateLimitConfig"),
	).Block(
		jen.If(
			jen.List(jen.Id("cfg"), jen.Id("ok")).Op(":=").Id("OperationRateLimits").Index(jen.Id("operationID")),
			jen.Id("ok"),
		).Block(
			jen.Return(jen.Op("&").Id("cfg")),
		),
		jen.Return(jen.Nil()),
	)

	f.Comment("GetOperationIDFromRoute returns the operation ID for a given route")
	f.Func().Id("GetOperationIDFromRoute").Params(
		jen.Id("method").String(),
		jen.Id("path").String(),
	).Params(
		jen.String(),
	).Block(
		jen.Id("key").Op(":=").Id("method").Op("+").Lit(":").Op("+").Id("path"),
		jen.If(
			jen.List(jen.Id("opID"), jen.Id("ok")).Op(":=").Id("RouteToOperation").Index(jen.Id("key")),
			jen.Id("ok"),
		).Block(
			jen.Return(jen.Id("opID")),
		),
		jen.Return(jen.Lit("")),
	)

	return f.Save(outfile)
}

func parseDurationToNanoseconds(s string) int64 {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0
	}
	return int64(d)
}

func convertPathFormat(openAPIPath string) string {
	// Convert OpenAPI path format {param} to Echo format :param
	result := ""
	for _, ch := range openAPIPath {
		if ch == '{' {
			result += ":"
		} else if ch == '}' {
			// Skip closing brace
			continue
		} else {
			result += string(ch)
		}
	}
	return result
}

