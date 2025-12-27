package wordblock_checker

import (
	"context"
	"sync"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

func TestWordBlockChecker_Check(t *testing.T) {
	tests := []struct {
		name           string
		blockedWords   []string
		content        string
		expectBlocked  bool
		expectedReason string
	}{
		{
			name:          "no blocked words - should pass",
			blockedWords:  []string{},
			content:       "This is a normal post about programming",
			expectBlocked: false,
		},
		{
			name:           "contains exact blocked word - should block",
			blockedWords:   []string{"spam", "scam"},
			content:        "This is spam content",
			expectBlocked:  true,
			expectedReason: "Content contains blocked word: spam",
		},
		{
			name:           "contains blocked word case insensitive - should block",
			blockedWords:   []string{"VIAGRA"},
			content:        "Buy viagra now!",
			expectBlocked:  true,
			expectedReason: "Content contains blocked word: VIAGRA",
		},
		{
			name:           "contains blocked word with extra spaces - should block",
			blockedWords:   []string{"  BadWord  "},
			content:        "This has badword in it",
			expectBlocked:  true,
			expectedReason: "Content contains blocked word:   BadWord  ",
		},
		{
			name:          "blocked word not present - should pass",
			blockedWords:  []string{"spam", "scam"},
			content:       "This is a legitimate post",
			expectBlocked: false,
		},
		{
			name:           "partial word match - should block",
			blockedWords:   []string{"bad"},
			content:        "This is badly written",
			expectBlocked:  true,
			expectedReason: "Content contains blocked word: bad",
		},
		{
			name:          "empty content - should pass",
			blockedWords:  []string{"spam"},
			content:       "",
			expectBlocked: false,
		},
		{
			name:          "empty blocked list - should pass",
			blockedWords:  []string{},
			content:       "Any content should pass",
			expectBlocked: false,
		},
		{
			name:           "multiple blocked words, one matches - should block",
			blockedWords:   []string{"spam", "scam", "fraud"},
			content:        "This is a scam website",
			expectBlocked:  true,
			expectedReason: "Content contains blocked word: scam",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := &WordBlockChecker{
				blockedWords:  tt.blockedWords,
				normalizedMap: make(map[string]string),
				enabled:       true,
				mu:            sync.RWMutex{},
			}
			checker.buildNormalizedMap()

			content, err := datagraph.NewRichText(tt.content)
			require.NoError(t, err)

			result, err := checker.Check(
				context.Background(),
				xid.New(),
				datagraph.KindThread,
				"",
				content,
			)

			require.NoError(t, err)
			assert.Equal(t, tt.expectBlocked, result.RequiresReview)
			if tt.expectBlocked {
				assert.Contains(t, result.Reason, "Content contains blocked word:")
			}
		})
	}
}

func TestWordBlockChecker_Name(t *testing.T) {
	checker := &WordBlockChecker{}
	assert.Equal(t, "wordblock_checker", checker.Name())
}

func TestWordBlockChecker_Enabled(t *testing.T) {
	checker := &WordBlockChecker{enabled: true}
	assert.True(t, checker.Enabled())

	checker.enabled = false
	assert.False(t, checker.Enabled())
}

func TestWordBlockChecker_BuildNormalizedMap(t *testing.T) {
	checker := &WordBlockChecker{
		blockedWords: []string{
			"SPAM",
			"  ScAm  ",
			"",
			"  ",
			"fraud",
		},
		normalizedMap: make(map[string]string),
	}

	checker.buildNormalizedMap()

	assert.Equal(t, 3, len(checker.normalizedMap))
	assert.Contains(t, checker.normalizedMap, "spam")
	assert.Contains(t, checker.normalizedMap, "scam")
	assert.Contains(t, checker.normalizedMap, "fraud")
	assert.Equal(t, "SPAM", checker.normalizedMap["spam"])
	assert.Equal(t, "  ScAm  ", checker.normalizedMap["scam"])
}

func TestWordBlockChecker_ThreadSafety(t *testing.T) {
	checker := &WordBlockChecker{
		blockedWords:  []string{"spam"},
		normalizedMap: make(map[string]string),
		enabled:       true,
		mu:            sync.RWMutex{},
	}
	checker.buildNormalizedMap()

	errChan := make(chan error, 10)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			content, err := datagraph.NewRichText("This is spam")
			if err != nil {
				errChan <- err
				return
			}
			_, err = checker.Check(
				context.Background(),
				xid.New(),
				datagraph.KindThread,
				"",
				content,
			)
			if err != nil {
				errChan <- err
				return
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		t.Error(err)
	}
}
