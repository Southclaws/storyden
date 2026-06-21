package pluginbuilder

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

const (
	storydenSDKPackage          = "github.com/Southclaws/storyden/sdk/go/storyden"
	storydenSDKRPCPackage       = "github.com/Southclaws/storyden/lib/plugin/rpc"
	storydenSDKOpenAPIPackage   = "github.com/Southclaws/storyden/app/transports/http/openapi"
	storydenSDKOperationPackage = "github.com/Southclaws/storyden/app/transports/http/openapi/operation"
)

type StorydenSDKSearchInput struct {
	Query      string `json:"query" jsonschema:"Case-insensitive literal substring to search for in Storyden plugin SDK and host API symbols"`
	Area       string `json:"area,omitempty" jsonschema:"Optional area: rpc, events, manifest, http_api, operations, or all"`
	MaxResults int    `json:"max_results,omitempty" jsonschema:"Maximum matching SDK symbols to return"`
}

type StorydenSDKSearchResult struct {
	Packages  []PackageInfo     `json:"packages"`
	Symbols   []GoSymbolSummary `json:"symbols"`
	Hints     []StorydenSDKHint `json:"hints,omitempty"`
	Truncated bool              `json:"truncated"`
}

type StorydenSDKHint struct {
	Message string `json:"message"`
}

type StorydenSDKEventsInput struct {
	Query     string `json:"query,omitempty" jsonschema:"Optional case-insensitive literal substring such as reply, thread, post, react, account, node, or report"`
	MaxEvents int    `json:"max_events,omitempty" jsonschema:"Maximum matching event definitions to return"`
}

type StorydenSDKEventsResult struct {
	ImportPath string                 `json:"import_path"`
	Events     []StorydenSDKEventInfo `json:"events"`
	Hints      []StorydenSDKHint      `json:"hints,omitempty"`
	Truncated  bool                   `json:"truncated"`
}

type StorydenSDKEventInfo struct {
	Event         string          `json:"event"`
	ManifestConst string          `json:"manifest_const"`
	PayloadType   GoSymbolSummary `json:"payload_type"`
	Fields        []GoFieldInfo   `json:"fields,omitempty"`
	FieldUsages   []SDKFieldUsage `json:"field_usages,omitempty"`
}

type SDKFieldUsage struct {
	Field      string `json:"field"`
	Type       string `json:"type"`
	Expression string `json:"expression"`
	Use        string `json:"use"`
}

func (a *Agent) addStorydenSDKTools(add toolAdder) error {
	if err := add(functiontool.New(functiontool.Config{
		Name: "plugin_storyden_sdk_events",
		Description: `Discover Storyden plugin event names and payload structs.

Use this before editing manifest.yaml events_consumed or writing an event
handler. This is a targeted Storyden SDK tool powered by Go package analysis,
not hard-coded docs. It searches only:
github.com/Southclaws/storyden/lib/plugin/rpc

Examples:
- query "reply" to find EventThreadReplyCreated, EventThreadReplyPublished,
  and related payload fields.
- query "react" to find post reaction events.
- omit query to list the first event definitions.

Use the event field value in manifest.yaml. Use the payload_type name in Go
handler signatures.`,
	}, func(ctx adktool.Context, args StorydenSDKEventsInput) (StorydenSDKEventsResult, error) {
		return a.StorydenSDKEvents(ctx, args)
	})); err != nil {
		return err
	}

	return add(functiontool.New(functiontool.Config{
		Name: "plugin_storyden_sdk_search",
		Description: `Search Storyden plugin SDK and host API Go symbols.

Use this before the generic Go search when working with Storyden-specific
imports. It searches only known Storyden SDK/API packages and omits unrelated
stdlib imports from the response.

Areas:
- plugin: plugin runtime symbols in sdk/go/storyden such as Plugin and BuildAPIClient
- events: plugin event names and payload structs in lib/plugin/rpc
- manifest: manifest and capability symbols in lib/plugin/rpc
- rpc: all lib/plugin/rpc symbols
- http_api: generated HTTP API client, request, and response symbols
- operations: generated operation ID symbols
- all: all of the above

The query is a literal case-insensitive substring, not regex or wildcard search.`,
	}, func(ctx adktool.Context, args StorydenSDKSearchInput) (StorydenSDKSearchResult, error) {
		return a.StorydenSDKSearch(ctx, args)
	}))
}

func (a *Agent) StorydenSDKEvents(ctx context.Context, in StorydenSDKEventsInput) (StorydenSDKEventsResult, error) {
	query := strings.ToLower(strings.TrimSpace(in.Query))
	if strings.ContainsAny(query, "*?[](){}^$|\\") {
		return StorydenSDKEventsResult{}, fmt.Errorf("query must be a literal substring, not a regex or wildcard pattern; search for a plain term such as %q", firstPlainSearchTerm(query))
	}

	maxEvents := in.MaxEvents
	if maxEvents <= 0 || maxEvents > 200 {
		maxEvents = 60
	}

	pkg, err := a.loadGoPackage(ctx, storydenSDKRPCPackage)
	if err != nil {
		return StorydenSDKEventsResult{}, err
	}

	docs := packageDocs(pkg)
	names := pkg.Types.Scope().Names()
	sort.Strings(names)

	payloads := map[string]types.Object{}
	for _, name := range names {
		obj := pkg.Types.Scope().Lookup(name)
		if _, ok := obj.(*types.TypeName); ok && strings.HasPrefix(name, "Event") {
			payloads[name] = obj
		}
	}

	out := StorydenSDKEventsResult{
		ImportPath: storydenSDKRPCPackage,
		Events:     []StorydenSDKEventInfo{},
		Hints: []StorydenSDKHint{
			{Message: "Use event values such as EventThreadReplyCreated in manifest.yaml events_consumed."},
			{Message: "Use payload type names such as *rpc.EventThreadReplyCreated in event handler signatures."},
			{Message: "When passing Storyden resource IDs to generated openapi parameters, use event.Field.String(); never use string(event.Field[:])."},
		},
	}

	for _, name := range names {
		if !strings.HasPrefix(name, "EventEvent") {
			continue
		}
		obj, ok := pkg.Types.Scope().Lookup(name).(*types.Const)
		if !ok {
			continue
		}
		eventValue := strings.Trim(constant.StringVal(obj.Val()), `"`)
		if query != "" && !strings.Contains(strings.ToLower(name+" "+eventValue), query) {
			continue
		}
		if len(out.Events) >= maxEvents {
			out.Truncated = true
			break
		}

		payloadObj := payloads[eventValue]
		info := StorydenSDKEventInfo{
			Event:         eventValue,
			ManifestConst: name,
		}
		if payloadObj != nil {
			info.PayloadType = symbolSummary(pkg.PkgPath, payloadObj, docs[eventValue])
			if typeName, ok := payloadObj.(*types.TypeName); ok {
				if s, ok := typeName.Type().Underlying().(*types.Struct); ok {
					info.Fields = structFields(s)
					info.FieldUsages = sdkFieldUsages(s)
				}
			}
		}
		out.Events = append(out.Events, info)
	}

	return out, nil
}

func (a *Agent) StorydenSDKSearch(ctx context.Context, in StorydenSDKSearchInput) (StorydenSDKSearchResult, error) {
	query := strings.ToLower(strings.TrimSpace(in.Query))
	if query == "" {
		return StorydenSDKSearchResult{}, errors.New("query is required")
	}
	if strings.ContainsAny(query, "*?[](){}^$|\\") {
		return StorydenSDKSearchResult{}, fmt.Errorf("query must be a literal substring, not a regex or wildcard pattern; search for a plain term such as %q", firstPlainSearchTerm(query))
	}

	maxResults := in.MaxResults
	if maxResults <= 0 || maxResults > 300 {
		maxResults = 100
	}

	paths, err := storydenSDKPackagesForArea(in.Area)
	if err != nil {
		return StorydenSDKSearchResult{}, err
	}

	out := StorydenSDKSearchResult{
		Packages: []PackageInfo{},
		Symbols:  []GoSymbolSummary{},
		Hints: []StorydenSDKHint{
			{Message: "Use plugin_storyden_sdk_events for manifest event names and event payload fields."},
			{Message: "Use plugin_go_symbol_detail with the returned import_path and symbol name when you need methods or full struct fields."},
			{Message: "For Storyden host HTTP API calls inside event handlers, use client, err := pl.BuildAPIClient(ctx); do not construct raw openapi clients from plugin internals."},
			{Message: "Prefer generated WithResponse methods when checking HTTP status. If status is non-2xx, return fmt.Errorf with resp.Status() instead of returning nil or a stale err."},
		},
	}

	visited := map[string]bool{}
	for _, importPath := range paths {
		pkg, err := a.loadGoPackage(ctx, importPath)
		if err != nil {
			return StorydenSDKSearchResult{}, err
		}
		out.Packages = append(out.Packages, packageInfoFromPackage(pkg))
		searchStorydenSDKPackageSymbols(pkg, storydenSDKSearchQueries(query), maxResults, &out, visited, in.Area)
		if out.Truncated {
			break
		}
	}

	return out, nil
}

func storydenSDKPackagesForArea(area string) ([]string, error) {
	switch strings.ToLower(strings.TrimSpace(area)) {
	case "", "all":
		return []string{storydenSDKPackage, storydenSDKRPCPackage, storydenSDKOpenAPIPackage, storydenSDKOperationPackage}, nil
	case "plugin", "runtime", "storyden":
		return []string{storydenSDKPackage}, nil
	case "rpc", "events", "manifest":
		return []string{storydenSDKRPCPackage}, nil
	case "http_api", "http", "api", "openapi":
		return []string{storydenSDKOpenAPIPackage}, nil
	case "operations", "operation":
		return []string{storydenSDKOperationPackage}, nil
	default:
		return nil, fmt.Errorf("unknown Storyden SDK area %q; use plugin, rpc, events, manifest, http_api, operations, or all", area)
	}
}

func searchStorydenSDKPackageSymbols(pkg *packages.Package, queries []string, max int, out *StorydenSDKSearchResult, visited map[string]bool, area string) {
	if pkg == nil || pkg.Types == nil || visited[pkg.PkgPath] {
		return
	}
	visited[pkg.PkgPath] = true

	docs := packageDocs(pkg)
	names := pkg.Types.Scope().Names()
	sort.Strings(names)
	for _, name := range names {
		if !ast.IsExported(name) || !storydenSDKAreaAllowsSymbol(area, name) {
			continue
		}
		obj := pkg.Types.Scope().Lookup(name)
		if typeName, ok := obj.(*types.TypeName); ok {
			searchStorydenSDKMethods(pkg.PkgPath, typeName, queries, max, out)
			if out.Truncated {
				return
			}
		}

		summary := symbolSummary(pkg.PkgPath, obj, docs[name])
		haystack := strings.ToLower(summary.ImportPath + " " + summary.Name + " " + summary.Kind + " " + storydenSDKSearchableSignature(obj, summary.Signature) + " " + summary.Doc)
		if !containsAnyLiteral(haystack, queries) {
			continue
		}
		if len(out.Symbols) >= max {
			out.Truncated = true
			return
		}
		out.Symbols = append(out.Symbols, summary)
	}
}

func storydenSDKSearchableSignature(obj types.Object, signature string) string {
	typeName, ok := obj.(*types.TypeName)
	if !ok {
		return signature
	}
	if _, ok := typeName.Type().Underlying().(*types.Interface); ok {
		return ""
	}
	return signature
}

func searchStorydenSDKMethods(importPath string, typeName *types.TypeName, queries []string, max int, out *StorydenSDKSearchResult) {
	named, ok := typeName.Type().(*types.Named)
	if !ok {
		return
	}
	for _, method := range namedMethods(named) {
		if !ast.IsExported(method.Name) {
			continue
		}
		haystack := strings.ToLower(importPath + " " + method.Name + " method " + method.Signature)
		if !containsAnyLiteral(haystack, queries) {
			continue
		}
		if len(out.Symbols) >= max {
			out.Truncated = true
			return
		}
		out.Symbols = append(out.Symbols, GoSymbolSummary{
			ImportPath: importPath,
			Name:       method.Name,
			Kind:       "method",
			Signature:  method.Signature,
		})
	}
}

func sdkFieldUsages(s *types.Struct) []SDKFieldUsage {
	usages := []SDKFieldUsage{}
	for _, field := range structFields(s) {
		if !strings.HasSuffix(field.Name, "ID") {
			continue
		}
		if !strings.HasSuffix(field.Type, ".ID") && field.Type != "xid.ID" {
			continue
		}
		usages = append(usages, SDKFieldUsage{
			Field:      field.Name,
			Type:       field.Type,
			Expression: "event." + field.Name + ".String()",
			Use:        "Use this string form for generated openapi path/query parameters and log messages.",
		})
	}
	return usages
}

func storydenSDKSearchQueries(query string) []string {
	queries := []string{query}
	replacements := map[string]string{
		"reaction":  "react",
		"reactions": "react",
		"replied":   "reply",
		"replies":   "reply",
	}
	for from, to := range replacements {
		if strings.Contains(query, from) {
			queries = append(queries, strings.ReplaceAll(query, from, to))
		}
	}
	for _, term := range strings.Fields(query) {
		queries = append(queries, term)
		if replacement, ok := replacements[term]; ok {
			queries = append(queries, replacement)
		}
	}
	return dedupeStrings(queries)
}

func containsAnyLiteral(haystack string, needles []string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}

func dedupeStrings(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func storydenSDKAreaAllowsSymbol(area, name string) bool {
	switch strings.ToLower(strings.TrimSpace(area)) {
	case "events":
		return strings.HasPrefix(name, "Event")
	case "manifest":
		return strings.Contains(name, "Manifest") ||
			strings.Contains(name, "Capability") ||
			strings.Contains(name, "Config") ||
			strings.Contains(name, "Validate")
	default:
		return true
	}
}
