package main

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func resolveFlags(commands []command, components map[string]parameter) error {
	for i := range commands {
		for j, flag := range commands[i].Flags {
			if flag.Ref == "" {
				continue
			}
			value, ok := components[strings.TrimPrefix(flag.Ref, "#/components/flags/")]
			if !ok {
				return fmt.Errorf("unknown flag reference %q", flag.Ref)
			}
			commands[i].Flags[j] = value
		}
		if err := resolveFlags(commands[i].Commands, components); err != nil {
			return err
		}
	}
	return nil
}

type apiOperation struct {
	spec       *openapi3.Operation
	parameters openapi3.Parameters
}

func validateBindings(commands []command, api *openapi3.T) error {
	operations := map[string]apiOperation{}
	for _, item := range api.Paths.Map() {
		for _, op := range item.Operations() {
			operations[op.OperationID] = apiOperation{op, append(append(openapi3.Parameters{}, item.Parameters...), op.Parameters...)}
		}
	}
	return validateCommands(commands, operations)
}

func validateCommands(commands []command, operations map[string]apiOperation) error {
	for _, c := range commands {
		if strings.HasPrefix(c.ID, "api") {
			op, ok := operations[strings.TrimPrefix(c.ID, "api")]
			if !ok {
				return fmt.Errorf("%s has no matching OpenAPI operation", c.ID)
			}
			if err := validateOperation(c, op); err != nil {
				return err
			}
		}
		if err := validateCommands(c.Commands, operations); err != nil {
			return err
		}
	}
	return nil
}

func validateOperation(c command, op apiOperation) error {
	inputs := map[string]parameter{}
	for _, arg := range c.Arguments {
		inputs["path:"+strings.ToLower(arg.Name)] = arg
	}
	hasBody := false
	for _, flag := range c.Flags {
		if flag.Name == "data" {
			hasBody = true
		}
		if isQuery(flag.Name) {
			inputs["query:"+strings.ReplaceAll(flag.Name, "-", "_")] = flag
		}
	}
	for _, param := range op.parameters {
		p := param.Value
		if p.In != "path" && p.In != "query" {
			continue
		}
		key := p.In + ":" + p.Name
		input, ok := inputs[key]
		if !ok {
			return fmt.Errorf("%s is missing %s", c.ID, key)
		}
		if err := validateParameter(c.ID, input, p); err != nil {
			return err
		}
		delete(inputs, key)
	}
	if len(inputs) > 0 {
		return fmt.Errorf("%s has unknown API parameters: %v", c.ID, inputs)
	}
	if hasBody != (op.spec.RequestBody != nil) {
		return fmt.Errorf("%s body flags do not match the OpenAPI operation", c.ID)
	}
	return nil
}

func validateParameter(id string, input parameter, p *openapi3.Parameter) error {
	if p.In != "query" {
		return nil
	}
	if !input.TrackChanged {
		return fmt.Errorf("%s flag %s needs trackChanged", id, input.Name)
	}
	if p.Schema != nil && p.Schema.Value.Type.Is("array") != input.Repeatable {
		return fmt.Errorf("%s flag %s has incorrect repeatability", id, input.Name)
	}
	return nil
}
