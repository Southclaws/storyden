package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func run() error {
	var inputFile, outputFile string
	var baseDir string

	if len(os.Args) >= 3 {
		inputFile = os.Args[1]
		outputFile = os.Args[2]
		baseDir = filepath.Dir(inputFile)
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}

		baseDir = "api"
		if filepath.Base(cwd) == "api" {
			baseDir = "."
		}

		inputFile = filepath.Join(baseDir, "robots.yaml")
		outputFile = filepath.Join(baseDir, "robots.json")
	}

	fmt.Printf("Dereferencing %s -> %s\n", inputFile, outputFile)

	yamlData, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("read input file: %w", err)
	}

	var data map[string]any
	if err := yaml.Unmarshal(yamlData, &data); err != nil {
		return fmt.Errorf("unmarshal YAML: %w", err)
	}

	rootDefs := make(map[string]any)
	if existing, ok := data["definitions"].(map[string]any); ok {
		for k, v := range existing {
			rootDefs[k] = v
		}
	}

	if err := dereferenceSchema(data, baseDir, rootDefs); err != nil {
		return fmt.Errorf("dereference schema: %w", err)
	}

	data["definitions"] = rootDefs

	if err := hoistNestedDefinitions(data, rootDefs, "root"); err != nil {
		return err
	}

	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	fmt.Printf("✓ Successfully dereferenced schema to %s\n", outputFile)
	return nil
}

func dereferenceSchema(data any, baseDir string, rootDefs map[string]any) error {
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
				if err := dereferenceSchema(refSchema, refBaseDir, rootDefs); err != nil {
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
				if err := dereferenceSchema(val, baseDir, rootDefs); err != nil {
					return err
				}
			}
		}
	case []any:
		for _, item := range v {
			if err := dereferenceSchema(item, baseDir, rootDefs); err != nil {
				return err
			}
		}
	}
	return nil
}

func hoistNestedDefinitions(data any, rootDefs map[string]any, path string) error {
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
			if err := hoistNestedDefinitions(val, rootDefs, childPath); err != nil {
				return err
			}
		}
	case []any:
		for i, item := range v {
			childPath := fmt.Sprintf("%s[%d]", path, i)
			if err := hoistNestedDefinitions(item, rootDefs, childPath); err != nil {
				return err
			}
		}
	}
	return nil
}
