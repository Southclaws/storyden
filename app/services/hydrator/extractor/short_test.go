package extractor

import (
	"testing"
)

func Test_makeShortBody(t *testing.T) {
	type args struct {
		long string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"short",
			args{"this is a paragraph"},
			"this is a paragraph",
		},
		{
			"paras",
			args{"this is a paragraph\nsecond\nbthird\nfourth"},
			"this is a paragraph second bthird fourth",
		},
		{
			"long",
			args{"abcdefeghijklmnopqrstuvwxyzabcdefeghijklmnopqrstuvwxyzabcdefeghijklmnopqrstuvwxyzabcdefeghijklmnopqrstuvwxyzabcdefeghijklmnopqrstuvwxyzabcdefeghijklmnopqrstuvwxyzabcdefeghijklmnopqrstuvwxyzabcdefeghijklmnopqrstuvwxyz"},
			"abcdefeghijklmnopqrstuvwxyzabcdefeghijklmnopqrstuvwxyzabcdefeghi...",
		},
		{
			"headings",
			args{`# heading

Here's a paragraph. It's pretty neat.

Here's the rest of the text.

![http://image.com](http://image.com)

neat photo right?

This is quite a long post, the summary, should just be the first 128 characters rounded down to the nearest space.
`},
			"Here's a paragraph. It's pretty neat. Here's the rest of the text. neat photo right? This is quite a long post, the summary...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// got := MakeShortBody(tt.args.long)
			// assert.Equal(t, tt.want, got)
		})
	}
}
