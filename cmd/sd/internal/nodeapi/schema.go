package nodeapi

import (
	"context"
	"fmt"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

// PreparePropertySchema preserves field identity and unmentioned fields because
// the API replaces the schema and deletes values belonging to omitted field IDs.
func PreparePropertySchema(ctx context.Context, client *openapi.ClientWithResponses, parent string, updates []openapi.PropertySchemaMutableProps) ([]openapi.PropertySchemaMutableProps, error) {
	node, err := Fetch(ctx, client, parent)
	if err != nil {
		return nil, err
	}
	if len(node.Children) == 0 {
		return nil, fmt.Errorf("page %q has no children; create a child before setting its schema", parent)
	}
	return mergePropertySchema(node.ChildPropertySchema, updates)
}

func mergePropertySchema(existing []openapi.PropertySchema, updates []openapi.PropertySchemaMutableProps) ([]openapi.PropertySchemaMutableProps, error) {
	result := make([]openapi.PropertySchemaMutableProps, 0, len(existing)+len(updates))
	positions := make(map[string]int, len(existing))
	for _, field := range existing {
		positions[field.Name] = len(result)
		result = append(result, openapi.PropertySchemaMutableProps{Fid: &field.Fid, Name: field.Name, Type: field.Type, Sort: field.Sort})
	}
	seen := make(map[string]bool, len(updates))
	for _, update := range updates {
		if seen[update.Name] {
			return nil, fmt.Errorf("duplicate schema field %q", update.Name)
		}
		seen[update.Name] = true
		if index, ok := positions[update.Name]; ok {
			update.Fid = result[index].Fid
			result[index] = update
		} else {
			result = append(result, update)
		}
	}
	return result, nil
}
