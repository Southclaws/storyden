package pluginbuilder

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"

	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool/functiontool"
)

type ASTInput struct {
	Path string `json:"path" jsonschema:"Workspace-relative Go source file path"`
}

type ASTResult struct {
	Package   string       `json:"package"`
	Imports   []string     `json:"imports"`
	Functions []SymbolInfo `json:"functions"`
	Types     []SymbolInfo `json:"types"`
}

type SymbolInfo struct {
	Name string `json:"name"`
	Line int    `json:"line"`
}

func (a *Agent) addASTTools(add toolAdder) error {
	return add(functiontool.New(functiontool.Config{
		Name:        "plugin_go_ast",
		Description: "Parse a Go source file and return its package name, imports, functions, and type declarations.",
	}, func(ctx adkagent.Context, args ASTInput) (ASTResult, error) {
		result, err := a.AST(ctx, args)
		if err != nil {
			return ASTResult{}, err
		}
		return result, nil
	}))
}

func (a *Agent) AST(ctx context.Context, in ASTInput) (ASTResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return ASTResult{}, err
	}

	source, err := workspace.ReadFile(ctx, in.Path, -1)
	if err != nil {
		return ASTResult{}, err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, in.Path, source.Content, parser.ParseComments)
	if err != nil {
		return ASTResult{}, err
	}

	out := ASTResult{
		Package: file.Name.Name,
		Imports: []string{},
	}

	for _, imp := range file.Imports {
		unquoted, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			unquoted = strings.Trim(imp.Path.Value, `"`)
		}
		out.Imports = append(out.Imports, unquoted)
	}

	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			out.Functions = append(out.Functions, SymbolInfo{
				Name: d.Name.Name,
				Line: fset.Position(d.Pos()).Line,
			})
		case *ast.GenDecl:
			if d.Tok != token.TYPE {
				continue
			}
			for _, spec := range d.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				out.Types = append(out.Types, SymbolInfo{
					Name: ts.Name.Name,
					Line: fset.Position(ts.Pos()).Line,
				})
			}
		}
	}

	return out, nil
}
