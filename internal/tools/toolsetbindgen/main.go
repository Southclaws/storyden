package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

const toolsetsImportPath = "github.com/Southclaws/storyden/app/services/semdex/robot/toolsets"

var systemToolsetIDPattern = regexp.MustCompile(`^system\.[a-z][a-z0-9_]*(?:\.[a-z][a-z0-9_]*)*$`)

type toolsetBinding struct {
	ID          string
	PackageName string
	ToolNames   []string
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: toolsetbindgen <input.json> <toolsets-dir>")
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("read %s: %w", args[0], err)
	}
	bindings, err := parseToolsetBindings(data)
	if err != nil {
		return err
	}
	if err := writeBindings(args[1], bindings); err != nil {
		return err
	}

	fmt.Printf("toolsetbindgen: generated %d system Toolsets → %s\n", len(bindings), args[1])
	return nil
}

func parseToolsetBindings(data []byte) ([]toolsetBinding, error) {
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("unmarshal Robot schema: %w", err)
	}
	definitions, ok := schema["definitions"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("Robot schema has no definitions")
	}

	toolNamesByToolset := map[string]map[string]struct{}{}
	definitionNames := make([]string, 0, len(definitions))
	for name := range definitions {
		definitionNames = append(definitionNames, name)
	}
	sort.Strings(definitionNames)

	for _, definitionName := range definitionNames {
		if !isToolDefinitionName(definitionName) {
			continue
		}
		definition, ok := definitions[definitionName].(map[string]any)
		if !ok {
			continue
		}
		if _, isTool := definition["x-storyden-tool"]; !isTool {
			continue
		}
		toolName, ok := definition["title"].(string)
		if !ok || strings.TrimSpace(toolName) == "" {
			return nil, fmt.Errorf("%s has no tool title", definitionName)
		}
		extensionRaw, exists := definition["x-storyden"]
		if !exists {
			return nil, fmt.Errorf("%s must declare x-storyden.toolsets", definitionName)
		}
		extension, ok := extensionRaw.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("%s x-storyden must be an object", definitionName)
		}
		toolsetsRaw, exists := extension["toolsets"]
		if !exists {
			return nil, fmt.Errorf("%s must declare x-storyden.toolsets", definitionName)
		}
		toolsets, ok := toolsetsRaw.([]any)
		if !ok {
			return nil, fmt.Errorf("%s x-storyden.toolsets must be a string array", definitionName)
		}
		for _, raw := range toolsets {
			id, ok := raw.(string)
			if !ok {
				return nil, fmt.Errorf("%s x-storyden.toolsets must contain only strings", definitionName)
			}
			if !systemToolsetIDPattern.MatchString(id) {
				return nil, fmt.Errorf("%s has invalid system Toolset ID %q", definitionName, id)
			}
			if toolNamesByToolset[id] == nil {
				toolNamesByToolset[id] = map[string]struct{}{}
			}
			toolNamesByToolset[id][toolName] = struct{}{}
		}
	}

	ids := make([]string, 0, len(toolNamesByToolset))
	for id := range toolNamesByToolset {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	bindings := make([]toolsetBinding, 0, len(ids))
	for _, id := range ids {
		toolNames := make([]string, 0, len(toolNamesByToolset[id]))
		for toolName := range toolNamesByToolset[id] {
			toolNames = append(toolNames, toolName)
		}
		sort.Strings(toolNames)
		bindings = append(bindings, toolsetBinding{
			ID:          id,
			PackageName: strings.ReplaceAll(id, ".", "_"),
			ToolNames:   toolNames,
		})
	}
	return bindings, nil
}

func isToolDefinitionName(name string) bool {
	return strings.HasPrefix(name, "Tool") &&
		len(name) > len("Tool") &&
		!strings.HasSuffix(name, "Input") &&
		!strings.HasSuffix(name, "Output")
}

func writeBindings(outputDir string, bindings []toolsetBinding) error {
	generatedPackages := make(map[string]struct{}, len(bindings))
	for _, binding := range bindings {
		packageDir := filepath.Join(outputDir, binding.PackageName)
		definitionPath := filepath.Join(packageDir, "definition.go")
		if _, err := os.Stat(definitionPath); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Toolset %q needs a hand-written definition at %s", binding.ID, definitionPath)
			}
			return err
		}
		generated, err := renderToolNames(binding)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(packageDir, "tools_gen.go"), generated, 0o644); err != nil {
			return fmt.Errorf("write generated Toolset tools for %s: %w", binding.ID, err)
		}
		generatedPackages[binding.PackageName] = struct{}{}
	}

	if err := removeStaleBindings(outputDir, generatedPackages); err != nil {
		return err
	}
	generated, err := renderSystemBindings(bindings)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outputDir, "system_bindings_gen.go"), generated, 0o644); err != nil {
		return fmt.Errorf("write generated system Toolset bindings: %w", err)
	}
	return nil
}

func removeStaleBindings(outputDir string, generatedPackages map[string]struct{}) error {
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "system_") {
			continue
		}
		if _, ok := generatedPackages[entry.Name()]; ok {
			continue
		}
		path := filepath.Join(outputDir, entry.Name(), "tools_gen.go")
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove stale generated binding %s: %w", path, err)
		}
	}
	return nil
}

func renderToolNames(binding toolsetBinding) ([]byte, error) {
	var output bytes.Buffer
	if err := toolNamesTemplate.Execute(&output, binding); err != nil {
		return nil, err
	}
	formatted, err := format.Source(output.Bytes())
	if err != nil {
		return nil, fmt.Errorf("format generated tools for %s: %w\n%s", binding.ID, err, output.String())
	}
	return formatted, nil
}

func renderSystemBindings(bindings []toolsetBinding) ([]byte, error) {
	var output bytes.Buffer
	if err := systemBindingsTemplate.Execute(&output, bindings); err != nil {
		return nil, err
	}
	formatted, err := format.Source(output.Bytes())
	if err != nil {
		return nil, fmt.Errorf("format generated system bindings: %w\n%s", err, output.String())
	}
	return formatted, nil
}

var toolNamesTemplate = template.Must(template.New("tool-names").Parse(`// Code generated by toolsetbindgen. DO NOT EDIT.

package {{.PackageName}}

var ToolNames = []string{
{{- range .ToolNames}}
	{{printf "%q" .}},
{{- end}}
}
`))

var systemBindingsTemplate = template.Must(template.New("system-bindings").Parse(`// Code generated by toolsetbindgen. DO NOT EDIT.

package toolsets

import (
{{- range .}}
	{{.PackageName}} "` + toolsetsImportPath + `/{{.PackageName}}"
{{- end}}
)

func systemDefinitions() []Definition {
	return []Definition{
{{- range .}}
		{
			ID: {{.PackageName}}.ID,
			Name: {{.PackageName}}.Name,
			Description: {{.PackageName}}.Description,
			Instruction: {{.PackageName}}.Instruction,
			ToolNames: append([]string(nil), {{.PackageName}}.ToolNames...),
			Source: SourceSystem,
		},
{{- end}}
	}
}
`))
