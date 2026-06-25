package pluginbuilder

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"

	"github.com/uber-go/gopatch/patch"
)

const maxPatchTargetBytes = 512_000

type ApplyPatchInput struct {
	Path  string `json:"path" jsonschema:"Workspace-relative Go source file path to patch. Must end with .go."`
	Patch string `json:"patch" jsonschema:"gopatch semantic Go patch. Do not use git diff headers or line-number hunks."`
}

type ApplyPatchResult struct {
	Path    string         `json:"path"`
	Changed bool           `json:"changed"`
	Bytes   int            `json:"bytes"`
	Message string         `json:"message,omitempty"`
	Vet     *CommandResult `json:"vet,omitempty"`
}

func (a *Agent) addPatchTools(add toolAdder) error {
	return add(functiontool.New(functiontool.Config{
		Name:        "plugin_apply_patch",
		Description: applyPatchDescription,
	}, func(ctx adktool.Context, args ApplyPatchInput) (ApplyPatchResult, error) {
		result, err := a.ApplyPatch(ctx, args)
		if err != nil {
			return ApplyPatchResult{}, err
		}
		return result, nil
	}))
}

const applyPatchDescription = `Apply a semantic Go source patch to one workspace-relative .go file.

This tool does not run git and does not accept git/unified diff file headers.
Use it for precise edits to existing Go files. Use plugin_file_write for new files,
manifest.yaml, go.mod, or full-file rewrites.

When the patch changes the file, this tool automatically runs go vet ./... and
returns the vet result. Use that feedback immediately before making another edit.

Patch format:
` + "```" + `
@@
optional metavariable declarations
@@
-Go code pattern to match
+replacement Go code
` + "```" + `

Rules:
- Always set path to the exact .go file being edited.
- Do not include lines like "diff --git", "--- a/file.go", "+++ b/file.go", or "@@ -1,3 +1,3 @@".
- The minus and plus sections must be valid Go syntax patterns. They are not line-based text hunks.
- Match complete statements or complete expressions. Do not patch half of a function call, half of a function signature, or an unclosed block.
- Use surrounding statement context when inserting lines inside a function.
- Use "..." inside calls, argument lists, statement lists, or composite literals to elide zero or more items.
- Declare metavariables when reusing matched expressions, statements, or identifiers.
- If a patch becomes hard to express, use plugin_file_read followed by plugin_file_write for a complete file rewrite.

Valid minimal insertion before an existing complete statement:
` + "```" + `
@@
@@
-pl, err := storyden.New(ctx)
+log.Println("plugin starting")
+pl, err := storyden.New(ctx)
` + "```" + `

Invalid line-diff-style insertion into an unclosed block:
` + "```" + `
@@
@@
-pl.OnThreadPublished(func(ctx context.Context, event *rpc.EventThreadPublished) error {
+pl.OnThreadPublished(func(ctx context.Context, event *rpc.EventThreadPublished) error {
+	log.Println("event")
` + "```" + `

The invalid example fails because the minus and plus patterns are incomplete Go
syntax. Include the full statement or patch a complete statement inside the block.

Example: replace one exact return value in main.go
` + "```" + `
@@
@@
-return "old"
+return "new"
` + "```" + `

Example: wrap all matching expressions while preserving the expression
` + "```" + `
@@
var x expression
@@
-log.Println(x)
+logger.Info(x)
` + "```" + `

Example: add work before an existing call while preserving all arguments
` + "```" + `
@@
@@
-client.ReplyCreateWithResponse(ctx, ...)
+if ctx.Err() != nil {
+	return ctx.Err()
+}
+client.ReplyCreateWithResponse(ctx, ...)
` + "```" + `

In a function call, "..." matches zero or more arguments. In a statement list or
block, "..." matches zero or more statements. It can match complex expressions;
you do not need to name every argument when the exact values are not changing.

The result reports changed=false if the patch was valid but did not match the file.
Syntax problems in the patch return an error instead of changed=false.
When changed=true, inspect vet.success, vet.output, and vet.error for immediate
compile/type/line-number feedback.`

func (a *Agent) ApplyPatch(ctx context.Context, in ApplyPatchInput) (ApplyPatchResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return ApplyPatchResult{}, err
	}
	path := strings.TrimSpace(in.Path)
	if path == "" {
		return ApplyPatchResult{}, errors.New("path is required")
	}
	if filepath.Ext(path) != ".go" {
		return ApplyPatchResult{}, fmt.Errorf("plugin_apply_patch only supports Go source files: %s", path)
	}
	if strings.TrimSpace(in.Patch) == "" {
		return ApplyPatchResult{}, errors.New("patch is required")
	}
	if strings.Contains(in.Patch, "\x00") {
		return ApplyPatchResult{}, errors.New("patch contains NUL byte")
	}

	source, err := workspace.ReadFile(ctx, path, maxPatchTargetBytes)
	if err != nil {
		return ApplyPatchResult{}, err
	}
	if source.Truncated {
		return ApplyPatchResult{}, fmt.Errorf("file %q exceeds %d byte patch limit", source.Path, maxPatchTargetBytes)
	}

	parsed, err := patch.Parse("plugin_builder.patch", []byte(in.Patch))
	if err != nil {
		return ApplyPatchResult{}, fmt.Errorf("patch parse error: %w. Ensure minus and plus sections are complete valid Go syntax patterns, not line-based hunks", err)
	}

	next, err := parsed.Apply(source.Path, source.Content)
	if err != nil {
		return ApplyPatchResult{}, fmt.Errorf("patch apply error: %w", err)
	}

	if bytes.Equal(source.Content, next) {
		return ApplyPatchResult{
			Path:    source.Path,
			Changed: false,
			Bytes:   len(source.Content),
			Message: "patch was valid but did not match any code in the target file",
		}, nil
	}

	written, err := workspace.WriteFile(ctx, source.Path, next)
	if err != nil {
		return ApplyPatchResult{}, err
	}

	vet, err := commandResult(workspace.Run(ctx, workspaceprovider.CommandSpec{Command: "go", Args: []string{"vet", "./..."}}))
	if err != nil {
		return ApplyPatchResult{}, err
	}

	return ApplyPatchResult{
		Path:    written.Path,
		Changed: true,
		Bytes:   written.Bytes,
		Message: "patch applied and go vet ./... completed; inspect vet for validation details",
		Vet:     &vet,
	}, nil
}
