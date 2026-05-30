package batch

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ReadIdentifiers reads one identifier per line from r. Lines may be plain
// slugs/xids or JSON objects with a `.slug` or `.id` field (matching the
// shape `sd node list --format jsonl` emits). Blank lines and JSON objects
// with neither field are skipped without erroring so users can pipe lossy
// streams.
func ReadIdentifiers(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var ids []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		id, ok, err := extractID(line)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		ids = append(ids, id)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func extractID(line string) (string, bool, error) {
	if strings.HasPrefix(line, "{") {
		var probe struct {
			ID   string `json:"id"`
			Slug string `json:"slug"`
		}
		if err := json.Unmarshal([]byte(line), &probe); err != nil {
			return "", false, fmt.Errorf("invalid JSONL: %w", err)
		}
		// Prefer slug because it's the most readable identifier; xid works in
		// the same code paths but slug round-trips through copy/paste better.
		if probe.Slug != "" {
			return probe.Slug, true, nil
		}
		if probe.ID != "" {
			return probe.ID, true, nil
		}
		return "", false, nil
	}
	return line, true, nil
}
