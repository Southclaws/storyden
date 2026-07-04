package pluginbuilder

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
)

var incompleteImplementationWordPattern = regexp.MustCompile(`\b(todo|fixme|stub)\b`)

type ValidateInput struct {
	SkipGo bool `json:"skip_go,omitempty" jsonschema:"Skip Go formatting, dependency, vet, lint, and test checks"`
}

type ValidateResult struct {
	Success bool              `json:"success"`
	Checks  []ValidationCheck `json:"checks"`
	Message string            `json:"message,omitempty"`
}

type ValidationCheck struct {
	Name       string `json:"name"`
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
	Command    string `json:"command,omitempty"`
	Output     string `json:"output,omitempty"`
	Error      string `json:"error,omitempty"`
	Truncated  bool   `json:"truncated,omitempty"`
	DurationMS int64  `json:"duration_ms,omitempty"`
}

func (a *Agent) addValidateTools(add toolAdder) error {
	return add(functiontool.New(functiontool.Config{
		Name:        "plugin_validate",
		Description: "Run plugin source validation: manifest schema, manifest/code consistency, Go formatting, dependencies, vet/lint, and tests. Packaging and compilation happen inside plugin_install.",
	}, func(ctx adktool.Context, args ValidateInput) (ValidateResult, error) {
		return a.Validate(ctx, args)
	}))
}

func (a *Agent) Validate(ctx context.Context, in ValidateInput) (ValidateResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return ValidateResult{}, err
	}

	result := ValidateResult{Checks: []ValidationCheck{}}

	mf, err := readProjectManifest(ctx, workspace)
	result.addError("manifest", "manifest.yaml parses and matches the plugin manifest schema", err)

	var files []workspaceprovider.FileInfo
	if mf != nil {
		files, err = packageWorkspaceFiles(ctx, workspace)
		result.addError("workspace_files", "workspace files can be listed for packaging", err)

		if err == nil {
			err = validateHostAPIAccessManifest(ctx, workspace, mf.Manifest, files)
			result.addError("manifest_code_consistency", "manifest access matches Storyden host API client usage", err)

			err = validateNoIncompleteImplementationMarkers(ctx, workspace, files)
			result.addError("implementation_completeness", "plugin source has no placeholder, stub, dry-run, or TODO implementation markers", err)
		}
	}

	if !in.SkipGo {
		command, err := a.GoFormat(ctx)
		result.addCommand("go_fmt", command, err)
		command, err = a.GoTidy(ctx)
		result.addCommand("go_tidy", command, err)
		command, err = a.GoVet(ctx)
		result.addCommand("go_vet", command, err)
		command, err = a.GoTest(ctx, GoTestInput{})
		result.addCommand("go_test", command, err)
	}

	result.Success = true
	for _, check := range result.Checks {
		if !check.Success {
			result.Success = false
			break
		}
	}
	if result.Success {
		result.Message = "plugin validation passed"
	} else {
		result.Message = validationFailureSummary(result)
	}

	return result, nil
}

func (r *ValidateResult) addError(name string, successMessage string, err error) {
	check := ValidationCheck{
		Name:    name,
		Success: err == nil,
		Message: successMessage,
	}
	if err != nil {
		check.Message = err.Error()
		check.Error = err.Error()
	}
	r.Checks = append(r.Checks, check)
}

func (r *ValidateResult) addCommand(name string, command CommandResult, err error) {
	check := ValidationCheck{
		Name:       name,
		Success:    err == nil && command.Success,
		Command:    command.Command,
		Output:     command.Output,
		Error:      command.Error,
		Truncated:  command.Truncated,
		DurationMS: command.DurationMS,
	}
	if err != nil {
		check.Message = err.Error()
		check.Error = err.Error()
	} else if !command.Success {
		check.Message = strings.TrimSpace(command.Output)
		if check.Message == "" {
			check.Message = strings.TrimSpace(command.Error)
		}
		if check.Message == "" {
			check.Message = fmt.Sprintf("%s failed", name)
		}
	} else {
		check.Message = fmt.Sprintf("%s passed", name)
	}
	r.Checks = append(r.Checks, check)
}

func validationFailureSummary(result ValidateResult) string {
	failures := []string{}
	for _, check := range result.Checks {
		if check.Success {
			continue
		}
		message := strings.TrimSpace(check.Message)
		if message == "" {
			message = strings.TrimSpace(check.Error)
		}
		if message == "" {
			message = "failed"
		}
		failures = append(failures, fmt.Sprintf("%s: %s", check.Name, firstLine(message)))
	}
	if len(failures) == 0 {
		return "plugin validation failed"
	}
	return "plugin validation failed: " + strings.Join(failures, "; ")
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return strings.TrimSpace(s)
}

func validateNoIncompleteImplementationMarkers(ctx context.Context, workspace workspaceprovider.Workspace, files []workspaceprovider.FileInfo) error {
	findings := []string{}
	for _, file := range files {
		if !strings.HasSuffix(file.Path, ".go") {
			continue
		}
		data, err := workspace.ReadFile(ctx, file.Path, -1)
		if err != nil {
			return err
		}
		if isGeneratedGoSource(data.Content) {
			continue
		}
		for lineIndex, line := range strings.Split(string(data.Content), "\n") {
			if marker, ok := incompleteImplementationMarker(line); ok {
				findings = append(findings, fmt.Sprintf("%s:%d contains incomplete implementation marker %q", file.Path, lineIndex+1, marker))
				if len(findings) >= 10 {
					return fmt.Errorf("incomplete implementation markers found: %s", strings.Join(findings, "; "))
				}
			}
		}
	}
	if len(findings) > 0 {
		return fmt.Errorf("incomplete implementation markers found: %s", strings.Join(findings, "; "))
	}
	return nil
}

func isGeneratedGoSource(content []byte) bool {
	prefix := string(content)
	if len(prefix) > 2048 {
		prefix = prefix[:2048]
	}
	return strings.Contains(prefix, "Code generated") && strings.Contains(prefix, "DO NOT EDIT")
}

func incompleteImplementationMarker(line string) (string, bool) {
	lower := strings.ToLower(line)
	if marker := incompleteImplementationWordPattern.FindString(lower); marker != "" {
		return marker, true
	}
	markers := []string{
		"not implemented",
		"not yet implemented",
		"placeholder",
		"dry run",
		"dry-run",
		"would create",
		"would update",
		"would delete",
		"would send",
		"would post",
		"would call",
		"would execute",
		"implement actual",
		"implement later",
		"finish later",
		"fix later",
		"to be done later",
		"done later",
	}
	for _, marker := range markers {
		if strings.Contains(lower, marker) {
			return marker, true
		}
	}
	return "", false
}
