package searchscore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScoreAppliesExactSubstringAndTermWeights(t *testing.T) {
	query := New("  Plugin Builder  ")

	score := query.Score(
		Field{Text: "Plugin Builder", ExactMatch: 100, SubstringMatch: 40, TermMatch: 8},
		Field{Text: "Build a Plugin Builder workflow", SubstringMatch: 20, TermMatch: 4},
		Field{Text: "Use the Plugin Builder tools", SubstringMatch: 5},
	)

	assert.Equal(t, 149, score)
}

func TestScorePreservesPerFieldWeights(t *testing.T) {
	query := New("discussion")

	score := query.Score(
		Field{Text: "thread_get", ExactMatch: 120, SubstringMatch: 50, TermMatch: 10},
		Field{Text: "Get discussion", ExactMatch: 100, SubstringMatch: 40, TermMatch: 8},
		Field{Text: "Read one community discussion by ID", SubstringMatch: 20, TermMatch: 4},
	)

	assert.Equal(t, 72, score)
}

func TestScoreReturnsZeroForEmptyQuery(t *testing.T) {
	score := New("   ").Score(
		Field{Text: "anything", SubstringMatch: 100},
	)

	assert.Equal(t, 0, score)
}
