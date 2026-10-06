package apiexec

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

func (e *Executor) Operations(out io.Writer) error {
	ops := make([]Operation, 0, len(e.operations))
	for _, op := range e.operations {
		ops = append(ops, op)
	}
	sort.Slice(ops, func(i, j int) bool { return ops[i].ID < ops[j].ID })
	return writeJSON(out, ops)
}

// Schema returns a self-contained OpenAPI excerpt with only reachable components.
func (e *Executor) Schema(out io.Writer, id string) error {
	op, ok := e.operations[id]
	if !ok {
		return fmt.Errorf("unknown API operation %q; use sd api operations", id)
	}
	raw, err := json.Marshal(e.spec)
	if err != nil {
		return err
	}
	var document map[string]any
	if err := json.Unmarshal(raw, &document); err != nil {
		return err
	}
	selected := map[string]any{"openapi": e.spec.OpenAPI, "info": e.spec.Info, "paths": map[string]any{op.Path: map[string]any{strings.ToLower(op.Method): op.spec, "parameters": op.item.Parameters}}}
	raw, err = json.Marshal(selected)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, &selected); err != nil {
		return err
	}
	components := map[string]any{}
	seen := map[string]bool{}
	var collect func(any) error
	collect = func(value any) error {
		switch v := value.(type) {
		case map[string]any:
			if ref, ok := v["$ref"].(string); ok && strings.HasPrefix(ref, "#/components/") && !seen[ref] {
				seen[ref] = true
				parts := strings.Split(strings.TrimPrefix(ref, "#/"), "/")
				var target any = document
				for _, part := range parts {
					target = target.(map[string]any)[part]
				}
				section, ok := components[parts[1]].(map[string]any)
				if !ok {
					section = map[string]any{}
					components[parts[1]] = section
				}
				section[parts[2]] = target
				if err := collect(target); err != nil {
					return err
				}
			}
			for _, child := range v {
				if err := collect(child); err != nil {
					return err
				}
			}
		case []any:
			for _, child := range v {
				if err := collect(child); err != nil {
					return err
				}
			}
		}
		return nil
	}
	if err := collect(selected); err != nil {
		return err
	}
	selected["components"] = components
	return writeJSON(out, selected)
}
