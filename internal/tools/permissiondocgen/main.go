package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
)

type PermissionDoc struct {
	ID          string
	Name        string
	Description string
}

type PermissionSchema struct {
	FrontmatterDescription string
	PageContent            string
	Permissions            []PermissionDoc
}

func main() {
	permissionSchema := flag.String("schema", filepath.Join("common", "permission.yaml"), "path to permission schema")
	openapi := flag.String("openapi", "openapi.yaml", "path to openapi schema")
	output := flag.String("out", filepath.Join("..", "home", "content", "docs", "introduction", "members", "permissions.mdx"), "path to generated permissions docs")
	flag.Parse()

	if err := run(*permissionSchema, *openapi, *output); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(permissionSchemaPath, openapiPath, outPath string) error {
	permissionSchema, err := loadPermissionSchema(permissionSchemaPath)
	if err != nil {
		return err
	}

	openapiPermissions, err := loadOpenAPIPermissions(openapiPath)
	if err != nil {
		return err
	}

	if err := validatePermissionsInSync(permissionSchema.Permissions, openapiPermissions); err != nil {
		return err
	}

	content := renderPermissionsMdx(
		permissionSchema,
		sourcePathForDoc(permissionSchemaPath),
		sourcePathForDoc(openapiPath),
	)

	return os.WriteFile(outPath, []byte(content), 0o644)
}

func loadPermissionSchema(path string) (PermissionSchema, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return PermissionSchema{}, err
	}

	var root map[string]any
	if err := yaml.Unmarshal(b, &root); err != nil {
		return PermissionSchema{}, fmt.Errorf("failed to parse permission schema: %w", err)
	}

	frontmatterDescription, pageContent, err := parsePageDescription(root["description"])
	if err != nil {
		return PermissionSchema{}, err
	}

	rawDefs, ok := asMap(root["$defs"])
	if !ok {
		return PermissionSchema{}, errors.New("permission schema has no $defs")
	}

	rawOneOf, ok := asList(root["oneOf"])
	if !ok {
		return PermissionSchema{}, errors.New("permission schema has no root oneOf")
	}

	seen := map[string]struct{}{}
	permissions := make([]PermissionDoc, 0, len(rawOneOf))

	for _, raw := range rawOneOf {
		entry, ok := asMap(raw)
		if !ok {
			return PermissionSchema{}, errors.New("invalid oneOf entry in permission schema")
		}

		ref := toString(entry["$ref"])
		if ref == "" {
			return PermissionSchema{}, errors.New("oneOf entry missing $ref in permission schema")
		}

		defName, ok := refDefName(ref)
		if !ok {
			return PermissionSchema{}, fmt.Errorf("oneOf entry has unsupported ref format: %s", ref)
		}

		defRaw, ok := asMap(rawDefs[defName])
		if !ok {
			return PermissionSchema{}, fmt.Errorf("permission schema missing $defs.%s", defName)
		}

		id := strings.TrimSpace(toString(defRaw["const"]))
		if id == "" {
			id = defName
		}
		if id != defName {
			return PermissionSchema{}, fmt.Errorf("permission schema mismatch: $defs.%s const is %s", defName, id)
		}
		if _, exists := seen[id]; exists {
			return PermissionSchema{}, fmt.Errorf("duplicate permission in oneOf list: %s", id)
		}
		seen[id] = struct{}{}

		description := normalizeDescription(defRaw["description"])
		if description == "" {
			return PermissionSchema{}, fmt.Errorf("permission %s has no description", id)
		}

		name := normalizeDescription(defRaw["title"])
		if name == "" {
			name = humanisePermissionName(id)
		}

		permissions = append(permissions, PermissionDoc{
			ID:          id,
			Name:        name,
			Description: description,
		})
	}

	return PermissionSchema{
		FrontmatterDescription: frontmatterDescription,
		PageContent:            pageContent,
		Permissions:            permissions,
	}, nil
}

func parsePageDescription(v any) (string, string, error) {
	description := strings.ReplaceAll(toString(v), "\r\n", "\n")
	if strings.TrimSpace(description) == "" {
		return "", "", errors.New("permission schema root description is empty")
	}

	lines := strings.Split(description, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			start = i
			break
		}
	}
	if start == -1 {
		return "", "", errors.New("permission schema root description is empty")
	}

	frontmatterDescription := strings.TrimSpace(lines[start])
	pageContent := strings.TrimSpace(strings.Join(lines[start+1:], "\n"))
	if pageContent == "" {
		return "", "", errors.New("permission schema root description must include page content after the first line")
	}

	return frontmatterDescription, pageContent, nil
}

func loadOpenAPIPermissions(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var root map[string]any
	if err := yaml.Unmarshal(b, &root); err != nil {
		return nil, fmt.Errorf("failed to parse openapi schema: %w", err)
	}

	components, ok := asMap(root["components"])
	if !ok {
		return nil, errors.New("openapi schema missing components")
	}

	schemas, ok := asMap(components["schemas"])
	if !ok {
		return nil, errors.New("openapi schema missing components.schemas")
	}

	var permissionSchema map[string]any
	if v, ok := asMap(schemas["Permissions"]); ok {
		permissionSchema = v
	} else if v, ok := asMap(schemas["Permission"]); ok {
		permissionSchema = v
	} else {
		return nil, errors.New("openapi schema missing components.schemas.Permission(s)")
	}

	rawEnum, ok := asList(permissionSchema["enum"])
	if !ok {
		return nil, errors.New("openapi permission schema has no enum")
	}

	seen := map[string]struct{}{}
	out := make([]string, 0, len(rawEnum))
	for _, raw := range rawEnum {
		value := strings.TrimSpace(toString(raw))
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			return nil, fmt.Errorf("duplicate permission in openapi enum: %s", value)
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}

	return out, nil
}

func validatePermissionsInSync(schemaPerms []PermissionDoc, openapiPerms []string) error {
	schemaSet := map[string]struct{}{}
	for _, p := range schemaPerms {
		schemaSet[p.ID] = struct{}{}
	}

	openapiSet := map[string]struct{}{}
	for _, p := range openapiPerms {
		openapiSet[p] = struct{}{}
	}

	var missingInOpenAPI []string
	for _, p := range schemaPerms {
		if _, ok := openapiSet[p.ID]; !ok {
			missingInOpenAPI = append(missingInOpenAPI, p.ID)
		}
	}

	var missingInSchema []string
	for _, p := range openapiPerms {
		if _, ok := schemaSet[p]; !ok {
			missingInSchema = append(missingInSchema, p)
		}
	}

	if len(missingInOpenAPI) == 0 && len(missingInSchema) == 0 {
		return nil
	}

	sort.Strings(missingInOpenAPI)
	sort.Strings(missingInSchema)

	return fmt.Errorf("permission schema and openapi are out of sync; missing in openapi: [%s], missing in schema: [%s]",
		strings.Join(missingInOpenAPI, ", "),
		strings.Join(missingInSchema, ", "),
	)
}

func renderPermissionsMdx(permissionSchema PermissionSchema, schemaPath, openapiPath string) string {
	var sb strings.Builder

	sb.WriteString(`---
title: Permissions
description: ` + frontmatterString(permissionSchema.FrontmatterDescription) + `
---

This page is generated from ` + "`" + schemaPath + "`" + ` and validated against ` + "`" + openapiPath + "`" + `.

` + permissionSchema.PageContent + `

<Callout type="info">
  Generated by ` + "`internal/tools/permissiondocgen`" + `. Edit schema, not this page.
</Callout>

| Name | ID | Description |
| --- | --- | --- |
`)

	for _, permission := range permissionSchema.Permissions {
		sb.WriteString("| **" + escapeMarkdownTableCell(permission.Name) + "** | `" + escapeMarkdownTableCell(permission.ID) + "` | " + escapeMarkdownTableCell(permission.Description) + " |\n")
	}

	return sb.String()
}

func frontmatterString(v string) string {
	return strconv.Quote(strings.TrimSpace(v))
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

func refDefName(ref string) (string, bool) {
	const prefix = "#/$defs/"
	if !strings.HasPrefix(ref, prefix) {
		return "", false
	}
	name := strings.TrimSpace(strings.TrimPrefix(ref, prefix))
	if name == "" {
		return "", false
	}
	return name, true
}

func humanisePermissionName(permission string) string {
	parts := strings.Split(permission, "_")
	for i, part := range parts {
		part = strings.ToLower(strings.TrimSpace(part))
		if part == "" {
			parts[i] = part
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}

func normalizeDescription(v any) string {
	return strings.TrimSpace(strings.ReplaceAll(toString(v), "\n", " "))
}

func escapeMarkdownTableCell(v string) string {
	v = strings.TrimSpace(strings.ReplaceAll(v, "\n", " "))
	if v == "" {
		return "—"
	}

	v = strings.ReplaceAll(v, "\\", "\\\\")
	v = strings.ReplaceAll(v, "|", "\\|")
	v = strings.ReplaceAll(v, "`", "\\`")
	return v
}

func asMap(v any) (map[string]any, bool) {
	switch t := v.(type) {
	case map[string]any:
		return t, true
	case map[any]any:
		out := make(map[string]any, len(t))
		for key, value := range t {
			out[toString(key)] = value
		}
		return out, true
	default:
		return nil, false
	}
}

func asList(v any) ([]any, bool) {
	t, ok := v.([]any)
	return t, ok
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		return fmt.Sprint(t)
	}
}
