package ai

import (
	"context"
	"math"
)

type Mock struct{}

func newMock() (*Mock, error) {
	return &Mock{}, nil
}

func (o *Mock) Prompt(ctx context.Context, input string) (*Result, error) {
	var output string
	if len(input) < 120 {
		output = input
	} else {
		output = input[:116] + "..."
	}
	return &Result{
		Answer: "An answer for " + output,
	}, nil
}

func (o *Mock) PromptStream(ctx context.Context, input string) (chan string, chan error) {
	return nil, nil
}

const mockEmbeddingSize = 3072

func (o *Mock) EmbeddingFunc() func(ctx context.Context, text string) ([]float32, error) {
	return func(ctx context.Context, text string) ([]float32, error) {
		embedding := make([]float32, mockEmbeddingSize)

		for _, v := range text {
			for i := range mockEmbeddingSize {
				c := ((float32(v % 256)) / 256) * float32(((i+1)*3071)%65535)
				embedding[i] = c
			}
		}

		return normalizeVector(embedding), nil
	}
}

func normalizeVector(vec []float32) []float32 {
	var n float32
	for _, v := range vec {
		n += v * v
	}
	n = float32(math.Sqrt(float64(n)))

	r := make([]float32, len(vec))
	for i, v := range vec {
		r[i] = v / n
	}

	return r
}
