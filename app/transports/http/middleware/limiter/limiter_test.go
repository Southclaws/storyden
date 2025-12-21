package limiter_test

import (
	"testing"
	"time"

	"github.com/Southclaws/storyden/app/transports/http/middleware/limiter"
)

func TestOperationRateLimits(t *testing.T) {
	// Test that specific operations have rate limit configs
	tests := []struct {
		operationID string
		wantCost    int
		wantLimit   int
		wantPeriod  time.Duration
	}{
		{
			operationID: "AuthEmailSignup",
			wantCost:    10,
			wantLimit:   10,
			wantPeriod:  time.Hour,
		},
		{
			operationID: "AuthPasswordSignup",
			wantCost:    5,
			wantLimit:   20,
			wantPeriod:  time.Hour,
		},
		{
			operationID: "AuthPasswordReset",
			wantCost:    5,
			wantLimit:   10,
			wantPeriod:  time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.operationID, func(t *testing.T) {
			cfg := limiter.GetOperationConfig(tt.operationID)
			if cfg == nil {
				t.Fatalf("expected config for %s, got nil", tt.operationID)
			}

			if cfg.Cost != tt.wantCost {
				t.Errorf("cost: got %d, want %d", cfg.Cost, tt.wantCost)
			}

			if cfg.Limit != tt.wantLimit {
				t.Errorf("limit: got %d, want %d", cfg.Limit, tt.wantLimit)
			}

			if cfg.Period != tt.wantPeriod {
				t.Errorf("period: got %v, want %v", cfg.Period, tt.wantPeriod)
			}
		})
	}
}

func TestGetOperationConfig_NotConfigured(t *testing.T) {
	// Test that operations without x-storyden config return nil
	cfg := limiter.GetOperationConfig("SomeRandomOperation")
	if cfg != nil {
		t.Errorf("expected nil for unconfigured operation, got %+v", cfg)
	}
}

func TestRouteToOperation(t *testing.T) {
	// Test that routes map to correct operation IDs
	tests := []struct {
		method      string
		path        string
		wantOpID    string
	}{
		{
			method:   "post",
			path:     "/api/auth/email/signup",
			wantOpID: "AuthEmailSignup",
		},
		{
			method:   "post",
			path:     "/api/auth/password/signup",
			wantOpID: "AuthPasswordSignup",
		},
		{
			method:   "post",
			path:     "/api/auth/password/reset",
			wantOpID: "AuthPasswordReset",
		},
		{
			method:   "get",
			path:     "/api/threads",
			wantOpID: "ThreadList",
		},
	}

	for _, tt := range tests {
		t.Run(tt.method+"_"+tt.path, func(t *testing.T) {
			opID := limiter.GetOperationIDFromRoute(tt.method, tt.path)
			if opID != tt.wantOpID {
				t.Errorf("got operation ID %q, want %q", opID, tt.wantOpID)
			}
		})
	}
}

func TestRouteToOperation_NotFound(t *testing.T) {
	// Test that unknown routes return empty string
	opID := limiter.GetOperationIDFromRoute("GET", "/api/nonexistent")
	if opID != "" {
		t.Errorf("expected empty string for unknown route, got %q", opID)
	}
}
