package pluginbuilder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
)

const (
	maxDiscoveryFiles    = 1000
	maxDiscoveryFileSize = 2_000_000
)

type PackageListInput struct {
	Pattern     string `json:"pattern,omitempty" jsonschema:"Go package pattern. Use ./... for workspace packages or a full import path such as github.com/Southclaws/storyden/lib/plugin/rpc/..."`
	IncludeDeps bool   `json:"include_deps,omitempty" jsonschema:"Include transitive dependencies returned by go list -deps"`
	MaxPackages int    `json:"max_packages,omitempty" jsonschema:"Maximum packages to return"`
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

type PackageSymbolsInput struct {
	ImportPath        string `json:"import_path" jsonschema:"Full Go package import path to inspect"`
	IncludeUnexported bool   `json:"include_unexported,omitempty" jsonschema:"Include unexported package symbols"`
	MaxSymbols        int    `json:"max_symbols,omitempty" jsonschema:"Maximum symbols to return"`
}

type PackageSymbolsResult struct {
	Package   PackageInfo       `json:"package"`
	Imports   []PackageInfo     `json:"imports"`
	Symbols   []GoSymbolSummary `json:"symbols"`
	Truncated bool              `json:"truncated"`
}

type SymbolDetailInput struct {
	ImportPath string `json:"import_path" jsonschema:"Full Go package import path that owns the symbol"`
	Symbol     string `json:"symbol" jsonschema:"Exported or unexported symbol name to inspect"`
}

type SymbolDetailResult struct {
	Package     PackageInfo       `json:"package"`
	Symbol      GoSymbolSummary   `json:"symbol"`
	Fields      []GoFieldInfo     `json:"fields,omitempty"`
	Methods     []GoMethodInfo    `json:"methods,omitempty"`
	Interface   []GoMethodInfo    `json:"interface_methods,omitempty"`
	RelatedInfo []GoSymbolSummary `json:"related_info,omitempty"`
}

type SymbolSearchInput struct {
	Query             string `json:"query" jsonschema:"Case-insensitive literal substring to search for. This is not regex or wildcard search."`
	Pattern           string `json:"pattern,omitempty" jsonschema:"Package pattern to search. Defaults to ./... for workspace packages."`
	IncludeDeps       bool   `json:"include_deps,omitempty" jsonschema:"Also search imported dependency packages"`
	IncludeUnexported bool   `json:"include_unexported,omitempty" jsonschema:"Include unexported symbols"`
	MaxResults        int    `json:"max_results,omitempty" jsonschema:"Maximum symbol matches to return"`
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

func (a *Agent) addGoDiscoveryTools(add toolAdder) error {
	if err := add(functiontool.New(functiontool.Config{
		Name: "plugin_go_packages",
		Description: `List Go packages visible from the plugin workspace.

Use this as a directory listing for Go code. Start with pattern "./..." to see
workspace packages. Use a full package prefix with "/..." to discover subpackages,
for example "github.com/Southclaws/storyden/lib/plugin/..." or
"github.com/bwmarrin/discordgo/...".

Every result includes import_path. Pass that exact import_path to
plugin_go_package_symbols or plugin_go_symbol_detail to dive deeper.`,
	}, func(ctx adktool.Context, args PackageListInput) (PackageListResult, error) {
		return a.ListGoPackages(ctx, args)
	})); err != nil {
		return err
	}

	if err := add(functiontool.New(functiontool.Config{
		Name: "plugin_go_package_symbols",
		Description: `List symbols exported by one Go package.

Use this after plugin_go_packages or when you know the import path. This is the
main discovery tool for finding SDK capabilities without guessing method names.
For example, inspect "github.com/Southclaws/storyden/lib/plugin/rpc" to discover
available event payload types, or inspect a third-party package after adding it
to go.mod and running plugin_go_tidy.

The response includes symbol names, kinds, signatures, docs, imports, and the
owning import_path so you can recursively inspect related packages.`,
	}, func(ctx adktool.Context, args PackageSymbolsInput) (PackageSymbolsResult, error) {
		return a.GoPackageSymbols(ctx, args)
	})); err != nil {
		return err
	}

	if err := add(functiontool.New(functiontool.Config{
		Name: "plugin_go_symbol_detail",
		Description: `Inspect one symbol in one Go package.

Use this when package_symbols shows a promising type or function. For types, this
returns fields, methods, and interface methods where available. This is the
closest tool to editor "go to definition" plus "show methods" for the plugin
builder.`,
	}, func(ctx adktool.Context, args SymbolDetailInput) (SymbolDetailResult, error) {
		return a.GoSymbolDetail(ctx, args)
	})); err != nil {
		return err
	}

	return add(functiontool.New(functiontool.Config{
		Name: "plugin_go_symbol_search",
		Description: `Search symbols by name or documentation across Go packages.

Use this when you know the concept but not the API name, such as "reply",
"reaction", "discord", "message", or "event". Start with pattern "./..." for
workspace code. Use a full package prefix with "/..." to search a dependency
family. Set include_deps only when you need broader dependency search because it
can return a lot of packages.

The query is a literal case-insensitive substring, not a regex, glob, or
wildcard. Do not send queries such as "Event.*Reply"; search for "Event",
"Reply", or another plain word instead.`,
	}, func(ctx adktool.Context, args SymbolSearchInput) (SymbolSearchResult, error) {
		return a.GoSymbolSearch(ctx, args)
	}))
}

func (a *Agent) ListGoPackages(ctx context.Context, in PackageListInput) (PackageListResult, error) {
	workspace, err := a.Workspace(ctx)
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

	result, err := commandResult(workspace.Run(ctx, workspaceCommand("go", args...)))
	if err != nil {
		return PackageListResult{}, err
	}
	if !result.Success {
		return PackageListResult{}, fmt.Errorf("go list failed: %s%s", result.Output, result.Error)
	}

	decoder := json.NewDecoder(strings.NewReader(result.Output))
	out := PackageListResult{Packages: []PackageInfo{}}
	for {
		var raw goListPackage
		if err := decoder.Decode(&raw); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return PackageListResult{}, fmt.Errorf("decode go list output: %w", err)
		}
		if len(out.Packages) >= maxPackages {
			out.Truncated = true
			break
		}
		out.Packages = append(out.Packages, packageInfoFromGoList(raw))
	}

	return out, nil
}

func (a *Agent) GoPackageSymbols(ctx context.Context, in PackageSymbolsInput) (PackageSymbolsResult, error) {
	importPath := strings.TrimSpace(in.ImportPath)
	if importPath == "" {
		return PackageSymbolsResult{}, errors.New("import_path is required")
	}
	maxSymbols := in.MaxSymbols
	if maxSymbols <= 0 || maxSymbols > 300 {
		maxSymbols = 100
	}

	pkg, err := a.loadGoPackage(ctx, importPath)
	if err != nil {
		return PackageSymbolsResult{}, err
	}
	docs := packageDocs(pkg)

	names := pkg.Types.Scope().Names()
	sort.Strings(names)

	out := PackageSymbolsResult{
		Package: packageInfoFromPackage(pkg),
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
		out.Symbols = append(out.Symbols, symbolSummary(pkg.PkgPath, obj, docs[name]))
	}

	return out, nil
}

func (a *Agent) GoSymbolDetail(ctx context.Context, in SymbolDetailInput) (SymbolDetailResult, error) {
	importPath := strings.TrimSpace(in.ImportPath)
	if importPath == "" {
		return SymbolDetailResult{}, errors.New("import_path is required")
	}
	name := strings.TrimSpace(in.Symbol)
	if name == "" {
		return SymbolDetailResult{}, errors.New("symbol is required")
	}

	pkg, err := a.loadGoPackage(ctx, importPath)
	if err != nil {
		return SymbolDetailResult{}, err
	}
	obj := pkg.Types.Scope().Lookup(name)
	if obj == nil {
		return SymbolDetailResult{}, fmt.Errorf("symbol %q not found in package %s", name, importPath)
	}

	docs := packageDocs(pkg)
	out := SymbolDetailResult{
		Package: packageInfoFromPackage(pkg),
		Symbol:  symbolSummary(pkg.PkgPath, obj, docs[name]),
	}

	if typeName, ok := obj.(*types.TypeName); ok {
		named, _ := typeName.Type().(*types.Named)
		underlying := typeName.Type().Underlying()
		switch t := underlying.(type) {
		case *types.Struct:
			out.Fields = structFields(t)
		case *types.Interface:
			out.Interface = interfaceMethods(t)
		}
		if named != nil {
			out.Methods = namedMethods(named)
		}
	}

	return out, nil
}

func (a *Agent) GoSymbolSearch(ctx context.Context, in SymbolSearchInput) (SymbolSearchResult, error) {
	query := strings.ToLower(strings.TrimSpace(in.Query))
	if query == "" {
		return SymbolSearchResult{}, errors.New("query is required")
	}
	if strings.ContainsAny(query, "*?[](){}^$|\\") {
		return SymbolSearchResult{}, fmt.Errorf("query must be a literal substring, not a regex or wildcard pattern; search for a plain term such as %q or %q", firstPlainSearchTerm(query), lastPlainSearchTerm(query))
	}
	pattern := strings.TrimSpace(in.Pattern)
	if pattern == "" {
		pattern = "./..."
	}
	maxResults := in.MaxResults
	if maxResults <= 0 || maxResults > 300 {
		maxResults = 100
	}

	pkgs, err := a.loadGoPackages(ctx, pattern, in.IncludeDeps)
	if err != nil {
		return SymbolSearchResult{}, err
	}

	out := SymbolSearchResult{Results: []GoSymbolSummary{}}
	visited := map[string]bool{}
	for _, pkg := range pkgs {
		searchPackageSymbols(pkg, query, in.IncludeUnexported, maxResults, &out, visited)
		if out.Truncated {
			break
		}
	}

	return out, nil
}

func (a *Agent) loadGoPackage(ctx context.Context, importPath string) (*packages.Package, error) {
	pkgs, err := a.loadGoPackages(ctx, importPath, false)
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

func (a *Agent) loadGoPackages(ctx context.Context, pattern string, includeDeps bool) ([]*packages.Package, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return nil, err
	}

	var out []*packages.Package
	err = withWorkspaceMirror(ctx, workspace, func(dir string) error {
		mode := packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedImports | packages.NeedModule
		cfg := &packages.Config{
			Context: ctx,
			Dir:     dir,
			Mode:    mode,
			Tests:   false,
			Env:     append(os.Environ(), "GOFLAGS=-mod=mod"),
		}
		pkgs, err := packages.Load(cfg, pattern)
		if err != nil {
			return err
		}
		if includeDeps {
			pkgs = flattenPackageDeps(pkgs, 250)
		}
		if packageLoadErrorCount(pkgs) > 0 {
			// Keep partially-loaded packages useful for discovery, but fail if none loaded.
			if len(pkgs) == 0 {
				return fmt.Errorf("go package load failed")
			}
		}
		out = pkgs
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func withWorkspaceMirror(ctx context.Context, workspace workspaceprovider.Workspace, fn func(dir string) error) error {
	root, err := os.MkdirTemp("", "storyden-plugin-discovery-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(root)

	files, err := workspace.List(ctx, workspaceprovider.ListOptions{MaxFiles: maxDiscoveryFiles})
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.Size > maxDiscoveryFileSize {
			continue
		}
		if strings.HasSuffix(file.Path, ".zip") {
			continue
		}
		data, err := workspace.ReadFile(ctx, file.Path, maxDiscoveryFileSize)
		if err != nil {
			return err
		}
		if data.Truncated {
			continue
		}
		dst := filepath.Join(root, filepath.FromSlash(data.Path))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, data.Content, 0o644); err != nil {
			return err
		}
	}

	return fn(root)
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

func packageInfoFromPackage(pkg *packages.Package) PackageInfo {
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
		imports = append(imports, packageInfoFromPackage(imported))
	}
	sort.Slice(imports, func(i, j int) bool { return imports[i].ImportPath < imports[j].ImportPath })
	return imports
}

func packageDocs(pkg *packages.Package) map[string]string {
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
		docs[f.Name] = cleanDoc(f.Doc)
	}
	for _, t := range docPkg.Types {
		docs[t.Name] = cleanDoc(t.Doc)
	}
	for _, v := range docPkg.Vars {
		for _, name := range v.Names {
			docs[name] = cleanDoc(v.Doc)
		}
	}
	for _, c := range docPkg.Consts {
		for _, name := range c.Names {
			docs[name] = cleanDoc(c.Doc)
		}
	}
	return docs
}

func symbolSummary(importPath string, obj types.Object, doc string) GoSymbolSummary {
	return GoSymbolSummary{
		ImportPath: importPath,
		Name:       obj.Name(),
		Kind:       symbolKind(obj),
		Signature:  truncate(types.ObjectString(obj, qualifier), 1200),
		Doc:        truncate(cleanDoc(doc), 600),
	}
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

func structFields(s *types.Struct) []GoFieldInfo {
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

func interfaceMethods(iface *types.Interface) []GoMethodInfo {
	iface = iface.Complete()
	methods := make([]GoMethodInfo, 0, iface.NumExplicitMethods())
	for i := 0; i < iface.NumExplicitMethods(); i++ {
		methods = append(methods, methodInfo(iface.ExplicitMethod(i)))
	}
	return methods
}

func namedMethods(named *types.Named) []GoMethodInfo {
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

func searchPackageSymbols(pkg *packages.Package, query string, includeUnexported bool, max int, out *SymbolSearchResult, visited map[string]bool) {
	if pkg == nil || pkg.Types == nil || visited[pkg.PkgPath] {
		return
	}
	visited[pkg.PkgPath] = true
	docs := packageDocs(pkg)
	names := pkg.Types.Scope().Names()
	sort.Strings(names)
	for _, name := range names {
		if !includeUnexported && !ast.IsExported(name) {
			continue
		}
		obj := pkg.Types.Scope().Lookup(name)
		docText := docs[name]
		summary := symbolSummary(pkg.PkgPath, obj, docText)
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

func firstPlainSearchTerm(query string) string {
	terms := plainSearchTerms(query)
	if len(terms) == 0 {
		return "event"
	}
	return terms[0]
}

func lastPlainSearchTerm(query string) string {
	terms := plainSearchTerms(query)
	if len(terms) == 0 {
		return "reply"
	}
	return terms[len(terms)-1]
}

func plainSearchTerms(query string) []string {
	return strings.FieldsFunc(query, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9')
	})
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

func workspaceCommand(command string, args ...string) workspaceprovider.CommandSpec {
	return workspaceprovider.CommandSpec{Command: command, Args: args}
}

func packageLoadErrorCount(pkgs []*packages.Package) int {
	count := 0
	for _, pkg := range pkgs {
		count += len(pkg.Errors)
	}
	return count
}

func cleanDoc(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return strings.TrimSpace(s[:max]) + "..."
}
