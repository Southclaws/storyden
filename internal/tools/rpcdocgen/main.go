package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/goccy/go-yaml"
)

type MethodDoc struct {
	Method              string
	Slug                string
	DirectionSlug       string
	DirectionLabel      string
	RequestType         string
	ResponseType        string
	RequestDescription  string
	ResponseDescription string
	SourceFile          string
	RequestRequired     []string
	ResponseRequired    []string
}

type Direction struct {
	Slug   string
	Label  string
	Schema string
}

type ManifestField struct {
	Name        string
	Type        string
	Ref         string
	Example     any
	Required    bool
	Description string
}

type ManifestObject struct {
	Name        string
	Description string
	Fields      []ManifestField
}

type EventDoc struct {
	Name        string
	Description string
	Fields      []EventField
}

type EventField struct {
	Name        string
	Type        string
	Format      string
	Ref         string
	Const       string
	Required    bool
	Description string
}

var methodTemplate = template.Must(template.New("method").Parse(`---
title: {{ .Method }}
description: {{ .DirectionLabel }} RPC method.
---

` + "`" + `{{ .Method }}` + "`" + ` is a **{{ .DirectionLabel }}** method.

<Callout type="info">
  This page is generated from ` + "`" + `{{ .SourceFile }}` + "`" + `.
</Callout>

## Types

- Request type: ` + "`" + `{{ .RequestType }}` + "`" + `
- Response type: ` + "`" + `{{ .ResponseType }}` + "`" + `

## Request contract

{{ if .RequestDescription -}}
{{ .RequestDescription }}
{{ else -}}
No request description was provided in the schema.
{{ end }}

{{ if .RequestRequired }}Required fields:
{{ range .RequestRequired }}- ` + "`" + `{{ . }}` + "`" + `
{{ end }}
{{ else }}No explicit required fields are defined on this method-specific request object.
{{ end }}

## Response contract

{{ if .ResponseDescription -}}
{{ .ResponseDescription }}
{{ else -}}
No response description was provided in the schema.
{{ end }}

{{ if .ResponseRequired }}Required fields:
{{ range .ResponseRequired }}- ` + "`" + `{{ . }}` + "`" + `
{{ end }}
{{ else }}No explicit required fields are defined on this method-specific response object.
{{ end }}

## Notes

- JSON-RPC base request/response envelope types are defined in ` + "`" + `api/rpc/rpc.yaml` + "`" + `.
- Union wiring between host/plugin method sets is defined in ` + "`" + `api/plugin.yaml` + "`" + `.
`))

var directionIndexTemplate = template.Must(template.New("direction").Parse(`---
title: {{ .Label }}
description: Generated method reference pages for {{ .Label }} RPC.
---

{{ .Label }} methods define one side of Storyden's plugin RPC transport.

## Methods

{{ range .Methods -}}
- [{{ .Method }}](/docs/extending/api/{{ .DirectionSlug }}/{{ .Slug }})
{{ end }}
`))

func main() {
	schemaRoot := flag.String("schema", filepath.Join("rpc"), "path to RPC schema root")
	outputRoot := flag.String("out", filepath.Join("..", "home", "content", "docs", "extending", "api"), "path to docs output root")
	pluginSchema := flag.String("plugin-schema", filepath.Join("plugin.yaml"), "path to plugin manifest schema")
	manifestOut := flag.String("manifest-out", filepath.Join("..", "home", "content", "docs", "extending", "manifest.mdx"), "path to generated manifest docs output")
	eventsSchema := flag.String("events-schema", filepath.Join("common", "events.yaml"), "path to events schema")
	eventsOut := flag.String("events-out", filepath.Join("..", "home", "content", "docs", "extending", "api", "events.mdx"), "path to generated events docs output")
	flag.Parse()

	if err := run(*schemaRoot, *outputRoot, *pluginSchema, *manifestOut, *eventsSchema, *eventsOut); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(schemaRoot, outputRoot, pluginSchema, manifestOut, eventsSchema, eventsOut string) error {
	directions := []Direction{
		{Slug: "host-to-plugin", Label: "Host to Plugin", Schema: filepath.Join(schemaRoot, "host-to-plugin")},
		{Slug: "plugin-to-host", Label: "Plugin to Host", Schema: filepath.Join(schemaRoot, "plugin-to-host")},
	}

	if err := os.MkdirAll(outputRoot, 0o755); err != nil {
		return err
	}

	for _, direction := range directions {
		methods, err := loadDirection(direction)
		if err != nil {
			return err
		}

		dirOut := filepath.Join(outputRoot, direction.Slug)
		if err := os.MkdirAll(dirOut, 0o755); err != nil {
			return err
		}

		for _, method := range methods {
			body, err := executeTemplate(methodTemplate, method)
			if err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(dirOut, method.Slug+".mdx"), body, 0o644); err != nil {
				return err
			}
		}

		directionIndex, err := executeTemplate(directionIndexTemplate, map[string]any{
			"Label":   direction.Label,
			"Methods": methods,
		})
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dirOut, "index.mdx"), directionIndex, 0o644); err != nil {
			return err
		}

		pages := make([]string, 0, len(methods))
		for _, method := range methods {
			pages = append(pages, method.Slug)
		}

		if err := writeJSON(filepath.Join(dirOut, "meta.json"), map[string]any{"pages": pages}); err != nil {
			return err
		}
	}

	if err := generateManifestDoc(pluginSchema, manifestOut); err != nil {
		return err
	}
	if err := generateEventsDoc(eventsSchema, eventsOut); err != nil {
		return err
	}

	return nil
}

func generateEventsDoc(schemaPath, outputPath string) error {
	b, err := os.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	var root map[string]any
	if err := yaml.Unmarshal(b, &root); err != nil {
		return fmt.Errorf("failed to parse events schema: %w", err)
	}

	defs, ok := asMap(root["$defs"])
	if !ok {
		return fmt.Errorf("events schema has no $defs")
	}

	events, err := parseEventDocs(defs)
	if err != nil {
		return err
	}

	content := renderEventsMdx(events, sourcePathForDoc(schemaPath))
	return os.WriteFile(outputPath, []byte(content), 0o644)
}

func parseEventDocs(defs map[string]any) ([]EventDoc, error) {
	eventEnumDef, ok := asMap(defs["Event"])
	if !ok {
		return nil, fmt.Errorf("events schema missing $defs.Event")
	}

	enumValues, ok := asList(eventEnumDef["enum"])
	if !ok {
		return nil, fmt.Errorf("events schema $defs.Event has no enum list")
	}

	events := make([]EventDoc, 0, len(enumValues))
	for _, raw := range enumValues {
		name := strings.TrimSpace(toString(raw))
		if name == "" {
			continue
		}

		def, ok := asMap(defs[name])
		if !ok {
			continue
		}

		fields := parseEventFields(def)
		sortEventFields(fields)

		events = append(events, EventDoc{
			Name:        name,
			Description: normalizeDescription(def["description"]),
			Fields:      fields,
		})
	}

	return events, nil
}

func parseEventFields(def map[string]any) []EventField {
	props, ok := asMap(def["properties"])
	if !ok {
		return nil
	}

	requiredSet := map[string]struct{}{}
	for _, name := range parseRequiredValues(def["required"]) {
		requiredSet[name] = struct{}{}
	}

	fields := make([]EventField, 0, len(props))
	for name, raw := range props {
		prop, ok := asMap(raw)
		if !ok {
			continue
		}

		_, required := requiredSet[name]
		fields = append(fields, EventField{
			Name:        name,
			Type:        inferEventFieldType(prop),
			Format:      toString(prop["format"]),
			Ref:         toString(prop["$ref"]),
			Const:       toString(prop["const"]),
			Required:    required,
			Description: normalizeDescription(prop["description"]),
		})
	}

	return fields
}

func inferEventFieldType(prop map[string]any) string {
	if ref := toString(prop["$ref"]); ref != "" {
		if refType := refTypeName(ref); refType != "" {
			return refType
		}
		return "object"
	}
	if t := toString(prop["type"]); t != "" {
		return t
	}
	if c := toString(prop["const"]); c != "" {
		return "const"
	}
	return "unknown"
}

func refTypeName(ref string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}

	if i := strings.LastIndex(ref, "#/$defs/"); i >= 0 {
		name := strings.TrimSpace(ref[i+len("#/$defs/"):])
		if name != "" {
			return name
		}
	}

	if i := strings.LastIndex(ref, "/"); i >= 0 && i+1 < len(ref) {
		name := strings.TrimSpace(ref[i+1:])
		if name != "" {
			return name
		}
	}

	return ref
}

func sortEventFields(fields []EventField) {
	sort.Slice(fields, func(i, j int) bool {
		if fields[i].Name == "event" {
			return true
		}
		if fields[j].Name == "event" {
			return false
		}
		return fields[i].Name < fields[j].Name
	})
}

func renderEventsMdx(events []EventDoc, schemaPath string) string {
	var sb strings.Builder

	sb.WriteString(`---
title: Events
description: Generated reference for Storyden plugin event payloads.
---

This page is generated from ` + "`" + schemaPath + "`" + `.

<Callout type="info">
  Generated by ` + "`internal/tools/rpcdocgen`" + `. Edit schema, not this page.
</Callout>

Related pages:
- [Manifest](/docs/extending/manifest#events_consumed)
- [API](/docs/extending/api)
- [Host to Plugin -> event](/docs/extending/api/host-to-plugin/event)

## Events

`)

	for _, event := range events {
		sb.WriteString("### " + event.Name + "\n\n")
		if event.Description != "" {
			sb.WriteString(event.Description + "\n\n")
		} else {
			sb.WriteString("Payload schema for `" + event.Name + "`.\n\n")
		}

		if len(event.Fields) == 0 {
			sb.WriteString("No payload fields defined in schema.\n\n")
			continue
		}

		sb.WriteString("Payload fields:\n\n")
		sb.WriteString("| Field | Type | Required | Description |\n")
		sb.WriteString("| --- | --- | --- | --- |\n")
		for _, field := range event.Fields {
			typeLabel := field.Type
			if field.Format != "" {
				typeLabel += " (" + field.Format + ")"
			}
			if field.Const != "" {
				typeLabel += ", const " + field.Const
			}

			description := field.Description
			if description == "" {
				description = "-"
			}

			sb.WriteString("| `" + escapeMarkdownTableCell(field.Name) + "` | `" + escapeMarkdownTableCell(typeLabel) + "` | " + yesNo(field.Required) + " | " + escapeMarkdownTableCell(description) + " |\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func escapeMarkdownTableCell(v string) string {
	v = strings.ReplaceAll(v, "\n", " ")
	v = strings.TrimSpace(v)
	if v == "" {
		return "—"
	}

	v = strings.ReplaceAll(v, "\\", "\\\\")
	v = strings.ReplaceAll(v, "|", "\\|")
	v = strings.ReplaceAll(v, "`", "\\`")
	return v
}

func generateManifestDoc(schemaPath, outputPath string) error {
	b, err := os.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	var root map[string]any
	if err := yaml.Unmarshal(b, &root); err != nil {
		return fmt.Errorf("failed to parse plugin schema: %w", err)
	}

	required := parseRequiredValues(root["required"])
	requiredSet := make(map[string]struct{}, len(required))
	for _, field := range required {
		requiredSet[field] = struct{}{}
	}

	props, ok := asMap(root["properties"])
	if !ok {
		return fmt.Errorf("plugin schema has no properties")
	}

	topFields := parseManifestFields(props, requiredSet)
	sortManifestTopFields(topFields)

	defs, _ := asMap(root["$defs"])
	objects := parseManifestObjects(defs)
	exampleYAML := renderManifestExampleYAML(topFields, defs)

	content := renderManifestMdx(topFields, required, objects, exampleYAML)
	return os.WriteFile(outputPath, []byte(content), 0o644)
}

func parseManifestObjects(defs map[string]any) []ManifestObject {
	if len(defs) == 0 {
		return nil
	}

	objects := make([]ManifestObject, 0, len(defs))
	for name, raw := range defs {
		if !strings.HasPrefix(name, "Manifest") {
			continue
		}
		def, ok := asMap(raw)
		if !ok {
			continue
		}

		props, _ := asMap(def["properties"])
		required := parseRequiredValues(def["required"])
		requiredSet := make(map[string]struct{}, len(required))
		for _, field := range required {
			requiredSet[field] = struct{}{}
		}

		fields := parseManifestFields(props, requiredSet)
		sortManifestFields(fields)

		objects = append(objects, ManifestObject{
			Name:        name,
			Description: normalizeDescription(def["description"]),
			Fields:      fields,
		})
	}

	sort.Slice(objects, func(i, j int) bool {
		if objects[i].Name == "ManifestAccess" {
			return true
		}
		if objects[j].Name == "ManifestAccess" {
			return false
		}
		return objects[i].Name < objects[j].Name
	})

	return objects
}

func parseManifestFields(props map[string]any, requiredSet map[string]struct{}) []ManifestField {
	fields := make([]ManifestField, 0, len(props))
	for name, raw := range props {
		prop, ok := asMap(raw)
		if !ok {
			continue
		}
		_, required := requiredSet[name]
		fields = append(fields, ManifestField{
			Name:        name,
			Type:        inferManifestFieldType(prop),
			Ref:         toString(prop["$ref"]),
			Example:     normalizeExampleValue(prop["example"]),
			Required:    required,
			Description: normalizeDescription(prop["description"]),
		})
	}
	return fields
}

func sortManifestTopFields(fields []ManifestField) {
	priority := map[string]int{
		"id":                   0,
		"name":                 1,
		"author":               2,
		"description":          3,
		"version":              4,
		"command":              5,
		"args":                 6,
		"events_consumed":      7,
		"access":               8,
		"configuration_schema": 9,
	}

	sort.Slice(fields, func(i, j int) bool {
		pi, iok := priority[fields[i].Name]
		pj, jok := priority[fields[j].Name]
		switch {
		case iok && jok:
			return pi < pj
		case iok:
			return true
		case jok:
			return false
		default:
			return fields[i].Name < fields[j].Name
		}
	})
}

func sortManifestFields(fields []ManifestField) {
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
}

func inferManifestFieldType(prop map[string]any) string {
	if ref := toString(prop["$ref"]); ref != "" {
		return "object"
	}

	switch t := toString(prop["type"]); t {
	case "array":
		items, ok := asMap(prop["items"])
		if !ok {
			return "array"
		}
		if itemRef := toString(items["$ref"]); itemRef != "" {
			return "array<object>"
		}
		if itemType := toString(items["type"]); itemType != "" {
			return "array<" + itemType + ">"
		}
		return "array"
	case "":
		if _, ok := prop["additionalProperties"]; ok {
			return "object"
		}
		return "unknown"
	default:
		return t
	}
}

func renderManifestMdx(fields []ManifestField, required []string, objects []ManifestObject, exampleYAML string) string {
	var sb strings.Builder

	sb.WriteString(`---
title: Manifest
description: Generated schema reference for plugin manifests.
---

This page is generated from ` + "`api/plugin.yaml`" + `.

The manifest is the source of truth for both supervised and external plugins.

<Callout type="info">
  Generated by ` + "`internal/tools/rpcdocgen`" + `. Edit schema, not this page.
</Callout>

Related pages:
- [Plugin Model](/docs/extending/model)
- [Capabilities and Limits](/docs/extending/capabilities)
- [Security](/docs/extending/security)
- [API](/docs/extending/api)

## Required fields

`)

	for _, name := range required {
		sb.WriteString(fmt.Sprintf("- [%s](#%s)\n", name, name))
	}

	if containsString(required, "command") {
		sb.WriteString("\n<Callout type=\"info\">\n")
		sb.WriteString("  `command` and `args` are runtime entrypoint fields used by Supervised plugins.\n")
		sb.WriteString("  External plugins can provide placeholder values.\n")
		sb.WriteString("</Callout>\n")
	}

	if exampleYAML != "" {
		sb.WriteString("\n## Example YAML\n\n")
		sb.WriteString("```yaml\n")
		sb.WriteString(exampleYAML)
		sb.WriteString("\n```\n")
	}

	sb.WriteString("\n## Top-level fields\n\n")
	for _, field := range fields {
		appendManifestField(&sb, field, 3)
	}

	if len(objects) > 0 {
		sb.WriteString("## Nested objects\n\n")
		sb.WriteString("These objects are referenced by top-level fields.\n\n")

		for _, object := range objects {
			sb.WriteString(fmt.Sprintf("## %s\n\n", object.Name))
			if object.Description != "" {
				sb.WriteString(object.Description + "\n\n")
			}
			for _, field := range object.Fields {
				appendManifestField(&sb, field, 3)
			}
		}
	}

	return sb.String()
}

func appendManifestField(sb *strings.Builder, field ManifestField, level int) {
	sb.WriteString(strings.Repeat("#", level) + " " + field.Name + "\n\n")
	sb.WriteString(fmt.Sprintf("- Type: `%s`\n", field.Type))
	sb.WriteString(fmt.Sprintf("- Required: %s\n", yesNo(field.Required)))
	if field.Ref != "" {
		sb.WriteString(fmt.Sprintf("- Schema ref: `%s`\n", field.Ref))
	}
	if field.Example != nil {
		sb.WriteString(fmt.Sprintf("- Example: `%s`\n", shortExample(field.Example)))
	}
	sb.WriteString("\n")
	if field.Description != "" {
		sb.WriteString(field.Description + "\n\n")
	} else {
		sb.WriteString("No description provided in schema.\n\n")
	}
}

func parseRequiredValues(v any) []string {
	list, ok := asList(v)
	if !ok {
		return nil
	}

	out := make([]string, 0, len(list))
	for _, entry := range list {
		if s := strings.TrimSpace(toString(entry)); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func yesNo(v bool) string {
	if v {
		return "Yes"
	}
	return "No"
}

func renderManifestExampleYAML(fields []ManifestField, defs map[string]any) string {
	var sb strings.Builder

	for _, field := range fields {
		if !shouldIncludeManifestExampleField(field) {
			continue
		}

		value, ok := manifestExampleValue(field, defs)
		if !ok {
			continue
		}

		writeYAMLPair(&sb, field.Name, value, 0)
	}

	return strings.TrimRight(sb.String(), "\n")
}

func shouldIncludeManifestExampleField(field ManifestField) bool {
	return field.Required || field.Example != nil
}

func manifestExampleValue(field ManifestField, defs map[string]any) (any, bool) {
	if field.Example != nil {
		if field.Name == "id" {
			if s, ok := field.Example.(string); ok && strings.Contains(strings.ToLower(s), "external") {
				return placeholderForField(field.Name), true
			}
		}
		return field.Example, true
	}

	if field.Ref != "" {
		if strings.HasPrefix(field.Ref, "#/$defs/") {
			name := strings.TrimPrefix(field.Ref, "#/$defs/")
			def, ok := asMap(defs[name])
			if !ok {
				return nil, false
			}
			return manifestObjectExample(def, defs)
		}

		if field.Ref == "common/plugin-configuration.yaml" {
			return map[string]any{
				"fields": []any{
					map[string]any{
						"id":          "example_field",
						"label":       "Example Field",
						"description": "Example configuration value.",
						"type":        "string",
					},
				},
			}, true
		}
	}

	switch field.Type {
	case "string":
		return placeholderForField(field.Name), true
	case "number":
		return 0, true
	case "boolean":
		return false, true
	default:
		if strings.HasPrefix(field.Type, "array") {
			return []any{}, true
		}
		if field.Type == "object" {
			return map[string]any{}, true
		}
		return nil, false
	}
}

func manifestObjectExample(def map[string]any, defs map[string]any) (any, bool) {
	if example := normalizeExampleValue(def["example"]); example != nil {
		return example, true
	}

	props, ok := asMap(def["properties"])
	if !ok {
		return map[string]any{}, true
	}

	required := parseRequiredValues(def["required"])
	requiredSet := make(map[string]struct{}, len(required))
	for _, name := range required {
		requiredSet[name] = struct{}{}
	}

	fields := parseManifestFields(props, requiredSet)
	sortManifestFields(fields)

	out := map[string]any{}
	for _, field := range fields {
		if !field.Required && field.Example == nil {
			continue
		}
		value, ok := manifestExampleValue(field, defs)
		if !ok {
			continue
		}
		out[field.Name] = value
	}

	return out, true
}

func writeYAMLPair(sb *strings.Builder, key string, value any, indent int) {
	writeIndent(sb, indent)
	sb.WriteString(key)
	sb.WriteString(":")

	switch v := value.(type) {
	case map[string]any:
		if len(v) == 0 {
			sb.WriteString(" {}\n")
			return
		}
		sb.WriteString("\n")
		writeYAMLMap(sb, v, indent+2)
	case []any:
		if len(v) == 0 {
			sb.WriteString(" []\n")
			return
		}
		sb.WriteString("\n")
		writeYAMLList(sb, v, indent+2)
	default:
		sb.WriteString(" ")
		sb.WriteString(yamlScalar(v))
		sb.WriteString("\n")
	}
}

func writeYAMLMap(sb *strings.Builder, m map[string]any, indent int) {
	keys := sortedMapKeys(m)
	for _, key := range keys {
		writeYAMLPair(sb, key, m[key], indent)
	}
}

func writeYAMLList(sb *strings.Builder, list []any, indent int) {
	for _, item := range list {
		writeIndent(sb, indent)
		sb.WriteString("-")

		switch v := item.(type) {
		case map[string]any:
			if len(v) == 0 {
				sb.WriteString(" {}\n")
				continue
			}
			keys := sortedMapKeys(v)
			firstKey := keys[0]
			firstVal := v[firstKey]

			sb.WriteString(" ")
			sb.WriteString(firstKey)
			sb.WriteString(":")

			switch fv := firstVal.(type) {
			case map[string]any:
				if len(fv) == 0 {
					sb.WriteString(" {}\n")
				} else {
					sb.WriteString("\n")
					writeYAMLMap(sb, fv, indent+4)
				}
			case []any:
				if len(fv) == 0 {
					sb.WriteString(" []\n")
				} else {
					sb.WriteString("\n")
					writeYAMLList(sb, fv, indent+4)
				}
			default:
				sb.WriteString(" ")
				sb.WriteString(yamlScalar(fv))
				sb.WriteString("\n")
			}

			for _, key := range keys[1:] {
				writeYAMLPair(sb, key, v[key], indent+2)
			}
		case []any:
			if len(v) == 0 {
				sb.WriteString(" []\n")
				continue
			}
			sb.WriteString("\n")
			writeYAMLList(sb, v, indent+2)
		default:
			sb.WriteString(" ")
			sb.WriteString(yamlScalar(v))
			sb.WriteString("\n")
		}
	}
}

func sortedMapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func writeIndent(sb *strings.Builder, indent int) {
	for range indent {
		sb.WriteByte(' ')
	}
}

func yamlScalar(v any) string {
	switch t := v.(type) {
	case nil:
		return "null"
	case string:
		return strconv.Quote(t)
	case bool:
		if t {
			return "true"
		}
		return "false"
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", t)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", t)
	case float32, float64:
		return fmt.Sprintf("%v", t)
	default:
		return strconv.Quote(fmt.Sprintf("%v", t))
	}
}

func shortExample(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []any:
		return fmt.Sprintf("array[%d]", len(t))
	case map[string]any:
		return "object"
	default:
		return fmt.Sprintf("%v", t)
	}
}

func placeholderForField(name string) string {
	switch name {
	case "id":
		return "my-plugin"
	case "name":
		return "My Plugin"
	case "author":
		return "you"
	case "description":
		return "Describe what this plugin does."
	case "version":
		return "0.1.0"
	case "command":
		return "./my-plugin"
	default:
		return "example"
	}
}

func normalizeExampleValue(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for key, value := range t {
			out[key] = normalizeExampleValue(value)
		}
		return out
	case map[any]any:
		out := make(map[string]any, len(t))
		for key, value := range t {
			out[toString(key)] = normalizeExampleValue(value)
		}
		return out
	case []any:
		out := make([]any, 0, len(t))
		for _, value := range t {
			out = append(out, normalizeExampleValue(value))
		}
		return out
	default:
		return t
	}
}

func loadDirection(direction Direction) ([]MethodDoc, error) {
	entries, err := os.ReadDir(direction.Schema)
	if err != nil {
		return nil, err
	}

	methods := make([]MethodDoc, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".yaml" {
			continue
		}

		schemaPath := filepath.Join(direction.Schema, entry.Name())
		doc, err := parseMethodFile(schemaPath, direction)
		if err != nil {
			return nil, err
		}
		methods = append(methods, doc)
	}

	sort.Slice(methods, func(i, j int) bool {
		return methods[i].Method < methods[j].Method
	})

	return methods, nil
}

func parseMethodFile(path string, direction Direction) (MethodDoc, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return MethodDoc{}, err
	}

	var root map[string]any
	if err := yaml.Unmarshal(b, &root); err != nil {
		return MethodDoc{}, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	defs, ok := asMap(root["$defs"])
	if !ok {
		return MethodDoc{}, fmt.Errorf("file %s has no $defs", path)
	}

	requestType, requestDef := findRequestDefinition(defs)
	if requestType == "" || requestDef == nil {
		return MethodDoc{}, fmt.Errorf("file %s has no RPCRequest* definition", path)
	}

	requestSuffix := strings.TrimPrefix(requestType, "RPCRequest")
	responseType, responseDef := findResponseDefinition(defs, requestSuffix)
	if responseType == "" || responseDef == nil {
		return MethodDoc{}, fmt.Errorf("file %s has no RPCResponse* definition", path)
	}

	methodName := findMethodConst(requestDef)
	if methodName == "" {
		methodName = findMethodConst(responseDef)
	}
	if methodName == "" {
		methodName = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	requestDesc := normalizeDescription(toString(requestDef["description"]))
	responseDesc := normalizeDescription(toString(responseDef["description"]))

	reqRequired := extractRequired(requestDef)
	respRequired := extractRequired(responseDef)

	repoPath := sourcePathForDoc(path)

	return MethodDoc{
		Method:              methodName,
		Slug:                methodName,
		DirectionSlug:       direction.Slug,
		DirectionLabel:      direction.Label,
		RequestType:         requestType,
		ResponseType:        responseType,
		RequestDescription:  requestDesc,
		ResponseDescription: responseDesc,
		SourceFile:          repoPath,
		RequestRequired:     reqRequired,
		ResponseRequired:    respRequired,
	}, nil
}

func sourcePathForDoc(path string) string {
	p := filepath.ToSlash(path)

	if idx := strings.Index(p, "/api/"); idx >= 0 {
		return strings.TrimPrefix(p[idx+1:], "/")
	}

	p = strings.TrimPrefix(p, "./")
	p = strings.TrimPrefix(p, "/")

	if strings.HasPrefix(p, "api/") {
		return p
	}

	return filepath.ToSlash(filepath.Join("api", p))
}

func findRequestDefinition(defs map[string]any) (string, map[string]any) {
	keys := make([]string, 0, len(defs))
	for k := range defs {
		if strings.HasPrefix(k, "RPCRequest") {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		v, ok := asMap(defs[k])
		if !ok {
			continue
		}
		if strings.HasSuffix(k, "Params") {
			continue
		}
		return k, v
	}
	return "", nil
}

func findResponseDefinition(defs map[string]any, suffix string) (string, map[string]any) {
	candidate := "RPCResponse" + suffix
	if v, ok := asMap(defs[candidate]); ok {
		return candidate, v
	}

	keys := make([]string, 0, len(defs))
	for k := range defs {
		if strings.HasPrefix(k, "RPCResponse") {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		if strings.HasSuffix(k, "Result") || strings.HasSuffix(k, "Error") {
			continue
		}
		v, ok := asMap(defs[k])
		if !ok {
			continue
		}
		return k, v
	}

	return "", nil
}

func findMethodConst(def map[string]any) string {
	return walkForMethodConst(def)
}

func walkForMethodConst(v any) string {
	m, ok := asMap(v)
	if ok {
		if props, ok := asMap(m["properties"]); ok {
			if method, ok := asMap(props["method"]); ok {
				if s := toString(method["const"]); s != "" {
					return s
				}
			}
		}

		for _, value := range m {
			if found := walkForMethodConst(value); found != "" {
				return found
			}
		}
	}

	if list, ok := asList(v); ok {
		for _, item := range list {
			if found := walkForMethodConst(item); found != "" {
				return found
			}
		}
	}

	return ""
}

func extractRequired(def map[string]any) []string {
	requiredSet := map[string]struct{}{}

	if req, ok := asList(def["required"]); ok {
		for _, entry := range req {
			if s := toString(entry); s != "" {
				requiredSet[s] = struct{}{}
			}
		}
	}

	if allOf, ok := asList(def["allOf"]); ok {
		for _, part := range allOf {
			partMap, ok := asMap(part)
			if !ok {
				continue
			}
			if req, ok := asList(partMap["required"]); ok {
				for _, entry := range req {
					if s := toString(entry); s != "" {
						requiredSet[s] = struct{}{}
					}
				}
			}
		}
	}

	if len(requiredSet) == 0 {
		return nil
	}

	out := make([]string, 0, len(requiredSet))
	for key := range requiredSet {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

func executeTemplate(t *template.Template, data any) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeJSON(path string, value any) error {
	b, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0o644)
}

func normalizeDescription(v any) string {
	return strings.TrimSpace(toString(v))
}

func asMap(v any) (map[string]any, bool) {
	m, ok := v.(map[string]any)
	if ok {
		return m, true
	}
	return nil, false
}

func asList(v any) ([]any, bool) {
	list, ok := v.([]any)
	if ok {
		return list, true
	}
	return nil, false
}

func toString(v any) string {
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}
