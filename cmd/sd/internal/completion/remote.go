package completion

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/carapace-sh/carapace"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

const lookupTimeout = 1500 * time.Millisecond
const maxResponseBytes = 2 << 20
const maxSuggestions = 100

func operation(doc *openapi3.T, id string) (string, *openapi3.Operation, openapi3.Parameters) {
	for path, item := range doc.Paths.Map() {
		for _, op := range item.Operations() {
			if op.OperationID == id {
				return path, op, append(append(openapi3.Parameters{}, item.Parameters...), op.Parameters...)
			}
		}
	}
	return "", nil, nil
}

func resourceAction(cmd *cobra.Command, store *config.Store, def command, src source) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		params := map[string]string{}
		for i, in := range def.Arguments {
			if i >= len(c.Args) {
				break
			}
			name := strings.ToLower(in.Name)
			if name == "slug" {
				name = "node_slug"
			}
			params[name] = c.Args[i]
		}
		return lookupAction(cmd, store, src, params, c.Value)
	})
}

func lookupAction(cmd *cobra.Command, store *config.Store, src source, params map[string]string, prefix string) carapace.Action {
	path, ok := lookupPath(src, params, prefix)
	if !ok {
		return carapace.ActionValues()
	}
	ctx, cancel := context.WithTimeout(cmd.Root().Context(), lookupTimeout)
	defer cancel()
	// Carapace parses flags without running the leaf's PersistentPreRunE. Copy the
	// store so the callback honors --context without altering the saved default.
	selected := *store
	if f := cmd.Flag("context"); f != nil {
		selected.SelectedContext = f.Value.String()
	}
	client, err := api.NewAuthenticatedClient(ctx, &selected, api.WithRateLimitWarnings(io.Discard))
	if err != nil {
		return carapace.ActionValues()
	}
	endpoint := strings.TrimRight(client.BaseURL, "/") + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return carapace.ActionValues()
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req, io.Discard)
	if err != nil {
		return carapace.ActionValues()
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return carapace.ActionValues()
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil || len(data) > maxResponseBytes {
		return carapace.ActionValues()
	}
	var body map[string]any
	if json.Unmarshal(data, &body) != nil {
		return carapace.ActionValues()
	}
	return carapace.ActionValuesDescribed(suggestions(body[src.collection], src)...)
}

// lookupPath binds only allowlisted GET operations to previously completed inputs.
func lookupPath(src source, params map[string]string, prefix string) (string, bool) {
	doc, err := openapi.GetSwagger()
	if err != nil {
		return "", false
	}
	path, op, parameters := operation(doc, src.operation)
	// A source must be an explicitly selected GET, even if the API changes later.
	if op == nil || doc.Paths.Value(path).Get != op {
		return "", false
	}
	query := url.Values{}
	for _, ref := range parameters {
		p := ref.Value
		if p.In == "path" {
			value := params[p.Name]
			if value == "" || value == "." || value == ".." {
				return "", false
			}
			path = strings.ReplaceAll(path, "{"+p.Name+"}", url.PathEscape(value))
		}
		if p.In == "query" && p.Name == "q" && src.value != "id" && prefix != "" {
			query.Set("q", prefix)
		}
	}
	if src.operation == "NodeList" {
		query.Set("format", "flat")
	}
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	return path, true
}

func suggestions(items any, src source) []string {
	var values []string
	seen := map[string]bool{}
	var visit func(any)
	visit = func(items any) {
		rows, _ := items.([]any)
		for _, row := range rows {
			if len(values)/2 >= maxSuggestions {
				return
			}
			item, ok := row.(map[string]any)
			if !ok {
				continue
			}
			value, _ := item[src.value].(string)
			if value != "" && len(value) <= 512 && !strings.ContainsFunc(value, unicode.IsControl) && !seen[value] {
				description, _ := item[src.description].(string)
				description = strings.Map(func(r rune) rune {
					if unicode.IsControl(r) {
						return ' '
					}
					return r
				}, description)
				if r := []rune(description); len(r) > 120 {
					description = string(r[:120]) + "…"
				}
				values = append(values, value, description)
				seen[value] = true
			}
			if src.children != "" {
				visit(item[src.children])
			}
		}
	}
	visit(items)
	return values
}
