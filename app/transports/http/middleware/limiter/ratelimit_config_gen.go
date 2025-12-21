// Package limiter contains rate limiting middleware.
// THIS FILE IS GENERATED. DO NOT EDIT MANUALLY.
// To edit rate limit configuration, edit the OpenAPI spec x-storyden extensions and run codegen.
package limiter

import "time"

// OperationRateLimitConfig defines per-operation rate limiting configuration
type OperationRateLimitConfig struct {
	// Cost is the number of requests this operation counts as
	Cost int
	// Limit is the maximum number of requests allowed in the period (0 means use global default)
	Limit int
	// Period is the time window for the limit (empty means use global default)
	Period time.Duration
}

// OperationRateLimits contains per-operation rate limit configurations extracted from OpenAPI spec
var OperationRateLimits = map[string]OperationRateLimitConfig{
	"AuthEmailSignup": {
		Cost:   10,
		Limit:  10,
		Period: time.Duration(int64(3600000000000)),
	},
	"AuthPasswordReset": {
		Cost:   5,
		Limit:  10,
		Period: time.Duration(int64(3600000000000)),
	},
	"AuthPasswordSignup": {
		Cost:   5,
		Limit:  20,
		Period: time.Duration(int64(3600000000000)),
	},
}

// GetOperationConfig returns the rate limit config for an operation, or nil if not configured
func GetOperationConfig(operationID string) *OperationRateLimitConfig {
	if cfg, ok := OperationRateLimits[operationID]; ok {
		return &cfg
	}
	return nil
}
