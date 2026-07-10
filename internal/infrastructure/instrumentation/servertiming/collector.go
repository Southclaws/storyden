package servertiming

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	maxMetrics        = 20
	minMetricDuration = 5 * time.Millisecond
)

type Collector struct {
	mu      sync.Mutex
	metrics map[string]metric
}

type metric struct {
	group string
	desc  string
	durMS float64
}

func NewCollector() *Collector {
	return &Collector{metrics: map[string]metric{}}
}

func (c *Collector) Observe(name string, duration time.Duration) {
	if duration < minMetricDuration {
		return
	}

	group, desc := normaliseMetricName(name)
	if group == "" || desc == "" {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	key := group + "\x00" + desc
	m := c.metrics[key]
	m.group = group
	m.desc = desc
	m.durMS += float64(duration) / float64(time.Millisecond)
	c.metrics[key] = m
}

func (c *Collector) HeaderValue() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.metrics) == 0 {
		return ""
	}

	metrics := make([]metric, 0, len(c.metrics))
	for _, m := range c.metrics {
		metrics = append(metrics, m)
	}

	sort.Slice(metrics, func(i, j int) bool {
		if metrics[i].durMS == metrics[j].durMS {
			if metrics[i].group == metrics[j].group {
				return metrics[i].desc < metrics[j].desc
			}
			return metrics[i].group < metrics[j].group
		}
		return metrics[i].durMS > metrics[j].durMS
	})

	if len(metrics) > maxMetrics {
		metrics = metrics[:maxMetrics]
	}

	parts := make([]string, 0, len(metrics))
	for _, m := range metrics {
		dur := strconv.FormatFloat(m.durMS, 'f', 2, 64)
		parts = append(parts, fmt.Sprintf("%s;dur=%s;desc=\"%s\"", m.group, dur, escapeDescription(m.desc)))
	}

	return strings.Join(parts, ", ")
}

func normaliseMetricName(name string) (string, string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ""
	}

	if group, desc, ok := normaliseStructuredMetric(name); ok {
		return group, desc
	}

	return normaliseFlatMetric(name)
}

func normaliseStructuredMetric(name string) (string, string, bool) {
	firstDot := strings.Index(name, ".")
	if firstDot <= 0 || firstDot >= len(name)-1 {
		return "", "", false
	}

	pkg := name[:firstDot]
	rest := name[firstDot+1:]

	var thing, method string
	switch {
	case strings.HasPrefix(rest, "("):
		closeIndex := strings.Index(rest, ").")
		if closeIndex <= 1 || closeIndex >= len(rest)-2 {
			return "", "", false
		}
		thing = rest[1:closeIndex]
		method = rest[closeIndex+2:]
	default:
		secondDot := strings.Index(rest, ".")
		if secondDot <= 0 || secondDot >= len(rest)-1 {
			return "", "", false
		}
		thing = rest[:secondDot]
		method = rest[secondDot+1:]
	}

	group, suffix := splitPackageName(pkg)
	group = sanitiseToken(group)
	thing = normaliseDescriptionComponent(strings.TrimPrefix(thing, "*"), true)
	if thing == "" {
		thing = normaliseDescriptionComponent(suffix, true)
	}

	method = normaliseMethodName(method)
	if group == "" || thing == "" || method == "" {
		return "", "", false
	}

	return group, thing + "." + method, true
}

func normaliseFlatMetric(name string) (string, string) {
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '/' || r == '.'
	})
	if len(parts) == 0 {
		return "", ""
	}

	group, suffix := splitPackageName(parts[0])
	group = sanitiseToken(group)

	descParts := make([]string, 0, len(parts))
	if len(parts) > 1 {
		for _, part := range parts[1:] {
			part = normaliseDescriptionComponent(part, false)
			if part == "" {
				continue
			}
			descParts = append(descParts, part)
		}
	} else {
		suffix = normaliseDescriptionComponent(suffix, true)
		if suffix != "" {
			descParts = append(descParts, suffix)
		}
	}

	if group == "" || len(descParts) == 0 {
		return "", ""
	}

	return group, strings.Join(descParts, ".")
}

func splitPackageName(value string) (string, string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ""
	}

	if slash := strings.LastIndex(value, "/"); slash >= 0 {
		value = value[slash+1:]
	}

	value = strings.Trim(value, "_")
	if value == "" {
		return "", ""
	}

	group, suffix, found := strings.Cut(value, "_")
	if !found {
		return value, ""
	}

	return group, strings.Trim(suffix, "_")
}

func normaliseMethodName(value string) string {
	if cut := strings.Index(value, "-"); cut > 0 {
		value = value[:cut]
	}

	return normaliseDescriptionComponent(value, false)
}

func normaliseDescriptionComponent(value string, toLower bool) string {
	if value == "" {
		return ""
	}

	value = strings.Trim(value, "*_() ")
	if value == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(value))
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_' || r == '.':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}

	value = strings.Trim(b.String(), "_")
	if toLower {
		value = strings.ToLower(value)
	}

	return value
}

func escapeDescription(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return value
}

func sanitiseToken(name string) string {
	if name == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(name))

	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case strings.ContainsRune("!#$%&'*+-.^_`|~", r):
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}

	return strings.Trim(b.String(), "_")
}
