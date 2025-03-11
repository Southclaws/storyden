package kv

import (
	"time"

	"github.com/Southclaws/dt"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
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

func (a Attrs) ToZap() []zap.Field {
	return dt.Map(a, func(attr Attr) zap.Field {
		switch attr.Value.Type() {
		case attribute.BOOL:
			return zap.Bool(string(attr.Key), attr.Value.AsBool())

		case attribute.INT64:
			return zap.Int64(string(attr.Key), attr.Value.AsInt64())

		case attribute.FLOAT64:
			return zap.Float64(string(attr.Key), attr.Value.AsFloat64())

		case attribute.STRING:
			return zap.String(string(attr.Key), attr.Value.AsString())

		case attribute.BOOLSLICE:
			return zap.Bools(string(attr.Key), attr.Value.AsBoolSlice())

		case attribute.INT64SLICE:
			return zap.Int64s(string(attr.Key), attr.Value.AsInt64Slice())

		case attribute.FLOAT64SLICE:
			return zap.Float64s(string(attr.Key), attr.Value.AsFloat64Slice())

		case attribute.STRINGSLICE:
			return zap.Strings(string(attr.Key), attr.Value.AsStringSlice())
		}

		return zap.Any(string(attr.Key), attr.Value.Emit())
	})
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

func Time(key string, d time.Time) Attr {
	return Attr(attribute.String(key, d.Format(time.RFC3339)))
}

func Duration(key string, d time.Duration) Attr {
	return Attr(attribute.String(key, d.String()))
}
