package batch

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunBestEffortContinuesOnFailure(t *testing.T) {
	r := require.New(t)
	var results []Result
	emit := func(res Result) { results = append(results, res) }

	ok := Run(context.Background(),
		[]string{"a", "b", "c"},
		func(ctx context.Context, id string) (string, error) {
			if id == "b" {
				return "", errors.New("boom")
			}
			return id + "-name", nil
		},
		Options{},
		emit,
		nil,
	)

	r.False(ok)
	r.Len(results, 3)
	r.True(results[0].OK)
	r.False(results[1].OK)
	r.Equal("boom", results[1].Error)
	r.True(results[2].OK)
}

func TestRunDryRunSkipsCall(t *testing.T) {
	r := require.New(t)
	called := 0
	var results []Result

	ok := Run(context.Background(),
		[]string{"a", "b"},
		func(ctx context.Context, id string) (string, error) {
			called++
			return "", nil
		},
		Options{DryRun: true},
		func(res Result) { results = append(results, res) },
		nil,
	)

	r.True(ok)
	r.Equal(0, called)
	r.Len(results, 2)
	for _, res := range results {
		r.True(res.OK)
		r.Equal("(dry run)", res.Name)
	}
}

func TestRunSkipsEmptyIDs(t *testing.T) {
	r := require.New(t)
	var results []Result
	Run(context.Background(),
		[]string{"", "a", ""},
		func(ctx context.Context, id string) (string, error) { return id, nil },
		Options{},
		func(res Result) { results = append(results, res) },
		nil,
	)
	r.Len(results, 1)
	r.Equal("a", results[0].ID)
}

func TestRunReportsSummary(t *testing.T) {
	r := require.New(t)
	var s, f int
	Run(context.Background(),
		[]string{"ok", "bad"},
		func(ctx context.Context, id string) (string, error) {
			if id == "bad" {
				return "", errors.New("x")
			}
			return id, nil
		},
		Options{},
		func(Result) {},
		func(succeeded, failed int) { s, f = succeeded, failed },
	)
	r.Equal(1, s)
	r.Equal(1, f)
}

func TestJSONLEmitterShape(t *testing.T) {
	r := require.New(t)
	var buf bytes.Buffer
	emit := JSONLEmitter(&buf)
	emit(Result{ID: "a", OK: true, Name: "Alpha"})
	emit(Result{ID: "b", OK: false, Error: "x"})

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	r.Len(lines, 2)
	r.JSONEq(`{"id":"a","ok":true,"name":"Alpha"}`, lines[0])
	r.JSONEq(`{"id":"b","ok":false,"error":"x"}`, lines[1])
}

func TestPlainEmitterShape(t *testing.T) {
	r := require.New(t)
	var buf bytes.Buffer
	emit := PlainEmitter(&buf)
	emit(Result{ID: "a", OK: true, Name: "Alpha"})
	emit(Result{ID: "b", OK: false, Error: "x"})
	emit(Result{ID: "c", OK: true})

	out := buf.String()
	r.Contains(out, "ok    a  Alpha")
	r.Contains(out, "fail  b  x")
	r.Contains(out, "ok    c\n")
}
