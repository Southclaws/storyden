package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/assert/yaml"
)

func main() {
	schemaFlag := flag.String("schema", "internal/config/config.yaml", "path to config schema")
	outputPkgFlag := flag.String("pkg", "internal/config/config.go", "path to output config.go struct file")
	outputDocFlag := flag.String("doc", "home/content/docs/reference/configuration.mdx", "path to output docs markdown file")

	flag.Parse()

	filename := *schemaFlag
	outputPkgFile := *outputPkgFlag
	outputDocFile := *outputDocFlag

	if err := run(filename, outputPkgFile, outputDocFile); err != nil {
		fmt.Printf("Error: %e\n", err)
		os.Exit(1)
	}
}

type Section struct {
	Name        string         `yaml:"section"`
	Description string         `yaml:"description"`
	Fields      []ConfigOption `yaml:"fields"`
}

type ConfigOption struct {
	Env         string  `yaml:"env"`
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Type        string  `yaml:"type"`
	Required    bool    `yaml:"required"`
	Default     *string `yaml:"default"`
}

func run(filename, outputPkgFile, outputDocFile string) error {
	spec, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var sections []Section

	err = yaml.Unmarshal(spec, &sections)
	if err != nil {
		return err
	}

	f := jen.NewFile("config")

	f.PackageComment(`Package config contains all environment variable based configuration.`)
	f.PackageComment(`THIS FILE IS GENERATED. DO NOT EDIT MANUALLY.`)
	f.PackageComment(`To edit configuration variables, edit the config.yaml file and run codegen.`)

	markdown := strings.Builder{}

	markdown.WriteString(`---
title: Configuration via Environment Variables
description: Reference for configuring Storyden with environment variables.
---

"Configuration" throughout the documentation and codebase refers to these variables which are statically set when the process launches. Changing them requires a restart, however they don't need to be changed much once set.

The term "Settings" is distinct from these and refers to runtime-configurable values which are stored in the database and not via environment variables. These can be changed at any time via the API or the Admin settings page.

`)

	f.Comment("Config represents environment variable configuration parameters").Line()
	f.Type().Id("Config").StructFunc(func(g *jen.Group) {
		for _, section := range sections {

			markdown.WriteString(fmt.Sprintf("## %s\n\n", section.Name))
			markdown.WriteString(fmt.Sprintf("%s\n\n", section.Description))

			g.Line().Comment("-")
			g.Comment(section.Name)
			g.Comment("-").Line()

			for _, f := range section.Fields {
				// Generate golang code

				fieldName := f.Name
				if f.Required {
					fieldName = fieldName + "Required"
				}

				tags := map[string]string{
					"envconfig": f.Env,
				}

				if f.Default != nil {
					tags["default"] = *f.Default
				}

				g.Comment(f.Description)

				switch getTypeType(f.Type) {
				case primitive:
					g.Id(fieldName).Id(f.Type).Tag(tags)

				case stdlib:
					path, sym := splitImportPath(f.Type)

					g.Id(fieldName).Qual(path, sym).Tag(tags)

				case thirdParty:
					path, sym := splitImportPath(f.Type)

					g.Id(fieldName).Qual(path, sym).Tag(tags)
				}

				properties := fmt.Sprintf(`<table>
<tr><td>type</td><td>%s</td></tr>
<tr><td>default</td><td>%s</td></tr>
</table>

`, getNonTechnicalTypeName(f.Type), getPrettyDefault(f.Default))

				// Generate markdown documentation
				markdown.WriteString(fmt.Sprintf("### `%s`\n\n", f.Env))
				markdown.WriteString(properties)
				// if f.Default != nil {
				// 	markdown.WriteString(fmt.Sprintf("> default: `%s`\n\n", *f.Default))
				// }
				markdown.WriteString(fmt.Sprintf("%s\n\n", f.Description))
				// markdown.WriteString(fmt.Sprintf("Type: %s\n\n", getNonTechnicalTypeName(f.Type)))
			}
		}
	})

	// Generate the markdown documentation
	docFile, err := os.Create(outputDocFile)
	if err != nil {
		return err
	}
	defer docFile.Close()
	_, err = docFile.WriteString(markdown.String())
	if err != nil {
		return err
	}

	return f.Save(outputPkgFile)
}

const (
	// Primitive types, string, int, etc
	primitive = iota
	// stdlib import types, time.Time, url.URL, etc
	stdlib
	// third party imports
	thirdParty
)

func getTypeType(s string) int {
	if strings.Contains(s, "/") {
		return thirdParty
	}

	if strings.Contains(s, ".") {
		return stdlib
	}

	return primitive
}

func splitImportPath(s string) (string, string) {
	// slashparts := strings.Split(s, "/")
	dotparts := strings.Split(s, ".")

	return strings.Join(dotparts[:len(dotparts)-1], "."), dotparts[len(dotparts)-1]
}

func getNonTechnicalTypeName(s string) string {
	switch s {
	case "log/slog.Level":
		return "`debug`, `info`, `warn`, `error`"
	case "time.Duration":
		return "duration (e.g. 1h, 1m, 1s)"
	case "net/url.URL":
		return "url (e.g. http://example.com)"
	case "string":
		return "`string`"
	case "int":
		return "`integer` (number without decimal point)"
	case "bool":
		return "boolean (`true` or `false`, case sensitive)"
	case "float64":
		return "float (e.g. `1.0`, `1.5`)"
	case "float32":
		return "float (e.g. `1.0`, `1.5`)"
	case "int64":
		return "integer (e.g. `1`, `2`, `3`)"
	case "int32":
		return "integer (e.g. `1`, `2`, `3`)"
	case "int16":
		return "integer (e.g. `1`, `2`, `3`)"
	case "int8":
		return "integer (e.g. `1`, `2`, `3`)"
	case "uint64":
		return "unsigned integer (e.g. `1`, `2`, `3`)"
	case "uint32":
		return "unsigned integer (e.g. `1`, `2`, `3`)"
	case "uint16":
		return "unsigned integer (e.g. `1`, `2`, `3`)"
	case "uint8":
		return "unsigned integer (e.g. `1`, `2`, `3`)"
	default:
		return fmt.Sprintf("`%s`", s)
	}
}

func getPrettyDefault(s *string) string {
	if s == nil {
		return "none"
	}

	if *s == "" {
		return "(empty string)"
	}

	return fmt.Sprintf("`%s`", *s)
}
