package plugin_manager

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func validateManifestConfiguration(manifest rpc.Manifest, config map[string]any) error {
	schema, hasSchema := manifest.ConfigurationSchema.Get()
	if !hasSchema || len(schema.Fields) == 0 {
		return nil
	}

	fieldsByID := map[string]string{}
	for _, field := range schema.Fields {
		id, typ, ok := extractConfigurationField(field)
		if !ok {
			continue
		}
		fieldsByID[id] = typ
	}

	internalErrors := make([]string, 0)
	externalErrors := make([]string, 0)
	for key, value := range config {
		wantType, ok := fieldsByID[key]
		if !ok {
			// Unknown keys are allowed for now.
			continue
		}

		if matchesConfigurationType(wantType, value) {
			continue
		}

		internalErrors = append(internalErrors, fmt.Sprintf(
			`configuration field %q expects %s but got %T`,
			key,
			wantType,
			value,
		))
		externalErrors = append(externalErrors, fmt.Sprintf(
			`Field %q must be a %s value.`,
			key,
			wantType,
		))
	}

	if len(internalErrors) == 0 {
		return nil
	}

	return fault.New(
		"plugin configuration validation failed",
		ftag.With(ftag.InvalidArgument),
		fmsg.WithDesc(
			strings.Join(internalErrors, "\n"),
			strings.Join(externalErrors, "\n"),
		),
	)
}

func extractConfigurationField(field rpc.PluginConfigurationFieldSchema) (string, string, bool) {
	switch v := field.PluginConfigurationFieldUnion.(type) {
	case *rpc.PluginConfigurationFieldString:
		id := strings.TrimSpace(v.ID)
		if id == "" {
			return "", "", false
		}
		return id, "string", true

	case *rpc.PluginConfigurationFieldNumber:
		id := strings.TrimSpace(v.ID)
		if id == "" {
			return "", "", false
		}
		return id, "number", true

	case *rpc.PluginConfigurationFieldBoolean:
		id := strings.TrimSpace(v.ID)
		if id == "" {
			return "", "", false
		}
		return id, "boolean", true

	default:
		return "", "", false
	}
}

func matchesConfigurationType(expected string, value any) bool {
	switch expected {
	case "string":
		_, ok := value.(string)
		return ok
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "number":
		return isConfigurationNumber(value)
	default:
		return false
	}
}

func isConfigurationNumber(value any) bool {
	switch value.(type) {
	case float64, float32,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		json.Number:
		return true
	default:
		return false
	}
}
