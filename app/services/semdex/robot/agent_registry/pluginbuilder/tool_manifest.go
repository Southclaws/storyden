package pluginbuilder

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/google/jsonschema-go/jsonschema"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	libplugin "github.com/Southclaws/storyden/lib/plugin"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type ManifestWriteResult struct {
	Path     string `json:"path"`
	ID       string `json:"id"`
	Bytes    int    `json:"bytes"`
	Revision string `json:"revision"`
	Message  string `json:"message,omitempty"`
}

const managedPluginCommand = "go"

var managedPluginArgs = []string{"run", "."}

func (a *Agent) addManifestTools(add toolAdder) error {
	return add(functiontool.New[map[string]any, ManifestWriteResult](functiontool.Config{
		Name:        "plugin_manifest_write",
		Description: manifestWriteDescription,
		InputSchema: pluginManifestToolInputSchema(),
	}, func(ctx adktool.Context, args map[string]any) (ManifestWriteResult, error) {
		result, err := a.WriteManifest(ctx, args)
		if err != nil {
			return ManifestWriteResult{}, err
		}
		return result, nil
	}))
}

const manifestWriteDescription = `Write and validate manifest.yaml from structured manifest fields.

Use this instead of plugin_file_edit or plugin_file_write when changing
manifest.yaml. Use configuration_schema for plugin settings, never
configuration. The tool validates the manifest before writing it and manages
runtime launch fields automatically.`

func (a *Agent) WriteManifest(ctx context.Context, raw map[string]any) (ManifestWriteResult, error) {
	if raw == nil {
		return ManifestWriteResult{}, errors.New("manifest input is required")
	}
	manifestRaw := normalizeManagedManifestRaw(raw)
	if err := validateManifestRaw(manifestRaw); err != nil {
		return ManifestWriteResult{}, err
	}

	manifest, err := rpc.ManifestFromMap(manifestRaw)
	if err != nil {
		return ManifestWriteResult{}, fmt.Errorf("validate manifest: %w", err)
	}

	target, ok, err := pluginBuildTargetFromContext(ctx)
	if err != nil {
		return ManifestWriteResult{}, err
	}
	if !ok {
		return ManifestWriteResult{}, errors.New("create or import a plugin before writing manifest.yaml")
	}
	if target.ManifestID != "" && target.ManifestID != manifest.ID {
		return ManifestWriteResult{}, errors.New(pluginBuildTargetDifferentPluginMessage)
	}
	if target.InstallationID != "" {
		if err := ensurePluginBuildTarget(ctx, manifest.ID, target.InstallationID); err != nil {
			return ManifestWriteResult{}, err
		}
	}

	workspace, err := a.Workspace(ctx)
	if err != nil {
		return ManifestWriteResult{}, err
	}

	out, err := yaml.Marshal(manifest.ToMap())
	if err != nil {
		return ManifestWriteResult{}, fmt.Errorf("render manifest.yaml: %w", err)
	}

	written, err := workspace.WriteFile(ctx, manifestYAMLFilename, out)
	if err != nil {
		return ManifestWriteResult{}, err
	}

	return ManifestWriteResult{
		Path:     written.Path,
		ID:       manifest.ID,
		Bytes:    written.Bytes,
		Revision: contentRevision(out),
		Message:  "manifest.yaml written and validated",
	}, nil
}

func normalizeManagedManifestRaw(raw map[string]any) map[string]any {
	normalised := make(map[string]any, len(raw)+2)
	for key, value := range raw {
		normalised[key] = value
	}
	normalised["command"] = managedPluginCommand
	normalised["args"] = managedPluginArgsAny()
	return normalised
}

func managedPluginArgsAny() []any {
	args := make([]any, len(managedPluginArgs))
	for i, arg := range managedPluginArgs {
		args[i] = arg
	}
	return args
}

func validateManifestRaw(raw map[string]any) error {
	schema := pluginManifestValidationSchema()

	if err := rejectUnknownManifestFields(raw, schema); err != nil {
		return err
	}
	if err := validateConfigurationSchemaRaw(raw, schema.Properties["configuration_schema"]); err != nil {
		return err
	}

	resolved, err := schema.Resolve(nil)
	if err != nil {
		return fmt.Errorf("load manifest schema: %w", err)
	}
	if err := resolved.Validate(raw); err != nil {
		return fmt.Errorf("validate manifest schema: %w", err)
	}

	return nil
}

func rejectUnknownManifestFields(raw map[string]any, schema *jsonschema.Schema) error {
	for key := range raw {
		if _, ok := schema.Properties[key]; ok {
			continue
		}
		if suggestion, ok := nearestSchemaProperty(key, schema.Properties); ok {
			return fmt.Errorf("unknown manifest field %q; did you mean %q?", key, suggestion)
		}
		return fmt.Errorf("unknown manifest field %q", key)
	}
	return nil
}

func validateConfigurationSchemaRaw(raw map[string]any, schema *jsonschema.Schema) error {
	value, ok := raw["configuration_schema"]
	if !ok || value == nil {
		return nil
	}
	if schema == nil {
		return nil
	}

	configurationSchema, ok := value.(map[string]any)
	if !ok {
		return errors.New("configuration_schema must be an object")
	}

	if err := rejectUnknownObjectFields("configuration_schema", configurationSchema, schema.Properties); err != nil {
		return err
	}

	fieldsValue, ok := configurationSchema["fields"]
	if !ok || fieldsValue == nil {
		return nil
	}

	fields, ok := fieldsValue.([]any)
	if !ok {
		return errors.New("configuration_schema.fields must be an array")
	}

	fieldsSchema := schema.Properties["fields"]
	if fieldsSchema == nil || fieldsSchema.Items == nil {
		return nil
	}

	for index, fieldValue := range fields {
		field, ok := fieldValue.(map[string]any)
		if !ok {
			return fmt.Errorf("configuration_schema.fields[%d] must be an object", index)
		}
		if err := validateConfigurationFieldRaw(index, field, fieldsSchema.Items); err != nil {
			return err
		}
	}

	return nil
}

func validateConfigurationFieldRaw(index int, field map[string]any, itemSchema *jsonschema.Schema) error {
	path := fmt.Sprintf("configuration_schema.fields[%d]", index)
	branch, allowedTypes, ok, err := schemaBranchForDiscriminator(itemSchema, field["type"])
	if err != nil {
		return fmt.Errorf("%s.%w", path, err)
	}
	if !ok {
		return fmt.Errorf("%s.type must be one of: %s", path, strings.Join(allowedTypes, ", "))
	}

	for _, key := range branch.Required {
		if _, ok := field[key]; !ok {
			return fmt.Errorf("%s.%s is required", path, key)
		}
	}

	if err := rejectUnknownObjectFields(path, field, branch.Properties); err != nil {
		return err
	}

	for key, value := range field {
		propertySchema := branch.Properties[key]
		if propertySchema == nil {
			continue
		}
		if err := validateSimpleSchemaValue(value, propertySchema); err != nil {
			return fmt.Errorf("%s.%s %w", path, key, err)
		}
	}

	return nil
}

func rejectUnknownObjectFields(path string, object map[string]any, properties map[string]*jsonschema.Schema) error {
	for key := range object {
		if _, ok := properties[key]; ok {
			continue
		}
		if suggestion, ok := nearestSchemaProperty(key, properties); ok {
			return fmt.Errorf("%s.%s is not a valid field; did you mean %q?", path, key, suggestion)
		}
		return fmt.Errorf("%s.%s is not a valid field", path, key)
	}
	return nil
}

func schemaBranchForDiscriminator(schema *jsonschema.Schema, discriminator any) (*jsonschema.Schema, []string, bool, error) {
	branches := configurationFieldBranches(schema)
	allowedTypes := schemaDiscriminatorValues(branches)
	sort.Strings(allowedTypes)

	if _, ok := discriminator.(string); !ok {
		return nil, allowedTypes, false, errors.New("type must be a string")
	}

	for _, branch := range branches {
		typeSchema := branch.Properties["type"]
		if typeSchema == nil || typeSchema.Const == nil {
			continue
		}
		if fmt.Sprint(*typeSchema.Const) == discriminator {
			return branch, allowedTypes, true, nil
		}
	}

	return nil, allowedTypes, false, nil
}

func configurationFieldBranches(schema *jsonschema.Schema) []*jsonschema.Schema {
	if len(schema.OneOf) > 0 {
		return schema.OneOf
	}

	common := []*jsonschema.Schema{}
	branches := []*jsonschema.Schema{}
	for _, item := range schema.AllOf {
		if len(item.OneOf) > 0 {
			for _, branch := range item.OneOf {
				branches = append(branches, mergeSchemaBranch(append(common, branch)...))
			}
			continue
		}
		common = append(common, item)
	}

	return branches
}

func mergeSchemaBranch(parts ...*jsonschema.Schema) *jsonschema.Schema {
	out := &jsonschema.Schema{
		Properties: map[string]*jsonschema.Schema{},
	}
	seenRequired := map[string]struct{}{}

	for _, part := range parts {
		if part == nil {
			continue
		}
		for key, property := range part.Properties {
			out.Properties[key] = property
		}
		for _, required := range part.Required {
			if _, ok := seenRequired[required]; ok {
				continue
			}
			out.Required = append(out.Required, required)
			seenRequired[required] = struct{}{}
		}
	}

	return out
}

func schemaDiscriminatorValues(branches []*jsonschema.Schema) []string {
	values := []string{}
	for _, branch := range branches {
		typeSchema := branch.Properties["type"]
		if typeSchema == nil || typeSchema.Const == nil {
			continue
		}
		values = append(values, fmt.Sprint(*typeSchema.Const))
	}
	return values
}

func validateSimpleSchemaValue(value any, schema *jsonschema.Schema) error {
	if schema.Const != nil && fmt.Sprint(*schema.Const) != fmt.Sprint(value) {
		return fmt.Errorf("must be %q", fmt.Sprint(*schema.Const))
	}

	switch schema.Type {
	case "", "object", "array":
		return nil
	case "string":
		if _, ok := value.(string); !ok {
			return errors.New("must be a string")
		}
	case "number":
		if !isManifestNumber(value) {
			return errors.New("must be a number")
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return errors.New("must be a boolean")
		}
	}

	return nil
}

func isManifestNumber(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	default:
		return false
	}
}

func nearestSchemaProperty(key string, properties map[string]*jsonschema.Schema) (string, bool) {
	for candidate := range properties {
		if strings.HasPrefix(candidate, key) || strings.HasPrefix(key, candidate) {
			return candidate, true
		}
	}

	bestDistance := 0
	best := ""
	for candidate := range properties {
		distance := levenshteinDistance(key, candidate)
		if best == "" || distance < bestDistance {
			best = candidate
			bestDistance = distance
		}
	}
	if best == "" {
		return "", false
	}

	maxDistance := len(key) / 3
	if maxDistance < 2 {
		maxDistance = 2
	}
	return best, bestDistance <= maxDistance
}

func levenshteinDistance(a, b string) int {
	if a == b {
		return 0
	}
	if a == "" {
		return len(b)
	}
	if b == "" {
		return len(a)
	}

	previous := make([]int, len(b)+1)
	current := make([]int, len(b)+1)
	for j := range previous {
		previous[j] = j
	}

	for i := 1; i <= len(a); i++ {
		current[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			current[j] = minInt(
				current[j-1]+1,
				previous[j]+1,
				previous[j-1]+cost,
			)
		}
		previous, current = current, previous
	}

	return previous[len(b)]
}

func minInt(values ...int) int {
	best := values[0]
	for _, value := range values[1:] {
		if value < best {
			best = value
		}
	}
	return best
}

func pluginManifestToolInputSchema() *jsonschema.Schema {
	schema := libplugin.GetManifestSchema()

	schema.Description = "Storyden plugin manifest fields managed by the plugin builder. Use configuration_schema for plugin settings; configuration is not a valid field. Runtime launch fields are managed automatically."
	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}
	delete(schema.Properties, "command")
	delete(schema.Properties, "args")
	schema.Required = removeRequiredFields(schema.Required, "command")

	return schema
}

func pluginManifestValidationSchema() *jsonschema.Schema {
	schema := libplugin.GetManifestSchema()

	schema.Description = "Complete Storyden plugin manifest. Use configuration_schema for plugin settings; configuration is not a valid field."
	if schema.Properties == nil {
		schema.Properties = map[string]*jsonschema.Schema{}
	}

	return schema
}

func removeRequiredFields(required []string, fields ...string) []string {
	remove := map[string]struct{}{}
	for _, field := range fields {
		remove[field] = struct{}{}
	}

	out := make([]string, 0, len(required))
	for _, field := range required {
		if _, ok := remove[field]; ok {
			continue
		}
		out = append(out, field)
	}
	return out
}
