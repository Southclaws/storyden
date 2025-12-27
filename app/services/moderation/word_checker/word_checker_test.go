package word_checker

import (
	"context"
	"sync"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/moderation/checker"
)

func TestWordChecker_BlockList(t *testing.T) {
	tests := []struct {
		name           string
		blockList      []string
		content        string
		expectedAction checker.Action
		expectedReason string
	}{
		{
			name:           "no blocked words - should allow",
			blockList:      []string{},
			content:        "This is a normal post about programming",
			expectedAction: checker.ActionAllow,
		},
		{
			name:           "contains exact blocked word - should reject",
			blockList:      []string{"spam", "scam"},
			content:        "This is spam content",
			expectedAction: checker.ActionReject,
			expectedReason: "Content violates community guidelines",
		},
		{
			name:           "contains blocked word case insensitive - should reject",
			blockList:      []string{"VIAGRA"},
			content:        "Buy viagra now!",
			expectedAction: checker.ActionReject,
			expectedReason: "Content violates community guidelines",
		},
		{
			name:           "contains blocked word with extra spaces - should reject",
			blockList:      []string{"  BadWord  "},
			content:        "This has badword in it",
			expectedAction: checker.ActionReject,
			expectedReason: "Content violates community guidelines",
		},
		{
			name:           "blocked word not present - should allow",
			blockList:      []string{"spam", "scam"},
			content:        "This is a legitimate post",
			expectedAction: checker.ActionAllow,
		},
		{
			name:           "partial word match - should reject",
			blockList:      []string{"bad"},
			content:        "This is badly written",
			expectedAction: checker.ActionReject,
			expectedReason: "Content violates community guidelines",
		},
		{
			name:           "empty content - should allow",
			blockList:      []string{"spam"},
			content:        "",
			expectedAction: checker.ActionAllow,
		},
		{
			name:           "multiple blocked words, one matches - should reject",
			blockList:      []string{"spam", "scam", "fraud"},
			content:        "This is a scam website",
			expectedAction: checker.ActionReject,
			expectedReason: "Content violates community guidelines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WordChecker{
				blockListNormalized:  normalizeWordList(tt.blockList),
				reportListNormalized: make(map[string]string),
				enabled:              true,
				mu:                   sync.RWMutex{},
			}

			content, err := datagraph.NewRichText(tt.content)
			require.NoError(t, err)

			result, err := wc.Check(
				context.Background(),
				xid.New(),
				datagraph.KindThread,
				"",
				content,
			)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedAction, result.Action)
			if tt.expectedAction == checker.ActionReject {
				assert.Equal(t, tt.expectedReason, result.Reason)
			}
		})
	}
}

func TestWordChecker_ReportList(t *testing.T) {
	tests := []struct {
		name           string
		reportList     []string
		content        string
		expectedAction checker.Action
		expectedReason string
	}{
		{
			name:           "no report words - should allow",
			reportList:     []string{},
			content:        "This is a normal post",
			expectedAction: checker.ActionAllow,
		},
		{
			name:           "contains report word - should report",
			reportList:     []string{"offensive"},
			content:        "This is offensive content",
			expectedAction: checker.ActionReport,
			expectedReason: "Content contains blocked word: offensive",
		},
		{
			name:           "contains report word case insensitive - should report",
			reportList:     []string{"SLUR"},
			content:        "This has slur in it",
			expectedAction: checker.ActionReport,
			expectedReason: "Content contains blocked word: SLUR",
		},
		{
			name:           "report word not present - should allow",
			reportList:     []string{"offensive"},
			content:        "This is a legitimate post",
			expectedAction: checker.ActionAllow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &WordChecker{
				blockListNormalized:  make(map[string]string),
				reportListNormalized: normalizeWordList(tt.reportList),
				enabled:              true,
				mu:                   sync.RWMutex{},
			}

			content, err := datagraph.NewRichText(tt.content)
			require.NoError(t, err)

			result, err := wc.Check(
				context.Background(),
				xid.New(),
				datagraph.KindThread,
				"",
				content,
			)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedAction, result.Action)
			if tt.expectedAction == checker.ActionReport {
				assert.Equal(t, tt.expectedReason, result.Reason)
			}
		})
	}
}

func TestWordChecker_BlockListTakesPrecedence(t *testing.T) {
	wc := &WordChecker{
		blockListNormalized:  normalizeWordList([]string{"bad"}),
		reportListNormalized: normalizeWordList([]string{"bad"}),
		enabled:              true,
		mu:                   sync.RWMutex{},
	}

	content, err := datagraph.NewRichText("This has bad content")
	require.NoError(t, err)

	result, err := wc.Check(
		context.Background(),
		xid.New(),
		datagraph.KindThread,
		"",
		content,
	)

	require.NoError(t, err)
	assert.Equal(t, checker.ActionReject, result.Action)
	assert.Equal(t, "Content violates community guidelines", result.Reason)
}

func TestWordChecker_Name(t *testing.T) {
	wc := &WordChecker{}
	assert.Equal(t, "word_checker", wc.Name())
}

func TestWordChecker_Enabled(t *testing.T) {
	wc := &WordChecker{enabled: true}
	assert.True(t, wc.Enabled())

	wc.enabled = false
	assert.False(t, wc.Enabled())
}

func TestNormalizeWordList(t *testing.T) {
	normalized := normalizeWordList([]string{
		"SPAM",
		"  ScAm  ",
		"",
		"  ",
		"fraud",
	})

	assert.Equal(t, 3, len(normalized))
	assert.Contains(t, normalized, "spam")
	assert.Contains(t, normalized, "scam")
	assert.Contains(t, normalized, "fraud")
	assert.Equal(t, "SPAM", normalized["spam"])
	assert.Equal(t, "  ScAm  ", normalized["scam"])
}

func TestWordChecker_ThreadSafety(t *testing.T) {
	wc := &WordChecker{
		blockListNormalized:  normalizeWordList([]string{"spam"}),
		reportListNormalized: make(map[string]string),
		enabled:              true,
		mu:                   sync.RWMutex{},
	}

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
			_, err = wc.Check(
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
