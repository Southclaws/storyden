package kv

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
)

type Attr attribute.KeyValue

type Attrs []Attr

func (a Attrs) ToAttributes() []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, len(a))
	for i, attr := range a {
		attrs[i] = attribute.KeyValue(attr)
	}
	return attrs
}

func (a Attrs) ToFault() []string {
	out := make([]string, 0, len(a)*2)
	for _, attr := range a {
		out = append(out, string(attr.Key), attr.Value.Emit())
	}
	return out
}

func String(key, value string) Attr {
	return Attr(attribute.String(key, value))
}

func Int(key string, value int) Attr {
	return Attr(attribute.Int(key, value))
}

func Float(key string, value float64) Attr {
	return Attr(attribute.Float64(key, value))
}

func Bool(key string, value bool) Attr {
	return Attr(attribute.Bool(key, value))
}

func Strings(key string, value []string) Attr {
	return Attr(attribute.StringSlice(key, value))
}

func Duration(key string, d time.Duration) Attr {
	return Attr(attribute.String(key, d.String()))
}
