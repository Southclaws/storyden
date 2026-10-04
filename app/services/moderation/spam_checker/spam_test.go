package spam_checker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

func TestDetector_Detect(t *testing.T) {
	for _, test := range []struct {
		file string
		spam bool
	}{
		{"post01_spam.txt", true},
		{"post02.txt", false},
		{"post03.txt", false},
		{"post04.txt", false},
		{"post05.txt", false},
		{"post06.txt", false},
		{"post07.txt", false},
		{"post08.txt", false},
		{"post09_spam.txt", true},
		{"post10_spam.txt", true},
	} {
		t.Run(test.file, func(t *testing.T) {
			content, err := os.ReadFile(test.file)
			require.NoError(t, err)

			spam, err := New().Detect(context.Background(), bytes.NewReader(content))
			require.NoError(t, err)
			require.Equal(t, test.spam, spam)
		})
	}
}

func TestDetector_GetRatioIncludesCompletedGZIPStream(t *testing.T) {
	d := repeatedContentDetector{}
	ratio, err := d.getRatio(strings.NewReader("a"))
	require.NoError(t, err)
	// A complete gzip stream has a 10-byte header, compressed data, and an 8-byte trailer.
	require.Greater(t, ratio, 18.0)
}

func TestDetector_EmptyInput(t *testing.T) {
	for _, input := range []string{"", " ", "\r\n \n"} {
		t.Run(input, func(t *testing.T) {
			d := repeatedContentDetector{}
			ratio, err := d.getRatio(strings.NewReader(input))
			require.NoError(t, err)
			require.Equal(t, 1.0, ratio)

			spam, err := New().Detect(context.Background(), strings.NewReader(input))
			require.NoError(t, err)
			require.False(t, spam)
		})
	}
}

func TestDetector_ReadError(t *testing.T) {
	wantErr := errors.New("read failed")
	for name, reader := range map[string]io.Reader{
		"before data": iotest.ErrReader(wantErr),
		"after data":  io.MultiReader(strings.NewReader(strings.Repeat("a", 4096)), iotest.ErrReader(wantErr)),
	} {
		t.Run(name, func(t *testing.T) {
			spam, err := New().Detect(context.Background(), reader)
			require.ErrorIs(t, err, wantErr)
			require.False(t, spam)
		})
	}
}
