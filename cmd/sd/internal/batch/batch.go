// Package batch is the per-identifier mutation runner used by node delete,
// node visibility, node move and node update. It defaults to best-effort: it
// keeps going after a per-item failure, reports each result, and exits
// non-zero if anything failed. Agents can re-run only the failed identifiers
// without having to recover state.
package batch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// Result describes the outcome of a single per-identifier operation.
type Result struct {
	ID    string `json:"id"`
	OK    bool   `json:"ok"`
	Name  string `json:"name,omitempty"`
	Error string `json:"error,omitempty"`
}

// Options tunes Run's behaviour. DryRun skips the API call entirely and still
// emits a Result so agents can preview the full plan.
type Options struct {
	DryRun bool
}

// Run executes fn for each id sequentially. The Action is invoked unless DryRun
// is set, in which case Run emits an "ok: true, name: dry-run" Result. Results
// are written through emit so callers can route to stderr (default plain mode)
// or stdout (JSONL mode). The summary line is written through summary at the
// end. Returns true iff every id succeeded.
func Run(ctx context.Context, ids []string, fn func(ctx context.Context, id string) (name string, err error), opts Options, emit func(Result), summary func(succeeded, failed int)) bool {
	succeeded, failed := 0, 0
	for _, id := range ids {
		if id == "" {
			continue
		}
		select {
		case <-ctx.Done():
			emit(Result{ID: id, OK: false, Error: ctx.Err().Error()})
			failed++
			continue
		default:
		}

		if opts.DryRun {
			emit(Result{ID: id, OK: true, Name: "(dry run)"})
			succeeded++
			continue
		}

		name, err := fn(ctx, id)
		if err != nil {
			emit(Result{ID: id, OK: false, Error: err.Error()})
			failed++
			continue
		}
		emit(Result{ID: id, OK: true, Name: name})
		succeeded++
	}
	if summary != nil {
		summary(succeeded, failed)
	}
	return failed == 0
}

// JSONLEmitter returns an emit function that writes each Result as one JSON
// object per line to w.
func JSONLEmitter(w io.Writer) func(Result) {
	encoder := json.NewEncoder(w)
	return func(r Result) {
		_ = encoder.Encode(r)
	}
}

// PlainEmitter returns an emit function that writes a one-line human summary
// of each Result to w.
func PlainEmitter(w io.Writer) func(Result) {
	return func(r Result) {
		if r.OK {
			if r.Name != "" {
				fmt.Fprintf(w, "ok    %s  %s\n", r.ID, r.Name)
				return
			}
			fmt.Fprintf(w, "ok    %s\n", r.ID)
			return
		}
		fmt.Fprintf(w, "fail  %s  %s\n", r.ID, r.Error)
	}
}

// PlainSummary returns a summary function that writes "N succeeded, M failed"
// to w. Suppressed when both counts are zero.
func PlainSummary(w io.Writer) func(int, int) {
	return func(succeeded, failed int) {
		if succeeded == 0 && failed == 0 {
			return
		}
		fmt.Fprintf(w, "%d succeeded, %d failed\n", succeeded, failed)
	}
}
