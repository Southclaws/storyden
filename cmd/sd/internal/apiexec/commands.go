package apiexec

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/rs/xid"
	"github.com/spf13/cobra"
)

func NewOperations(e *Executor) cligen.ListApiHandler {
	return func(ctx context.Context, cmd *cobra.Command, streams cligen.IO, p cligen.ListApiParams) error {
		return e.Operations(streams.Out)
	}
}
func NewSchema(e *Executor) cligen.DescribeApiHandler {
	return func(ctx context.Context, cmd *cobra.Command, streams cligen.IO, p cligen.DescribeApiParams) error {
		return e.Schema(streams.Out, p.Operation)
	}
}
func NewRequest(e *Executor) cligen.RequestApiHandler {
	return func(ctx context.Context, cmd *cobra.Command, streams cligen.IO, p cligen.RequestApiParams) error {
		path, err := parsePairs(p.Path)
		if err != nil {
			return err
		}
		paths := map[string]string{}
		for k, values := range path {
			if len(values) != 1 {
				return fmt.Errorf("path parameter %q must appear once", k)
			}
			paths[k] = values[0]
		}
		query, err := parsePairs(p.Query)
		if err != nil {
			return err
		}
		return e.Run(ctx, streams, Request{Operation: p.Operation, Path: paths, Query: query, Data: p.Data, File: p.File, ContentType: p.ContentType, Output: string(p.Output), OutputFile: p.OutputFile, Timeout: p.Timeout})
	}
}
func parsePairs(pairs []string) (url.Values, error) {
	values := url.Values{}
	for _, pair := range pairs {
		key, value, ok := strings.Cut(pair, "=")
		if !ok || key == "" {
			return nil, fmt.Errorf("expected KEY=VALUE, got %q", pair)
		}
		values.Add(key, value)
	}
	return values, nil
}

func NewChat(e *Executor) cligen.RobotChatHandler {
	return func(ctx context.Context, cmd *cobra.Command, streams cligen.IO, p cligen.RobotChatParams) error {
		if p.Message != "" && p.MessageFile != "" {
			return fmt.Errorf("use MESSAGE or --message-file, not both")
		}
		message := p.Message
		if p.MessageFile != "" {
			var data []byte
			var err error
			if p.MessageFile == "-" {
				data, err = io.ReadAll(streams.In)
			} else {
				data, err = os.ReadFile(p.MessageFile)
			}
			if err != nil {
				return err
			}
			message = string(data)
		}
		if strings.TrimSpace(message) == "" {
			return fmt.Errorf("provide a nonempty MESSAGE or --message-file")
		}
		session := p.Session
		if session == "" {
			session = xid.New().String()
		} else if _, err := xid.FromString(session); err != nil {
			return fmt.Errorf("--session must be a valid session XID")
		}
		payload := map[string]any{"id": session, "sessionId": session, "messages": []any{map[string]any{"id": xid.New().String(), "role": "user", "parts": []any{map[string]any{"type": "text", "text": message}}}}}
		if p.Robot != "" {
			payload["robotId"] = p.Robot
		}
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		return e.Run(ctx, streams, Request{Operation: "RobotSessionCreate", Data: string(data), Output: string(p.Output), Timeout: p.Timeout})
	}
}
