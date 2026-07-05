package pluginbuilder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
	storydensdk "github.com/Southclaws/storyden/sdk/go/storyden"
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
	HandlerMethod string          `json:"handler_method"`
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
	_ = ctx
	query := strings.ToLower(strings.TrimSpace(in.Query))
	if strings.ContainsAny(query, "*?[](){}^$|\\") {
		return StorydenSDKEventsResult{}, fmt.Errorf("query must be a literal substring, not a regex or wildcard pattern; search for a plain term such as %q", plugindev.FirstPlainSearchTerm(query))
	}

	maxEvents := in.MaxEvents
	if maxEvents <= 0 || maxEvents > 200 {
		maxEvents = 60
	}

	out := StorydenSDKEventsResult{
		ImportPath: storydenSDKRPCPackage,
		Events:     []StorydenSDKEventInfo{},
		Hints: []StorydenSDKHint{
			{Message: "Use event values such as EventThreadReplyCreated in manifest.yaml events_consumed."},
			{Message: "Use payload type names such as *rpc.EventThreadReplyCreated in event handler signatures."},
			{Message: "Register event handlers with the returned handler_method, for example pl.OnThreadReplyCreated(...). Do not call nonexistent methods such as HandleEventRPC."},
			{Message: "Activity events are Storyden activity records, not Discord gateway messages or Discord user activity. For Discord message handling, register handlers with the Discord client and do not add Storyden Activity events unless the requested behavior is explicitly about Storyden activity records."},
			{Message: "When passing Storyden resource IDs to generated openapi parameters, use event.Field.String(); never use string(event.Field[:])."},
		},
	}

	for _, event := range rpc.EventValues {
		eventValue := string(event)
		name := "Event" + eventValue
		if query != "" && !strings.Contains(strings.ToLower(name+" "+eventValue), query) {
			continue
		}
		if len(out.Events) >= maxEvents {
			out.Truncated = true
			break
		}

		out.Events = append(out.Events, storydenSDKEventInfo(eventValue))
	}

	return out, nil
}

func (a *Agent) StorydenSDKSearch(ctx context.Context, in StorydenSDKSearchInput) (StorydenSDKSearchResult, error) {
	query := strings.ToLower(strings.TrimSpace(in.Query))
	if query == "" {
		return StorydenSDKSearchResult{}, errors.New("query is required")
	}
	if strings.ContainsAny(query, "*?[](){}^$|\\") {
		return StorydenSDKSearchResult{}, fmt.Errorf("query must be a literal substring, not a regex or wildcard pattern; search for a plain term such as %q", plugindev.FirstPlainSearchTerm(query))
	}

	maxResults := in.MaxResults
	if maxResults <= 0 || maxResults > 300 {
		maxResults = 100
	}

	if result, ok := staticStorydenSDKSearch(query, in.Area, maxResults); ok {
		return result, nil
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
			{Message: "For Storyden host HTTP API calls, use client, err := pl.BuildAPIClient(ctx) and reuse the resulting client where appropriate; do not construct raw API clients from plugin internals."},
			{Message: "If code uses BuildAPIClient, manifest.yaml must include access with a stable bot account handle, display name, and narrow Storyden permission names for the API operations being called."},
			{Message: "Prefer generated WithResponse methods when checking HTTP status. If status is non-2xx, return fmt.Errorf with resp.Status() instead of returning nil or a stale err."},
		},
	}

	visited := map[string]bool{}
	for _, importPath := range paths {
		pkg, err := a.loadGoPackage(ctx, importPath)
		if err != nil {
			return StorydenSDKSearchResult{}, err
		}
		out.Packages = append(out.Packages, plugindev.PackageInfoFromPackage(pkg))
		searchStorydenSDKPackageSymbols(pkg, storydenSDKSearchQueries(query), maxResults, &out, visited, in.Area)
		if out.Truncated {
			break
		}
	}

	return out, nil
}

func staticStorydenSDKSearch(query, area string, maxResults int) (StorydenSDKSearchResult, bool) {
	normalisedArea := strings.ToLower(strings.TrimSpace(area))
	if isRobotRunQuery(query) {
		switch normalisedArea {
		case "", "all", "plugin", "runtime", "storyden", "rpc", "http_api", "http", "api", "openapi", "operations", "operation":
			return staticStorydenPluginSearch(query, maxResults), true
		}
	}

	switch normalisedArea {
	case "plugin", "runtime", "storyden":
		return staticStorydenPluginSearch(query, maxResults), true
	case "events":
		return staticStorydenEventSearch(query, maxResults), true
	default:
		return StorydenSDKSearchResult{}, false
	}
}

func isRobotRunQuery(query string) bool {
	normalised := strings.ToLower(strings.TrimSpace(query))
	if normalised == "run" || normalised == "robot" || normalised == "robots" {
		return true
	}
	return strings.Contains(normalised, "robot_run") ||
		strings.Contains(normalised, "runrobot") ||
		strings.Contains(normalised, "run robot") ||
		strings.Contains(normalised, "robot run") ||
		strings.Contains(normalised, "run a robot") ||
		strings.Contains(normalised, "call robot")
}

func staticStorydenPluginSearch(query string, maxResults int) StorydenSDKSearchResult {
	out := StorydenSDKSearchResult{
		Packages: []PackageInfo{{
			ImportPath: storydenSDKPackage,
			Name:       "storyden",
		}},
		Symbols: []GoSymbolSummary{},
		Hints: []StorydenSDKHint{
			{Message: "Use plugin_storyden_sdk_events for manifest event names, event payload fields, and event handler methods."},
			{Message: "Use pl.RunRobot(ctx, robotID, message) for robot_run; manifest access.permissions must include USE_ROBOTS."},
			{Message: "Do not use generated HTTP RobotChatSSE, RobotChatSSEWithResponse, or UI chat streaming endpoints from plugins; those are for UI clients, not plugin-to-host robot execution."},
			{Message: "For Storyden host HTTP API calls, use client, err := pl.BuildAPIClient(ctx) and reuse the resulting client where appropriate; do not construct raw API clients from plugin internals."},
			{Message: "If code uses BuildAPIClient or RunRobot, manifest.yaml must include access with a stable bot account handle, display name, and narrow Storyden permission names for the API operations being called."},
		},
	}

	queries := storydenSDKSearchQueries(query)
	pluginType := reflect.TypeOf((*storydensdk.Plugin)(nil))
	for i := 0; i < pluginType.NumMethod(); i++ {
		method := pluginType.Method(i)
		haystack := strings.ToLower(storydenSDKPackage + " " + method.Name + " " + camelToSnake(method.Name) + " method " + method.Type.String())
		if !containsAnyLiteral(haystack, queries) {
			continue
		}
		if len(out.Symbols) >= maxResults {
			out.Truncated = true
			break
		}
		out.Symbols = append(out.Symbols, GoSymbolSummary{
			ImportPath: storydenSDKPackage,
			Name:       method.Name,
			Kind:       "method",
			Signature:  method.Type.String(),
		})
	}

	return out
}

func staticStorydenEventSearch(query string, maxResults int) StorydenSDKSearchResult {
	out := StorydenSDKSearchResult{
		Packages: []PackageInfo{{
			ImportPath: storydenSDKRPCPackage,
			Name:       "rpc",
		}},
		Symbols: []GoSymbolSummary{},
		Hints: []StorydenSDKHint{
			{Message: "Use plugin_storyden_sdk_events for complete event records including handler_method and payload fields."},
			{Message: "Activity events are Storyden activity records, not Discord gateway messages or Discord user activity."},
		},
	}

	for _, event := range rpc.EventValues {
		eventValue := string(event)
		info := storydenSDKEventInfo(eventValue)
		haystack := strings.ToLower(info.Event + " " + info.ManifestConst + " " + info.HandlerMethod + " " + info.PayloadType.Name)
		if !containsAnyLiteral(haystack, storydenSDKSearchQueries(query)) {
			continue
		}
		if len(out.Symbols) >= maxResults {
			out.Truncated = true
			break
		}
		out.Symbols = append(out.Symbols, info.PayloadType)
	}

	return out
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

	docs := plugindev.PackageDocs(pkg)
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

		summary := plugindev.SymbolSummary(pkg.PkgPath, obj, docs[name])
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
	for _, method := range plugindev.NamedMethods(named) {
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
	for _, field := range plugindev.StructFields(s) {
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

func storydenSDKEventInfo(event string) StorydenSDKEventInfo {
	info := StorydenSDKEventInfo{
		Event:         event,
		ManifestConst: "Event" + event,
		HandlerMethod: storydenEventHandlerMethod(event),
		PayloadType: GoSymbolSummary{
			ImportPath: storydenSDKRPCPackage,
			Name:       event,
			Kind:       "type",
			Signature:  "struct",
		},
	}

	payloadType := eventPayloadType(event)
	if payloadType == nil {
		return info
	}

	info.PayloadType.Name = payloadType.Name()
	info.Fields = reflectStructFields(payloadType)
	info.FieldUsages = reflectSDKFieldUsages(info.Fields)
	return info
}

func eventPayloadType(event string) reflect.Type {
	data, err := json.Marshal(map[string]any{"event": event})
	if err != nil {
		return nil
	}

	var payload rpc.EventPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil
	}
	if payload.EventPayloadUnion == nil {
		return nil
	}

	t := reflect.TypeOf(payload.EventPayloadUnion)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	return t
}

func reflectStructFields(t reflect.Type) []GoFieldInfo {
	fields := []GoFieldInfo{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		fields = append(fields, GoFieldInfo{
			Name: field.Name,
			Type: field.Type.String(),
			Tag:  string(field.Tag),
		})
	}
	return fields
}

func reflectSDKFieldUsages(fields []GoFieldInfo) []SDKFieldUsage {
	usages := []SDKFieldUsage{}
	for _, field := range fields {
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

func storydenEventHandlerMethod(event string) string {
	return "On" + strings.TrimPrefix(event, "Event")
}

func storydenSDKSearchQueries(query string) []string {
	queries := []string{query}
	replacements := map[string]string{
		"reaction":  "react",
		"reactions": "react",
		"replied":   "reply",
		"replies":   "reply",
		"robot_run": "runrobot",
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

func camelToSnake(value string) string {
	var b strings.Builder
	for i, r := range value {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
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
