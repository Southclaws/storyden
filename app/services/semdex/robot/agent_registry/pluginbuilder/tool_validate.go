package pluginbuilder

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strconv"
	"strings"

	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

var incompleteImplementationWordPattern = regexp.MustCompile(`\b(todo|fixme|stub)\b`)

type ValidateInput struct {
	SkipGo bool `json:"skip_go,omitempty" jsonschema:"Skip Go formatting, dependency, vet, lint, and test checks"`
}

type ValidateResult struct {
	Success    bool              `json:"success"`
	Checks     []ValidationCheck `json:"checks"`
	Message    string            `json:"message,omitempty"`
	NextAction string            `json:"next_action,omitempty"`
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
		Description: "Check whether plugin source is ready to install. Runs manifest schema checks, manifest/code consistency checks, incomplete-implementation checks, gofmt, go mod tidy, go vet, plugin semantic lint, and go test. Use while iterating on source. Does not compile, package, upload, activate, or read runtime logs; plugin_install performs the single compile/package/install path.",
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

			err = validateConfigurationImplementation(ctx, workspace, mf.Manifest, files)
			result.addError("configuration_implementation", "manifest configuration fields are handled by plugin source", err)

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
		result.NextAction = "Source validation passed. Use plugin_install to compile once, package, upload or update, and activate when requested."
	} else {
		result.Message = validationFailureSummary(result)
		result.NextAction = validationNextAction(result)
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
		failures = append(failures, fmt.Sprintf("%s: %s", check.Name, validationCheckSummaryLine(check)))
	}
	if len(failures) == 0 {
		return "plugin validation failed"
	}
	return "plugin validation failed: " + strings.Join(failures, "; ")
}

func validationNextAction(result ValidateResult) string {
	if next := validationGoFailureNextAction(result); next != "" {
		return next
	}

	for _, check := range result.Checks {
		if check.Success {
			continue
		}

		switch check.Name {
		case "manifest":
			return "Fix manifest fields with plugin_manifest_write, then rerun plugin_validate."
		case "workspace_files":
			return "Inspect the workspace with plugin_file_list and plugin_file_read, then repair missing or unreadable files."
		case "manifest_code_consistency":
			return "Align manifest.yaml with code. If code uses BuildAPIClient, add access with the narrow required permissions; otherwise remove unnecessary access/code."
		case "configuration_implementation":
			return "Handle every manifest configuration_schema field in Go source. Parse required settings from the raw configuration map, keep the plugin running while settings are missing, and rerun plugin_validate."
		case "implementation_completeness":
			return "Replace placeholders, TODOs, dry-run logic, or stub behavior with real plugin behavior before installing."
		case "go_fmt":
			return "Run plugin_go_fmt or fix the formatting error, then rerun plugin_validate."
		case "go_tidy":
			return "Fix module/dependency issues, rerun plugin_go_tidy, then rerun plugin_validate."
		case "go_vet":
			return "Fix Go vet or Plugin Builder semantic lint failures, then rerun plugin_validate. If a method, field, type, or package is missing, use plugin_go_package_symbols, plugin_go_symbol_detail, or plugin_go_symbol_search to discover the actual API instead of asking the user."
		case "go_test":
			return "Fix failing tests or compilation errors, then rerun plugin_validate. If a method, field, type, or package is missing, use plugin_go_package_symbols, plugin_go_symbol_detail, or plugin_go_symbol_search to discover the actual API instead of asking the user."
		default:
			return fmt.Sprintf("Fix the failed %s check, then rerun plugin_validate.", check.Name)
		}
	}

	return "Review validation checks and rerun plugin_validate after making changes."
}

func validationGoFailureNextAction(result ValidateResult) string {
	for _, check := range result.Checks {
		if check.Success || (check.Name != "go_vet" && check.Name != "go_test") {
			continue
		}

		text := strings.ToLower(strings.Join([]string{check.Message, check.Output, check.Error}, "\n"))
		switch {
		case strings.Contains(text, "robotrunwithresponse") ||
			strings.Contains(text, "robotchatssewithresponse") ||
			strings.Contains(text, "robotchatsse"):
			return "Replace generated HTTP robot chat calls with pl.RunRobot(ctx, robotID, message), ensure manifest access.permissions includes USE_ROBOTS, then rerun plugin_validate."
		case strings.Contains(text, "undefined") ||
			strings.Contains(text, "unknown field") ||
			strings.Contains(text, "has no field or method"):
			return "Fix missing Go methods, fields, or types before semantic cleanup. Use plugin_storyden_sdk_search for Storyden APIs and plugin_go_package_symbols, plugin_go_symbol_detail, or plugin_go_symbol_search for external APIs; do not ask the user to choose around compile errors."
		}
	}
	return ""
}

func validationCheckSummaryLine(check ValidationCheck) string {
	candidates := []string{check.Message, check.Output, check.Error}
	for _, candidate := range candidates {
		for _, line := range strings.Split(candidate, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "# ") || line == "FAIL" {
				continue
			}
			return line
		}
	}

	return "failed"
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return strings.TrimSpace(s)
}

func validateConfigurationImplementation(ctx context.Context, workspace workspaceprovider.Workspace, manifest rpc.Manifest, files []workspaceprovider.FileInfo) error {
	fieldIDs := manifestConfigurationFieldIDs(manifest)
	if len(fieldIDs) == 0 {
		return nil
	}

	handledFields := map[string]struct{}{}
	emptyConfigStructs := []string{}
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
		emptyConfigStructs = append(emptyConfigStructs, goEmptyConfigurationStructs(file.Path, data.Content)...)
		for _, fieldID := range goConfigurationFieldReads(file.Path, data.Content) {
			handledFields[fieldID] = struct{}{}
		}
	}
	if len(emptyConfigStructs) > 0 {
		return fmt.Errorf("empty configuration structs are placeholders, not runtime configuration handling: %s", strings.Join(emptyConfigStructs, ", "))
	}

	missing := []string{}
	for _, fieldID := range fieldIDs {
		if _, ok := handledFields[fieldID]; !ok {
			missing = append(missing, fieldID)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("configuration_schema fields are not read from runtime configuration in Go source: %s", strings.Join(missing, ", "))
	}
	return nil
}

func goEmptyConfigurationStructs(path string, source []byte) []string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, source, parser.ParseComments)
	if err != nil {
		return nil
	}

	var names []string
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || !strings.Contains(strings.ToLower(typeSpec.Name.Name), "config") {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok || structType.Fields == nil || len(structType.Fields.List) > 0 {
				continue
			}
			names = append(names, typeSpec.Name.Name)
		}
	}
	return names
}

func goConfigurationFieldReads(path string, source []byte) []string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, source, parser.ParseComments)
	if err != nil {
		return nil
	}

	stringConstants := map[string]string{}
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || len(valueSpec.Values) == 0 {
				continue
			}
			for index, name := range valueSpec.Names {
				valueIndex := index
				if valueIndex >= len(valueSpec.Values) {
					valueIndex = len(valueSpec.Values) - 1
				}
				lit, ok := valueSpec.Values[valueIndex].(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					continue
				}
				value, err := strconv.Unquote(lit.Value)
				if err != nil {
					continue
				}
				stringConstants[name.Name] = value
			}
		}
	}

	configStructFields := goConfigStructJSONFields(file)
	configVars := map[string]string{}
	jsonMarshalVars := map[string]struct{}{}

	values := []string{}
	ast.Inspect(file, func(node ast.Node) bool {
		if valueSpec, ok := node.(*ast.ValueSpec); ok {
			if typeIdent, ok := valueSpec.Type.(*ast.Ident); ok {
				if _, ok := configStructFields[typeIdent.Name]; ok {
					for _, name := range valueSpec.Names {
						configVars[name.Name] = typeIdent.Name
					}
				}
			}
			return true
		}

		if assign, ok := node.(*ast.AssignStmt); ok {
			for i, rhs := range assign.Rhs {
				if call, ok := rhs.(*ast.CallExpr); ok && isJSONCall(call, "Marshal") && i < len(assign.Lhs) {
					if name, ok := assign.Lhs[i].(*ast.Ident); ok {
						jsonMarshalVars[name.Name] = struct{}{}
					}
					continue
				}
				if configTypeName, ok := configTypeFromExpression(rhs); ok && i < len(assign.Lhs) {
					if name, ok := assign.Lhs[i].(*ast.Ident); ok {
						configVars[name.Name] = configTypeName
					}
				}
			}
			return true
		}

		if call, ok := node.(*ast.CallExpr); ok && isJSONCall(call, "Unmarshal") && len(call.Args) >= 2 {
			if sourceIdent, ok := call.Args[0].(*ast.Ident); ok {
				if _, ok := jsonMarshalVars[sourceIdent.Name]; ok {
					if configTypeName, ok := configVarTypeFromUnmarshalTarget(call.Args[1], configVars); ok {
						values = append(values, configStructFields[configTypeName]...)
					}
				}
			}
			return true
		}

		if rangeStmt, ok := node.(*ast.RangeStmt); ok {
			key, ok := rangeStmt.Key.(*ast.Ident)
			if !ok || key.Name == "_" || rangeStmt.Body == nil {
				return true
			}
			ast.Inspect(rangeStmt.Body, func(child ast.Node) bool {
				switchStmt, ok := child.(*ast.SwitchStmt)
				if !ok {
					return true
				}
				tag, ok := switchStmt.Tag.(*ast.Ident)
				if !ok || tag.Name != key.Name {
					return true
				}
				for _, stmt := range switchStmt.Body.List {
					caseClause, ok := stmt.(*ast.CaseClause)
					if !ok {
						continue
					}
					for _, expr := range caseClause.List {
						switch typed := expr.(type) {
						case *ast.BasicLit:
							if typed.Kind != token.STRING {
								continue
							}
							value, err := strconv.Unquote(typed.Value)
							if err == nil {
								values = append(values, value)
							}
						case *ast.Ident:
							if value, ok := stringConstants[typed.Name]; ok {
								values = append(values, value)
							}
						}
					}
				}
				return true
			})
			return true
		}

		indexExpr, ok := node.(*ast.IndexExpr)
		if !ok {
			return true
		}

		switch index := indexExpr.Index.(type) {
		case *ast.BasicLit:
			if index.Kind != token.STRING {
				return true
			}
			value, err := strconv.Unquote(index.Value)
			if err != nil {
				return true
			}
			values = append(values, value)
		case *ast.Ident:
			if value, ok := stringConstants[index.Name]; ok {
				values = append(values, value)
			}
		}
		return true
	})
	return values
}

func goConfigStructJSONFields(file *ast.File) map[string][]string {
	out := map[string][]string{}
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || !strings.Contains(strings.ToLower(typeSpec.Name.Name), "config") {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok || structType.Fields == nil {
				continue
			}
			for _, field := range structType.Fields.List {
				if field.Tag == nil {
					continue
				}
				tag, err := strconv.Unquote(field.Tag.Value)
				if err != nil {
					continue
				}
				jsonName := strings.Split(strings.TrimPrefix(tag, "json:\""), ",")[0]
				jsonName = strings.TrimSuffix(jsonName, "\"")
				if jsonName == "" || jsonName == "-" {
					continue
				}
				out[typeSpec.Name.Name] = append(out[typeSpec.Name.Name], jsonName)
			}
		}
	}
	return out
}

func configTypeFromExpression(expr ast.Expr) (string, bool) {
	switch typed := expr.(type) {
	case *ast.CompositeLit:
		if ident, ok := typed.Type.(*ast.Ident); ok {
			return ident.Name, true
		}
	case *ast.UnaryExpr:
		if typed.Op == token.AND {
			return configTypeFromExpression(typed.X)
		}
	}
	return "", false
}

func configVarTypeFromUnmarshalTarget(expr ast.Expr, configVars map[string]string) (string, bool) {
	switch typed := expr.(type) {
	case *ast.Ident:
		value, ok := configVars[typed.Name]
		return value, ok
	case *ast.UnaryExpr:
		if typed.Op == token.AND {
			return configVarTypeFromUnmarshalTarget(typed.X, configVars)
		}
	}
	return "", false
}

func isJSONCall(call *ast.CallExpr, name string) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != name {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	return ok && pkg.Name == "json"
}

func manifestConfigurationFieldIDs(manifest rpc.Manifest) []string {
	schema, ok := manifest.ConfigurationSchema.Get()
	if !ok {
		return nil
	}

	ids := []string{}
	for _, field := range schema.Fields {
		switch typed := field.PluginConfigurationFieldUnion.(type) {
		case *rpc.PluginConfigurationFieldString:
			ids = append(ids, typed.ID)
		case *rpc.PluginConfigurationFieldNumber:
			ids = append(ids, typed.ID)
		case *rpc.PluginConfigurationFieldBoolean:
			ids = append(ids, typed.ID)
		}
	}
	return ids
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
		"for now",
		"add manifest configuration_schema fields here",
		"canned summary",
		"fake summary",
		"placeholder summary",
		"requested from robot system",
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
