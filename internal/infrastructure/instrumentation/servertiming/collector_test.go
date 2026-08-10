package servertiming

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCollectorHeaderValue(t *testing.T) {
	t.Parallel()

	c := NewCollector()
	c.Observe("thread.(*service).List", 8*time.Millisecond)
	c.Observe("thread.(*service).List", 7*time.Millisecond)
	c.Observe("thread_querier.(*Querier).List", 9*time.Millisecond)
	c.Observe("thread.(*service).Get", 4*time.Millisecond)

	header := c.HeaderValue()

	assert.Contains(t, header, "thread;dur=15.00;desc=\"service.List\"")
	assert.Contains(t, header, "thread;dur=9.00;desc=\"querier.List\"")
	assert.NotContains(t, header, "service.Get")
	assert.True(t, strings.Index(header, "service.List") < strings.Index(header, "querier.List"))
}

func TestCollectorSkipsEmptyMetrics(t *testing.T) {
	t.Parallel()

	c := NewCollector()
	c.Observe("///", 10*time.Millisecond)

	assert.Equal(t, "", c.HeaderValue())
}

func TestCollectorCapsToTopMetrics(t *testing.T) {
	t.Parallel()

	c := NewCollector()
	for i := 0; i < 25; i++ {
		c.Observe("db.(*Query).Step"+strconv.Itoa(i), time.Duration(i+5)*time.Millisecond)
	}

	header := c.HeaderValue()
	parts := strings.Split(header, ", ")

	assert.Len(t, parts, maxMetrics)
	assert.NotContains(t, header, "desc=\"query.Step0\"")
	assert.NotContains(t, header, "desc=\"query.Step1\"")
	assert.NotContains(t, header, "desc=\"query.Step2\"")
	assert.NotContains(t, header, "desc=\"query.Step3\"")
	assert.NotContains(t, header, "desc=\"query.Step4\"")
	assert.Contains(t, header, "desc=\"query.Step24\"")
}
