package main

import (
	"fmt"
	"go/token"
	"log/slog"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dave/jennifer/jen"
)

func validateSchema(sections []Section) error {
	names := make(map[string]bool)
	variables := make(map[string]bool)

	for _, section := range sections {
		for _, option := range section.Fields {
			if !token.IsIdentifier(option.Name) || names[option.Name] {
				return fmt.Errorf("invalid or duplicate field name %q", option.Name)
			}

			if option.Env == "" || strings.ContainsAny(option.Env, "=\x00") || variables[option.Env] {
				return fmt.Errorf("invalid or duplicate environment variable %q", option.Env)
			}

			names[option.Name] = true
			variables[option.Env] = true

			if err := validateOption(option); err != nil {
				return fmt.Errorf("%s: %w", option.Env, err)
			}
		}
	}

	return nil
}

func validateOption(option ConfigOption) error {
	switch option.Type {
	case "string", "[]byte", "[]string", "bool", "int", "int32", "int64", "float64", "time.Duration", "net/url.URL", "log/slog.Level":
	default:
		return fmt.Errorf("unsupported configuration type %q", option.Type)
	}

	if option.ZeroIsDefault && (!option.ExportDefault || option.Default == nil) {
		return fmt.Errorf("zero_is_default requires an exported default")
	}

	if err := validateBounds(option); err != nil {
		return err
	}

	if option.Default != nil {
		if err := validateDefault(option); err != nil {
			return fmt.Errorf("invalid default: %w", err)
		}
	}

	return nil
}

func validateBounds(option ConfigOption) error {
	if option.Minimum == nil && option.Maximum == nil {
		return nil
	}

	if !numericType(option.Type) {
		return fmt.Errorf("bounds require a numeric type")
	}

	for _, bound := range []*float64{option.Minimum, option.Maximum} {
		if bound == nil {
			continue
		}

		if err := validateBound(option.Type, *bound); err != nil {
			return err
		}
	}

	if option.Minimum != nil && option.Maximum != nil && *option.Minimum > *option.Maximum {
		return fmt.Errorf("minimum exceeds maximum")
	}

	return nil
}

func validateBound(typeName string, bound float64) error {
	if math.IsNaN(bound) || math.IsInf(bound, 0) {
		return fmt.Errorf("bounds must be finite")
	}

	if typeName == "float64" {
		return nil
	}

	if math.Trunc(bound) != bound {
		return fmt.Errorf("integer bounds must be whole numbers")
	}

	_, err := strconv.ParseInt(strconv.FormatFloat(bound, 'f', -1, 64), 10, integerBits(typeName))
	if err != nil {
		return fmt.Errorf("bound is outside the range of %s", typeName)
	}

	return nil
}

func numericType(name string) bool {
	switch name {
	case "int", "int32", "int64", "float64", "time.Duration":
		return true
	default:
		return false
	}
}

func validateDefault(option ConfigOption) error {
	value := *option.Default
	var number float64
	var err error

	switch option.Type {
	case "int", "int32", "int64":
		var parsed int64
		parsed, err = strconv.ParseInt(value, 0, integerBits(option.Type))
		number = float64(parsed)
	case "float64":
		number, err = strconv.ParseFloat(value, 64)
	case "bool":
		_, err = strconv.ParseBool(value)
	case "time.Duration":
		var parsed time.Duration
		parsed, err = time.ParseDuration(value)
		number = float64(parsed)
	case "net/url.URL":
		_, err = url.Parse(value)
	case "log/slog.Level":
		var level slog.Level
		err = level.UnmarshalText([]byte(value))
	}

	if err != nil {
		return err
	}

	if option.Required && value == "" {
		return fmt.Errorf("required value cannot be empty")
	}

	if math.IsNaN(number) || math.IsInf(number, 0) {
		return fmt.Errorf("value must be finite")
	}

	if option.Minimum != nil && number < *option.Minimum {
		return fmt.Errorf("value is below minimum")
	}

	if option.Maximum != nil && number > *option.Maximum {
		return fmt.Errorf("value exceeds maximum")
	}

	return nil
}

func integerBits(name string) int {
	switch name {
	case "int32":
		return 32
	case "int64":
		return 64
	default:
		return 0
	}
}

func generateLoading(file *jen.File, sections []Section) {
	file.Var().Id("errInvalidConfiguration").Op("=").Qual("errors", "New").Call(jen.Lit("invalid configuration"))
	file.Line()

	file.Func().Id("Load").Params().Params(jen.Id("Config"), jen.Error()).Block(
		jen.Return(jen.Id("load").Call(jen.Qual("os", "LookupEnv"))),
	)
	file.Line()

	file.Func().Id("load").Params(
		jen.Id("lookup").Func().Params(jen.String()).Params(jen.String(), jen.Bool()),
	).Params(jen.Id("Config"), jen.Error()).BlockFunc(func(body *jen.Group) {
		body.Var().Id("c").Id("Config")

		for _, section := range sections {
			for _, option := range section.Fields {
				body.Line().BlockFunc(func(field *jen.Group) {
					generateLookup(field, option)
				})
			}
		}

		body.Line().Return(jen.Id("ApplyDefaults").Call(jen.Id("c")))
	})
	file.Line()

	file.Func().Id("ApplyDefaults").Params(jen.Id("c").Id("Config")).Params(jen.Id("Config"), jen.Error()).BlockFunc(func(body *jen.Group) {
		for _, section := range sections {
			for _, option := range section.Fields {
				if !option.ZeroIsDefault {
					continue
				}

				body.If(jen.Id("c").Dot(option.Name).Op("==").Lit(0)).Block(
					jen.Id("c").Dot(option.Name).Op("=").Id("Default" + option.Name),
				).Line()
			}
		}

		body.Return(jen.Id("c"), jen.Id("c").Dot("Validate").Call())
	})
	file.Line()

	file.Func().Params(jen.Id("c").Id("Config")).Id("Validate").Params().Error().BlockFunc(func(body *jen.Group) {
		for _, section := range sections {
			for _, option := range section.Fields {
				generateValidation(body, option)
			}
		}

		body.Return(jen.Nil())
	})
}

func generateLookup(body *jen.Group, option ConfigOption) {
	body.List(jen.Id("value"), jen.Id("present")).Op(":=").Id("lookup").Call(jen.Lit(option.Env))

	if option.Default != nil {
		body.If(jen.Op("!").Id("present")).Block(jen.Id("value").Op("=").Lit(*option.Default)).Line()
	}

	if option.Required {
		condition := jen.Id("value").Op("==").Lit("")
		if option.Default == nil {
			condition = jen.Op("!").Id("present").Op("||").Add(condition)
		}

		body.If(condition).Block(
			jen.Return(jen.Id("c"), configurationError(option, "is required")),
		).Line()
	}

	if option.Default != nil || option.Required {
		generateParse(body, option)
		return
	}

	body.If(jen.Id("present")).BlockFunc(func(field *jen.Group) {
		generateParse(field, option)
	})
}

func generateParse(body *jen.Group, option ConfigOption) {
	target := jen.Id("c").Dot(option.Name)
	value := jen.Id("value")
	var parser *jen.Statement

	switch option.Type {
	case "string":
		body.Add(target).Op("=").Add(value)
		return
	case "[]byte":
		body.Add(target).Op("=").Index().Byte().Call(value)
		return
	case "[]string":
		body.Add(target).Op("=").Index().String().Values()
		body.If(jen.Qual("strings", "TrimSpace").Call(value).Op("!=").Lit("")).Block(
			jen.Add(target).Op("=").Qual("strings", "Split").Call(value, jen.Lit(",")),
		)
		return
	case "log/slog.Level":
		body.If(jen.Err().Op(":=").Add(target).Dot("UnmarshalText").Call(jen.Index().Byte().Call(value)), jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Id("c"), configurationError(option, "must be a log level")),
		)
		return
	case "net/url.URL":
		parser = jen.Qual("net/url", "Parse").Call(value)
	case "time.Duration":
		parser = jen.Qual("time", "ParseDuration").Call(value)
	case "bool":
		parser = jen.Qual("strconv", "ParseBool").Call(value)
	case "float64":
		parser = jen.Qual("strconv", "ParseFloat").Call(value, jen.Lit(64))
	default:
		parser = jen.Qual("strconv", "ParseInt").Call(value, jen.Lit(0), jen.Lit(integerBits(option.Type)))
	}

	body.List(jen.Id("parsed"), jen.Err()).Op(":=").Add(parser)
	body.If(jen.Err().Op("!=").Nil()).Block(
		jen.Return(jen.Id("c"), configurationError(option, "has an invalid "+option.Type+" value")),
	).Line()

	switch option.Type {
	case "net/url.URL":
		body.Add(target).Op("=").Op("*").Id("parsed")
	case "int", "int32":
		body.Add(target).Op("=").Id(option.Type).Call(jen.Id("parsed"))
	default:
		body.Add(target).Op("=").Id("parsed")
	}
}

func generateValidation(body *jen.Group, option ConfigOption) {
	value := jen.Id("c").Dot(option.Name)
	if option.Type == "float64" {
		body.If(jen.Qual("math", "IsNaN").Call(value).Op("||").Qual("math", "IsInf").Call(value, jen.Lit(0))).Block(
			jen.Return(configurationError(option, "must be finite")),
		).Line()
	}

	if option.Minimum != nil {
		body.If(jen.Add(value).Op("<").Lit(*option.Minimum)).Block(
			jen.Return(configurationError(option, fmt.Sprintf("must be at least %g", *option.Minimum))),
		).Line()
	}

	if option.Maximum != nil {
		body.If(jen.Add(value).Op(">").Lit(*option.Maximum)).Block(
			jen.Return(configurationError(option, fmt.Sprintf("must be at most %g", *option.Maximum))),
		).Line()
	}
}

func configurationError(option ConfigOption, message string) *jen.Statement {
	return jen.Qual("fmt", "Errorf").Call(jen.Lit("%w: "+option.Env+" "+message), jen.Id("errInvalidConfiguration"))
}
