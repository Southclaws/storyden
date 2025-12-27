package checker

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/rs/xid"
)

// Checker is an interface for content moderation checks.
// Each implementation can perform a specific type of content validation.
type Checker interface {
	// Check examines the content and returns a result indicating whether
	// the content should be flagged for review.
	Check(ctx context.Context, targetID xid.ID, targetKind datagraph.Kind, name string, content datagraph.Content) (*Result, error)

	// Name returns a unique identifier for this checker.
	Name() string

	// Enabled returns whether this checker is currently active.
	Enabled() bool
}

type Result struct {
	RequiresReview bool
	Reason         string
}

// Registry holds all registered content moderation checkers.
type Registry struct {
	checkers []Checker
}

// NewRegistry creates a new checker registry.
func NewRegistry(checkers ...Checker) *Registry {
	return &Registry{
		checkers: checkers,
	}
}

// Add registers a new checker.
func (r *Registry) Add(checker Checker) {
	r.checkers = append(r.checkers, checker)
}

// GetEnabled returns all enabled checkers.
func (r *Registry) GetEnabled() []Checker {
	enabled := make([]Checker, 0, len(r.checkers))
	for _, c := range r.checkers {
		if c.Enabled() {
			enabled = append(enabled, c)
		}
	}
	return enabled
}
