package completion

import (
	"fmt"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

func operationNames() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		doc, err := openapi.GetSwagger()
		if err != nil {
			return carapace.ActionValues()
		}
		var values []string
		for path, item := range doc.Paths.Map() {
			for method, op := range item.Operations() {
				values = append(values, op.OperationID, method+" "+path)
			}
		}
		return carapace.ActionValuesDescribed(values...)
	})
}

// Complete parameter names before '=', then schema enums or resource identifiers.
func apiInput(cmd *cobra.Command, store *config.Store, name string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if len(c.Args) == 0 {
			return carapace.ActionValues()
		}
		doc, err := openapi.GetSwagger()
		if err != nil {
			return carapace.ActionValues()
		}
		_, op, params := operation(doc, c.Args[0])
		if op == nil {
			return carapace.ActionValues()
		}
		if name == "content_type" && op.RequestBody != nil {
			var values []string
			for mime := range op.RequestBody.Value.Content {
				values = append(values, mime)
			}
			return carapace.ActionValues(values...)
		}
		if name != "path" && name != "query" {
			return carapace.ActionValues()
		}

		path := ""
		for _, def := range definitions() {
			if def.Operation == c.Args[0] {
				path = def.Path
				break
			}
		}
		return parameterAction(cmd, store, path, name, params)
	})
}

func parameterAction(cmd *cobra.Command, store *config.Store, path, location string, params openapi3.Parameters) carapace.Action {
	return carapace.ActionMultiParts("=", func(c carapace.Context) carapace.Action {
		var values []string
		for _, ref := range params {
			p := ref.Value
			if p.In != location {
				continue
			}
			if len(c.Parts) == 0 {
				values = append(values, p.Name+"=")
				continue
			}
			if c.Parts[0] == p.Name {
				return parameterValues(cmd, store, path, p, c.Value)
			}
		}
		return carapace.ActionValues(values...)
	})
}

func parameterValues(cmd *cobra.Command, store *config.Store, path string, p *openapi3.Parameter, prefix string) carapace.Action {
	if choices := enumValues(p.Schema); len(choices) > 0 {
		return carapace.ActionValues(choices...)
	}
	src, ok := resourceSource(path, p.Name)
	if !ok {
		return carapace.ActionValues()
	}
	bindings := map[string]string{}
	assignments, _ := cmd.Flags().GetStringArray("path")
	for _, assignment := range assignments {
		k, v, ok := strings.Cut(assignment, "=")
		if ok {
			bindings[k] = v
		}
	}
	return lookupAction(cmd, store, src, bindings, prefix)
}

func enumValues(ref *openapi3.SchemaRef) []string {
	if ref == nil || ref.Value == nil {
		return nil
	}
	s := ref.Value
	if s.Type.Is("array") {
		return enumValues(s.Items)
	}
	var values []string
	for _, v := range s.Enum {
		values = append(values, fmt.Sprint(v))
	}
	if len(values) == 0 && s.Type.Is("boolean") {
		return []string{"true", "false"}
	}
	for _, sub := range s.AllOf {
		values = append(values, enumValues(sub)...)
	}
	for _, sub := range s.OneOf {
		values = append(values, enumValues(sub)...)
	}
	for _, sub := range s.AnyOf {
		values = append(values, enumValues(sub)...)
	}
	return values
}
