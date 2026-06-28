package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/oasdiff/yaml"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

type options struct {
	inputFile  string
	outputFile string
	rootOnly   bool
}

type resolver struct {
	cache map[string]map[string]any
	stack map[string]struct{}
}

func run() error {
	opts, err := parseOptions(os.Args[1:])
	if err != nil {
		return err
	}

	fmt.Printf("Dereferencing %s -> %s\n", opts.inputFile, opts.outputFile)

	r := resolver{
		cache: map[string]map[string]any{},
		stack: map[string]struct{}{},
	}

	inputPath, err := filepath.Abs(opts.inputFile)
	if err != nil {
		return fmt.Errorf("resolve input path: %w", err)
	}

	root, err := r.loadDocument(inputPath)
	if err != nil {
		return err
	}

	if !opts.rootOnly {
		rootDefs := make(map[string]any)
		if existing, ok := root["definitions"].(map[string]any); ok {
			for k, v := range existing {
				rootDefs[k] = v
			}
		}

		data := root
		if err := dereferenceSchemaLegacy(data, filepath.Dir(inputPath), rootDefs); err != nil {
			return fmt.Errorf("dereference schema: %w", err)
		}
		data["definitions"] = rootDefs
		if err := hoistNestedDefinitionsLegacy(data, rootDefs, "root"); err != nil {
			return err
		}
		return writeSchema(opts.outputFile, data)
	}

	output, err := r.dereferenceRoot(root, inputPath)
	if err != nil {
		return fmt.Errorf("dereference schema: %w", err)
	}

	if err := stripNestedDefinitions(output, !opts.rootOnly, true); err != nil {
		return err
	}

	data, ok := output.(map[string]any)
	if !ok {
		return fmt.Errorf("schema root must be an object")
	}

	return writeSchema(opts.outputFile, data)
}

func writeSchema(outputFile string, data map[string]any) error {
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}

	if err := os.WriteFile(outputFile, encoded, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	fmt.Printf("✓ Successfully dereferenced schema to %s\n", outputFile)
	return nil
}

func parseOptions(args []string) (options, error) {
	opts := options{}
	filtered := []string{}
	for _, arg := range args {
		switch arg {
		case "--root-only":
			opts.rootOnly = true
		default:
			filtered = append(filtered, arg)
		}
	}

	if len(filtered) >= 2 {
		opts.inputFile = filtered[0]
		opts.outputFile = filtered[1]
		return opts, nil
	}
	if len(filtered) != 0 {
		return opts, fmt.Errorf("usage: schemaderef [--root-only] <input.yaml> <output.json>")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return opts, fmt.Errorf("get working directory: %w", err)
	}

	baseDir := "api"
	if filepath.Base(cwd) == "api" {
		baseDir = "."
	}

	opts.inputFile = filepath.Join(baseDir, "robots.yaml")
	opts.outputFile = filepath.Join(baseDir, "robots.json")
	return opts, nil
}

func (r *resolver) loadDocument(path string) (map[string]any, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve schema path %s: %w", path, err)
	}
	if doc, ok := r.cache[abs]; ok {
		return doc, nil
	}

	data, err := os.ReadFile(abs)
	if err != nil {
		return nil, fmt.Errorf("read schema %s: %w", path, err)
	}

	var doc map[string]any
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("unmarshal schema %s: %w", path, err)
	}

	r.cache[abs] = doc
	return doc, nil
}

func (r *resolver) dereferenceRoot(root map[string]any, currentFile string) (any, error) {
	output := copyMapWithoutDefinitions(root)
	return r.dereferenceNode(output, currentFile, "")
}

func (r *resolver) dereferenceNode(node any, currentFile string, path string) (any, error) {
	switch v := node.(type) {
	case map[string]any:
		if ref, ok := v["$ref"].(string); ok {
			resolved, err := r.resolveRef(ref, currentFile)
			if err != nil {
				return nil, err
			}
			return mergeRefSiblings(resolved, v), nil
		}

		out := make(map[string]any, len(v))
		for key, value := range v {
			resolved, err := r.dereferenceNode(value, currentFile, joinPath(path, key))
			if err != nil {
				return nil, err
			}
			out[key] = resolved
		}
		return out, nil
	case []any:
		out := make([]any, 0, len(v))
		for index, item := range v {
			resolved, err := r.dereferenceNode(item, currentFile, fmt.Sprintf("%s[%d]", path, index))
			if err != nil {
				return nil, err
			}
			out = append(out, resolved)
		}
		return out, nil
	default:
		return v, nil
	}
}

func (r *resolver) resolveRef(ref string, currentFile string) (any, error) {
	targetFile, fragment, err := resolveReferenceLocation(ref, currentFile)
	if err != nil {
		return nil, err
	}

	key := targetFile + "#" + fragment
	if _, ok := r.stack[key]; ok {
		return map[string]any{"$ref": ref}, nil
	}

	doc, err := r.loadDocument(targetFile)
	if err != nil {
		return nil, err
	}

	target, err := selectFragment(doc, fragment)
	if err != nil {
		return nil, fmt.Errorf("resolve ref %s: %w", ref, err)
	}

	r.stack[key] = struct{}{}
	defer delete(r.stack, key)

	return r.dereferenceNode(deepCopy(target), targetFile, "")
}

func resolveReferenceLocation(ref string, currentFile string) (string, string, error) {
	parsed, err := url.Parse(ref)
	if err != nil {
		return "", "", fmt.Errorf("parse ref %s: %w", ref, err)
	}
	if parsed.Scheme != "" || parsed.Host != "" {
		return "", "", fmt.Errorf("external URL refs are not supported: %s", ref)
	}

	fragment := parsed.Fragment
	refPath := parsed.Path
	if refPath == "" {
		return currentFile, fragment, nil
	}

	target := refPath
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(currentFile), target)
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return "", "", fmt.Errorf("resolve ref path %s: %w", ref, err)
	}
	return abs, fragment, nil
}

func selectFragment(doc map[string]any, fragment string) (any, error) {
	if fragment == "" {
		return doc, nil
	}
	if !strings.HasPrefix(fragment, "/") {
		return nil, fmt.Errorf("unsupported fragment %q", fragment)
	}

	var current any = doc
	for _, token := range strings.Split(strings.TrimPrefix(fragment, "/"), "/") {
		key := strings.ReplaceAll(strings.ReplaceAll(token, "~1", "/"), "~0", "~")
		object, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("fragment %q does not point to an object field", fragment)
		}
		value, ok := object[key]
		if !ok {
			return nil, fmt.Errorf("fragment %q field %q not found", fragment, key)
		}
		current = value
	}

	return current, nil
}

func stripNestedDefinitions(node any, preserveRoot bool, root bool) error {
	switch v := node.(type) {
	case map[string]any:
		if !(root && preserveRoot) {
			delete(v, "$schema")
			delete(v, "$defs")
			delete(v, "definitions")
		}
		for _, value := range v {
			if err := stripNestedDefinitions(value, preserveRoot, false); err != nil {
				return err
			}
		}
	case []any:
		for _, item := range v {
			if err := stripNestedDefinitions(item, preserveRoot, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func mergeRefSiblings(resolved any, refObject map[string]any) any {
	resolvedObject, ok := resolved.(map[string]any)
	if !ok {
		return resolved
	}

	out := make(map[string]any, len(resolvedObject)+len(refObject))
	for key, value := range resolvedObject {
		if key == "$schema" {
			continue
		}
		out[key] = value
	}
	for key, value := range refObject {
		if key == "$ref" {
			continue
		}
		out[key] = deepCopy(value)
	}
	return out
}

func copyMapWithoutDefinitions(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for key, value := range in {
		if key == "$defs" || key == "definitions" {
			continue
		}
		out[key] = deepCopy(value)
	}
	return out
}

func deepCopy(value any) any {
	switch v := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for key, item := range v {
			out[key] = deepCopy(item)
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, item := range v {
			out[i] = deepCopy(item)
		}
		return out
	default:
		return v
	}
}

func joinPath(parent, child string) string {
	if parent == "" {
		return child
	}
	return parent + "." + child
}

func dereferenceSchemaLegacy(data any, baseDir string, rootDefs map[string]any) error {
	switch v := data.(type) {
	case map[string]any:
		if ref, ok := v["$ref"].(string); ok {
			if strings.HasPrefix(ref, "./") || strings.HasPrefix(ref, "../") {
				refPath := strings.TrimPrefix(ref, "./")
				fullPath := filepath.Join(baseDir, refPath)

				refData, err := os.ReadFile(fullPath)
				if err != nil {
					return fmt.Errorf("read ref %s: %w", ref, err)
				}

				var refSchema map[string]any
				if err := yaml.Unmarshal(refData, &refSchema); err != nil {
					return fmt.Errorf("unmarshal ref %s: %w", ref, err)
				}

				refBaseDir := filepath.Dir(fullPath)
				if err := dereferenceSchemaLegacy(refSchema, refBaseDir, rootDefs); err != nil {
					return err
				}

				delete(v, "$ref")
				for k, val := range refSchema {
					if k != "$schema" {
						v[k] = val
					}
				}
			}
		} else {
			for _, val := range v {
				if err := dereferenceSchemaLegacy(val, baseDir, rootDefs); err != nil {
					return err
				}
			}
		}
	case []any:
		for _, item := range v {
			if err := dereferenceSchemaLegacy(item, baseDir, rootDefs); err != nil {
				return err
			}
		}
	}
	return nil
}

func hoistNestedDefinitionsLegacy(data any, rootDefs map[string]any, path string) error {
	switch v := data.(type) {
	case map[string]any:
		if nestedDefs, ok := v["definitions"].(map[string]any); ok {
			_, isRoot := v["$schema"]
			if !isRoot {
				for name := range nestedDefs {
					if _, exists := rootDefs[name]; exists {
						return fmt.Errorf("duplicate definition %q found at %s - move shared definitions to ./common/ and reference them", name, path)
					}
				}
				for name, def := range nestedDefs {
					rootDefs[name] = def
				}
				delete(v, "definitions")
			}
		}
		for key, val := range v {
			childPath := path + "." + key
			if err := hoistNestedDefinitionsLegacy(val, rootDefs, childPath); err != nil {
				return err
			}
		}
	case []any:
		for i, item := range v {
			childPath := fmt.Sprintf("%s[%d]", path, i)
			if err := hoistNestedDefinitionsLegacy(item, rootDefs, childPath); err != nil {
				return err
			}
		}
	}
	return nil
}
