package tracing

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	coltracev1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	tracev1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func newOTLPJSONLClient(path string) otlptrace.Client {
	return &otlpJSONLClient{
		path: path,
	}
}

type otlpJSONLClient struct {
	path string

	mu   sync.Mutex
	file *os.File
}

func (c *otlpJSONLClient) Start(context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.file != nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return err
	}

	file, err := os.OpenFile(c.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}

	c.file = file

	return nil
}

func (c *otlpJSONLClient) Stop(context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.file == nil {
		return nil
	}

	err := c.file.Close()
	c.file = nil

	return err
}

func (c *otlpJSONLClient) UploadTraces(_ context.Context, protoSpans []*tracev1.ResourceSpans) error {
	payload, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(
		&coltracev1.ExportTraceServiceRequest{
			ResourceSpans: protoSpans,
		},
	)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.file == nil {
		return errors.New("otlp jsonl client not started")
	}

	if _, err := c.file.Write(payload); err != nil {
		return err
	}

	if _, err := c.file.WriteString("\n"); err != nil {
		return err
	}

	return nil
}
