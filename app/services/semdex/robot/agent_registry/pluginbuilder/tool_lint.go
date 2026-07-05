package pluginbuilder

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
)

type PluginLintResult struct {
	Success bool              `json:"success"`
	Issues  []PluginLintIssue `json:"issues,omitempty"`
}

func (r PluginLintResult) Format() string {
	if r.Success {
		return "plugin lint passed"
	}

	lines := []string{"plugin lint failed:"}
	for _, issue := range r.Issues {
		lines = append(lines, fmt.Sprintf("%s:%d: %s", issue.Path, issue.Line, issue.Message))
	}
	return strings.Join(lines, "\n")
}

type PluginLintIssue struct {
	Path    string `json:"path"`
	Line    int    `json:"line"`
	Message string `json:"message"`
}

func (a *Agent) PluginLint(ctx context.Context) (PluginLintResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return PluginLintResult{}, err
	}

	files, err := workspace.List(ctx, workspaceprovider.ListOptions{MaxFiles: 1000})
	if err != nil {
		return PluginLintResult{}, err
	}

	var issues []PluginLintIssue
	for _, file := range files {
		if !shouldLintGoFile(file.Path) {
			continue
		}

		source, err := workspace.ReadFile(ctx, file.Path, -1)
		if err != nil {
			return PluginLintResult{}, err
		}

		fileIssues, err := lintGoSource(file.Path, source.Content)
		if err != nil {
			return PluginLintResult{}, err
		}
		issues = append(issues, fileIssues...)
	}

	return PluginLintResult{
		Success: len(issues) == 0,
		Issues:  issues,
	}, nil
}

func shouldLintGoFile(path string) bool {
	if !strings.HasSuffix(path, ".go") {
		return false
	}
	if strings.HasSuffix(path, "_test.go") {
		return false
	}
	for _, part := range strings.Split(path, "/") {
		if part == "vendor" || part == ".git" {
			return false
		}
	}
	return true
}

func lintGoSource(path string, source []byte) ([]PluginLintIssue, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, source, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var issues []PluginLintIssue
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		if issue, ok := lintUnsupportedStorydenSDKCall(path, fset, call); ok {
			issues = append(issues, issue)
			return true
		}
		if !isEventHandlerRegistration(call) {
			return true
		}

		for _, arg := range call.Args {
			fn, ok := arg.(*ast.FuncLit)
			if !ok || !funcReturnsError(fn) {
				continue
			}
			issues = append(issues, lintEventHandlerBody(path, fset, fn.Body)...)
		}

		return true
	})

	return issues, nil
}

func lintUnsupportedStorydenSDKCall(path string, fset *token.FileSet, call *ast.CallExpr) (PluginLintIssue, bool) {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return PluginLintIssue{}, false
	}

	switch selector.Sel.Name {
	case "HandleEventRPC":
		return PluginLintIssue{
			Path:    path,
			Line:    fset.Position(selector.Sel.Pos()).Line,
			Message: "Storyden plugin SDK has no HandleEventRPC method; use typed event helpers such as pl.OnThreadReplyCreated(...) or pl.OnActivityCreated(...), or use pl.On(eventName, handler) only when a typed helper is unavailable",
		}, true
	case "RobotRunWithResponse":
		return PluginLintIssue{
			Path:    path,
			Line:    fset.Position(selector.Sel.Pos()).Line,
			Message: "Storyden generated HTTP API has no RobotRunWithResponse plugin API; use pl.RunRobot(ctx, robotID, message) and add USE_ROBOTS access instead",
		}, true
	case "RobotChatSSE", "RobotChatSSEWithResponse":
		return PluginLintIssue{
			Path:    path,
			Line:    fset.Position(selector.Sel.Pos()).Line,
			Message: "Storyden RobotChatSSE is a UI streaming endpoint, not plugin-to-host robot execution; use pl.RunRobot(ctx, robotID, message) and use the returned summary",
		}, true
	default:
		return PluginLintIssue{}, false
	}
}

func isEventHandlerRegistration(call *ast.CallExpr) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	return ok && strings.HasPrefix(selector.Sel.Name, "On")
}

func funcReturnsError(fn *ast.FuncLit) bool {
	if fn.Type.Results == nil {
		return false
	}
	for _, field := range fn.Type.Results.List {
		ident, ok := field.Type.(*ast.Ident)
		if ok && ident.Name == "error" {
			return true
		}
	}
	return false
}

func lintEventHandlerBody(path string, fset *token.FileSet, body *ast.BlockStmt) []PluginLintIssue {
	var issues []PluginLintIssue
	ast.Inspect(body, func(node ast.Node) bool {
		ifs, ok := node.(*ast.IfStmt)
		if !ok || !isFailureCondition(fset, ifs.Cond) {
			return true
		}

		switch {
		case blockReturnsNil(ifs.Body):
			issues = append(issues, PluginLintIssue{
				Path:    path,
				Line:    fset.Position(ifs.Pos()).Line,
				Message: "event handler swallows a required action failure; return an error instead of nil",
			})
		case isResponseFailureCondition(fset, ifs.Cond) && blockReturnsBareErr(ifs.Body):
			issues = append(issues, PluginLintIssue{
				Path:    path,
				Line:    fset.Position(ifs.Pos()).Line,
				Message: "event handler returns err for an HTTP response failure after err is already nil; return a descriptive error instead",
			})
		}
		return true
	})
	return issues
}

func isFailureCondition(fset *token.FileSet, expr ast.Expr) bool {
	condition := strings.ToLower(exprString(fset, expr))
	compact := strings.NewReplacer(" ", "", "\t", "", "\n", "").Replace(condition)

	return strings.Contains(compact, "err!=nil") ||
		isResponseFailureCondition(fset, expr)
}

func isResponseFailureCondition(fset *token.FileSet, expr ast.Expr) bool {
	condition := strings.ToLower(exprString(fset, expr))
	compact := strings.NewReplacer(" ", "", "\t", "", "\n", "").Replace(condition)

	return strings.Contains(condition, "statuscode") ||
		strings.Contains(condition, "json200") ||
		strings.Contains(condition, "jsondefault") ||
		strings.Contains(compact, ">=400") ||
		strings.Contains(compact, ">299") ||
		strings.Contains(compact, "<200")
}

func blockReturnsNil(block *ast.BlockStmt) bool {
	return blockReturnsIdent(block, "nil")
}

func blockReturnsBareErr(block *ast.BlockStmt) bool {
	return blockReturnsIdent(block, "err")
}

func blockReturnsIdent(block *ast.BlockStmt, name string) bool {
	found := false
	ast.Inspect(block, func(node ast.Node) bool {
		ret, ok := node.(*ast.ReturnStmt)
		if !ok {
			return true
		}
		if len(ret.Results) != 1 {
			return true
		}
		ident, ok := ret.Results[0].(*ast.Ident)
		if ok && ident.Name == name {
			found = true
			return false
		}
		return true
	})
	return found
}

func exprString(fset *token.FileSet, expr ast.Expr) string {
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, expr); err != nil {
		return ""
	}
	return buf.String()
}
