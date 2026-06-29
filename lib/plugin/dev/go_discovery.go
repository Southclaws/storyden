package dev

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

type PackageListOptions struct {
	Pattern     string
	IncludeDeps bool
	MaxPackages int
}

type PackageListResult struct {
	Packages  []PackageInfo `json:"packages"`
	Truncated bool          `json:"truncated"`
}

type PackageInfo struct {
	ImportPath    string `json:"import_path"`
	Name          string `json:"name"`
	ModulePath    string `json:"module_path,omitempty"`
	ModuleVersion string `json:"module_version,omitempty"`
	Standard      bool   `json:"standard"`
	Error         string `json:"error,omitempty"`
}

type PackageSymbolsOptions struct {
	ImportPath        string
	IncludeUnexported bool
	MaxSymbols        int
}

type PackageSymbolsResult struct {
	Package   PackageInfo       `json:"package"`
	Imports   []PackageInfo     `json:"imports"`
	Symbols   []GoSymbolSummary `json:"symbols"`
	Truncated bool              `json:"truncated"`
}

type SymbolDetailOptions struct {
	ImportPath string
	Symbol     string
}

type SymbolDetailResult struct {
	Package     PackageInfo       `json:"package"`
	Symbol      GoSymbolSummary   `json:"symbol"`
	Fields      []GoFieldInfo     `json:"fields,omitempty"`
	Methods     []GoMethodInfo    `json:"methods,omitempty"`
	Interface   []GoMethodInfo    `json:"interface_methods,omitempty"`
	RelatedInfo []GoSymbolSummary `json:"related_info,omitempty"`
}

type SymbolSearchOptions struct {
	Query             string
	Pattern           string
	IncludeDeps       bool
	IncludeUnexported bool
	MaxResults        int
}

type SymbolSearchResult struct {
	Results   []GoSymbolSummary `json:"results"`
	Truncated bool              `json:"truncated"`
}

type GoSymbolSummary struct {
	ImportPath string `json:"import_path"`
	Name       string `json:"name"`
	Kind       string `json:"kind"`
	Signature  string `json:"signature,omitempty"`
	Doc        string `json:"doc,omitempty"`
}

type GoFieldInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Tag  string `json:"tag,omitempty"`
}

type GoMethodInfo struct {
	Name      string `json:"name"`
	Receiver  string `json:"receiver,omitempty"`
	Signature string `json:"signature"`
}

func ListGoPackages(ctx context.Context, dir string, in PackageListOptions) (PackageListResult, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return PackageListResult{}, err
	}

	pattern := strings.TrimSpace(in.Pattern)
	if pattern == "" {
		pattern = "./..."
	}
	maxPackages := in.MaxPackages
	if maxPackages <= 0 || maxPackages > 200 {
		maxPackages = 100
	}

	args := []string{"list", "-json"}
	if in.IncludeDeps {
		args = append(args, "-deps")
	}
	args = append(args, pattern)

	output, err := runGo(ctx, root, args...)
	if err != nil {
		return PackageListResult{}, err
	}

	decoder := json.NewDecoder(strings.NewReader(output))
	out := PackageListResult{Packages: []PackageInfo{}}
	for {
		var raw goListPackage
		if err := decoder.Decode(&raw); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return PackageListResult{}, fmt.Errorf("decode go list output: %w", err)
		}
		if in.IncludeDeps && raw.Standard {
			continue
		}
		if len(out.Packages) >= maxPackages {
			out.Truncated = true
			break
		}
		out.Packages = append(out.Packages, packageInfoFromGoList(raw))
	}

	return out, nil
}

func GoPackageSymbols(ctx context.Context, dir string, in PackageSymbolsOptions) (PackageSymbolsResult, error) {
	importPath := strings.TrimSpace(in.ImportPath)
	if importPath == "" {
		return PackageSymbolsResult{}, errors.New("import_path is required")
	}
	maxSymbols := in.MaxSymbols
	if maxSymbols <= 0 || maxSymbols > 300 {
		maxSymbols = 100
	}

	pkg, err := LoadGoPackage(ctx, dir, importPath)
	if err != nil {
		return PackageSymbolsResult{}, err
	}
	docs := PackageDocs(pkg)

	names := pkg.Types.Scope().Names()
	sort.Strings(names)

	out := PackageSymbolsResult{
		Package: PackageInfoFromPackage(pkg),
		Imports: packageImports(pkg),
		Symbols: []GoSymbolSummary{},
	}
	for _, name := range names {
		if !in.IncludeUnexported && !ast.IsExported(name) {
			continue
		}
		if len(out.Symbols) >= maxSymbols {
			out.Truncated = true
			break
		}
		obj := pkg.Types.Scope().Lookup(name)
		out.Symbols = append(out.Symbols, SymbolSummary(pkg.PkgPath, obj, docs[name]))
	}

	return out, nil
}

func GoSymbolDetail(ctx context.Context, dir string, in SymbolDetailOptions) (SymbolDetailResult, error) {
	importPath := strings.TrimSpace(in.ImportPath)
	if importPath == "" {
		return SymbolDetailResult{}, errors.New("import_path is required")
	}
	name := strings.TrimSpace(in.Symbol)
	if name == "" {
		return SymbolDetailResult{}, errors.New("symbol is required")
	}

	pkg, err := LoadGoPackage(ctx, dir, importPath)
	if err != nil {
		return SymbolDetailResult{}, err
	}
	obj := pkg.Types.Scope().Lookup(name)
	if obj == nil {
		return SymbolDetailResult{}, fmt.Errorf("symbol %q not found in package %s", name, importPath)
	}

	docs := PackageDocs(pkg)
	out := SymbolDetailResult{
		Package: PackageInfoFromPackage(pkg),
		Symbol:  SymbolSummary(pkg.PkgPath, obj, docs[name]),
	}

	if typeName, ok := obj.(*types.TypeName); ok {
		named, _ := typeName.Type().(*types.Named)
		underlying := typeName.Type().Underlying()
		switch t := underlying.(type) {
		case *types.Struct:
			out.Fields = StructFields(t)
		case *types.Interface:
			out.Interface = InterfaceMethods(t)
		}
		if named != nil {
			out.Methods = NamedMethods(named)
		}
	}

	return out, nil
}

func GoSymbolSearch(ctx context.Context, dir string, in SymbolSearchOptions) (SymbolSearchResult, error) {
	query := strings.ToLower(strings.TrimSpace(in.Query))
	if query == "" {
		return SymbolSearchResult{}, errors.New("query is required")
	}
	if strings.ContainsAny(query, "*?[](){}^$|\\") {
		return SymbolSearchResult{}, fmt.Errorf("query must be a literal substring, not a regex or wildcard pattern; search for a plain term such as %q or %q", FirstPlainSearchTerm(query), LastPlainSearchTerm(query))
	}
	pattern := strings.TrimSpace(in.Pattern)
	if pattern == "" {
		pattern = "./..."
	}
	maxResults := in.MaxResults
	if maxResults <= 0 || maxResults > 300 {
		maxResults = 100
	}

	pkgs, err := LoadGoPackages(ctx, dir, pattern, in.IncludeDeps)
	if err != nil {
		return SymbolSearchResult{}, err
	}

	out := SymbolSearchResult{Results: []GoSymbolSummary{}}
	visited := map[string]bool{}
	for _, pkg := range pkgs {
		SearchPackageSymbols(pkg, query, in.IncludeUnexported, maxResults, &out, visited)
		if out.Truncated {
			break
		}
	}

	return out, nil
}

func LoadGoPackage(ctx context.Context, dir string, importPath string) (*packages.Package, error) {
	pkgs, err := LoadGoPackages(ctx, dir, importPath, false)
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		if pkg.PkgPath == importPath || pkg.ID == importPath {
			return pkg, nil
		}
	}
	if len(pkgs) == 1 {
		return pkgs[0], nil
	}
	return nil, fmt.Errorf("package %q not found", importPath)
}

func LoadGoPackages(ctx context.Context, dir string, pattern string, includeDeps bool) ([]*packages.Package, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	mode := packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedImports | packages.NeedModule
	cfg := &packages.Config{
		Context: ctx,
		Dir:     root,
		Mode:    mode,
		Tests:   false,
		Env:     append(os.Environ(), "GOFLAGS=-mod=mod"),
	}
	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		return nil, err
	}
	if includeDeps {
		pkgs = flattenPackageDeps(pkgs, 250)
	}
	if packageLoadErrorCount(pkgs) > 0 && len(pkgs) == 0 {
		return nil, fmt.Errorf("go package load failed")
	}
	return pkgs, nil
}

func PackageInfoFromPackage(pkg *packages.Package) PackageInfo {
	info := PackageInfo{
		ImportPath: pkg.PkgPath,
		Name:       pkg.Name,
		Standard:   pkg.Module == nil && isStandardImportPath(pkg.PkgPath),
	}
	if pkg.Module != nil {
		info.ModulePath = pkg.Module.Path
		info.ModuleVersion = pkg.Module.Version
	}
	return info
}

func PackageDocs(pkg *packages.Package) map[string]string {
	files := map[string]*ast.File{}
	for i, file := range pkg.Syntax {
		name := fmt.Sprintf("%s_%d.go", pkg.Name, i)
		if i < len(pkg.GoFiles) {
			name = pkg.GoFiles[i]
		}
		files[name] = file
	}
	docs := map[string]string{}
	if len(files) == 0 {
		return docs
	}
	docPkg := doc.New(&ast.Package{Name: pkg.Name, Files: files}, pkg.PkgPath, doc.AllDecls)
	for _, f := range docPkg.Funcs {
		docs[f.Name] = CleanDoc(f.Doc)
	}
	for _, t := range docPkg.Types {
		docs[t.Name] = CleanDoc(t.Doc)
	}
	for _, v := range docPkg.Vars {
		for _, name := range v.Names {
			docs[name] = CleanDoc(v.Doc)
		}
	}
	for _, c := range docPkg.Consts {
		for _, name := range c.Names {
			docs[name] = CleanDoc(c.Doc)
		}
	}
	return docs
}

func SymbolSummary(importPath string, obj types.Object, doc string) GoSymbolSummary {
	return GoSymbolSummary{
		ImportPath: importPath,
		Name:       obj.Name(),
		Kind:       symbolKind(obj),
		Signature:  Truncate(types.ObjectString(obj, qualifier), 1200),
		Doc:        Truncate(CleanDoc(doc), 600),
	}
}

func StructFields(s *types.Struct) []GoFieldInfo {
	fields := make([]GoFieldInfo, 0, s.NumFields())
	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)
		fields = append(fields, GoFieldInfo{
			Name: field.Name(),
			Type: types.TypeString(field.Type(), qualifier),
			Tag:  s.Tag(i),
		})
	}
	return fields
}

func InterfaceMethods(iface *types.Interface) []GoMethodInfo {
	iface = iface.Complete()
	methods := make([]GoMethodInfo, 0, iface.NumExplicitMethods())
	for i := 0; i < iface.NumExplicitMethods(); i++ {
		methods = append(methods, methodInfo(iface.ExplicitMethod(i)))
	}
	return methods
}

func NamedMethods(named *types.Named) []GoMethodInfo {
	seen := map[string]bool{}
	methods := []GoMethodInfo{}
	for _, typ := range []types.Type{named, types.NewPointer(named)} {
		set := types.NewMethodSet(typ)
		for i := 0; i < set.Len(); i++ {
			fn, ok := set.At(i).Obj().(*types.Func)
			if !ok || seen[fn.FullName()] {
				continue
			}
			seen[fn.FullName()] = true
			methods = append(methods, methodInfo(fn))
		}
	}
	sort.Slice(methods, func(i, j int) bool { return methods[i].Name < methods[j].Name })
	return methods
}

func SearchPackageSymbols(pkg *packages.Package, query string, includeUnexported bool, max int, out *SymbolSearchResult, visited map[string]bool) {
	if pkg == nil || pkg.Types == nil || visited[pkg.PkgPath] {
		return
	}
	visited[pkg.PkgPath] = true
	docs := PackageDocs(pkg)
	names := pkg.Types.Scope().Names()
	sort.Strings(names)
	for _, name := range names {
		if !includeUnexported && !ast.IsExported(name) {
			continue
		}
		obj := pkg.Types.Scope().Lookup(name)
		docText := docs[name]
		summary := SymbolSummary(pkg.PkgPath, obj, docText)
		haystack := strings.ToLower(summary.ImportPath + " " + summary.Name + " " + summary.Kind + " " + summary.Signature + " " + summary.Doc)
		if !strings.Contains(haystack, query) {
			continue
		}
		if len(out.Results) >= max {
			out.Truncated = true
			return
		}
		out.Results = append(out.Results, summary)
	}
}

func FirstPlainSearchTerm(query string) string {
	terms := PlainSearchTerms(query)
	if len(terms) == 0 {
		return "event"
	}
	return terms[0]
}

func LastPlainSearchTerm(query string) string {
	terms := PlainSearchTerms(query)
	if len(terms) == 0 {
		return "reply"
	}
	return terms[len(terms)-1]
}

func PlainSearchTerms(query string) []string {
	return strings.FieldsFunc(query, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9')
	})
}

func CleanDoc(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return strings.TrimSpace(s[:max]) + "..."
}

func runGo(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("go %s failed: %w: %s%s", strings.Join(args, " "), err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

type goListPackage struct {
	ImportPath string        `json:"ImportPath"`
	Name       string        `json:"Name"`
	Standard   bool          `json:"Standard"`
	Module     *goListModule `json:"Module"`
	Error      *goListError  `json:"Error"`
}

type goListModule struct {
	Path    string `json:"Path"`
	Version string `json:"Version"`
}

type goListError struct {
	Err string `json:"Err"`
}

func packageInfoFromGoList(raw goListPackage) PackageInfo {
	info := PackageInfo{
		ImportPath: raw.ImportPath,
		Name:       raw.Name,
		Standard:   raw.Standard,
	}
	if raw.Module != nil {
		info.ModulePath = raw.Module.Path
		info.ModuleVersion = raw.Module.Version
	}
	if raw.Error != nil {
		info.Error = raw.Error.Err
	}
	return info
}

func isStandardImportPath(importPath string) bool {
	if importPath == "" || strings.Contains(importPath, ".") {
		return false
	}
	first, _, _ := strings.Cut(importPath, "/")
	return first != "vendor"
}

func packageImports(pkg *packages.Package) []PackageInfo {
	imports := make([]PackageInfo, 0, len(pkg.Imports))
	for _, imported := range pkg.Imports {
		imports = append(imports, PackageInfoFromPackage(imported))
	}
	sort.Slice(imports, func(i, j int) bool { return imports[i].ImportPath < imports[j].ImportPath })
	return imports
}

func symbolKind(obj types.Object) string {
	switch obj.(type) {
	case *types.Func:
		return "func"
	case *types.TypeName:
		return "type"
	case *types.Const:
		return "const"
	case *types.Var:
		return "var"
	default:
		return "symbol"
	}
}

func qualifier(pkg *types.Package) string {
	if pkg == nil {
		return ""
	}
	return pkg.Name()
}

func methodInfo(fn *types.Func) GoMethodInfo {
	info := GoMethodInfo{
		Name:      fn.Name(),
		Signature: types.ObjectString(fn, qualifier),
	}
	if sig, ok := fn.Type().(*types.Signature); ok && sig.Recv() != nil {
		info.Receiver = types.TypeString(sig.Recv().Type(), qualifier)
	}
	return info
}

func flattenPackageDeps(roots []*packages.Package, max int) []*packages.Package {
	out := []*packages.Package{}
	seen := map[string]bool{}
	var visit func(*packages.Package)
	visit = func(pkg *packages.Package) {
		if pkg == nil || seen[pkg.PkgPath] || len(out) >= max {
			return
		}
		seen[pkg.PkgPath] = true
		out = append(out, pkg)
		imports := make([]*packages.Package, 0, len(pkg.Imports))
		for _, imp := range pkg.Imports {
			imports = append(imports, imp)
		}
		sort.Slice(imports, func(i, j int) bool { return imports[i].PkgPath < imports[j].PkgPath })
		for _, imp := range imports {
			visit(imp)
		}
	}
	for _, root := range roots {
		visit(root)
	}
	return out
}

func packageLoadErrorCount(pkgs []*packages.Package) int {
	count := 0
	for _, pkg := range pkgs {
		count += len(pkg.Errors)
	}
	return count
}
