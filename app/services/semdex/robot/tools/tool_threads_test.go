package tools

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/profile"
)

func TestLoadThreadDocumentReadsEveryReplyPage(t *testing.T) {
	threadID := post.ID(xid.New())
	reader := &pagedThreadDocumentReader{
		pages: map[int]*thread.Thread{
			1: threadPage(t, threadID, 3, "First reply"),
			2: threadPage(t, threadID, 3, "Second reply"),
			3: threadPage(t, threadID, 3, "Third reply"),
		},
	}

	opened, content, err := loadThreadDocument(context.Background(), reader, threadID)
	require.NoError(t, err)
	assert.Equal(t, "Paged thread", opened.Title)
	assert.Equal(t, []int{1, 2, 3}, reader.calls)
	assert.Equal(t, []int{reply.RepliesPerPage, reply.RepliesPerPage, reply.RepliesPerPage}, reader.sizes)
	assert.Contains(t, content.Plaintext(), "Original body")
	assert.Contains(t, content.Plaintext(), "Reply 1 by @member")
	assert.Contains(t, content.Plaintext(), "First reply")
	assert.Contains(t, content.Plaintext(), "Reply 2 by @member")
	assert.Contains(t, content.Plaintext(), "Second reply")
	assert.Contains(t, content.Plaintext(), "Reply 3 by @member")
	assert.Contains(t, content.Plaintext(), "Third reply")
}

type pagedThreadDocumentReader struct {
	pages map[int]*thread.Thread
	calls []int
	sizes []int
}

func (r *pagedThreadDocumentReader) Get(_ context.Context, _ post.ID, params pagination.Parameters) (*thread.Thread, error) {
	r.calls = append(r.calls, params.PageOneIndexed())
	r.sizes = append(r.sizes, params.Size())
	return r.pages[params.PageOneIndexed()], nil
}

func threadPage(t *testing.T, id post.ID, totalPages int, replyBody string) *thread.Thread {
	t.Helper()
	original, err := datagraph.NewRichText(`<p>Original body.</p>`)
	require.NoError(t, err)
	body, err := datagraph.NewRichText(`<p>` + replyBody + `.</p>`)
	require.NoError(t, err)

	return &thread.Thread{
		Post: post.Post{
			ID:      id,
			Content: original,
			Author:  profile.Ref{Handle: "author"},
		},
		Title: "Paged thread",
		Replies: pagination.Result[*reply.Reply]{
			TotalPages: totalPages,
			Items: []*reply.Reply{{
				Post: post.Post{
					ID:      post.ID(xid.New()),
					Content: body,
					Author:  profile.Ref{Handle: "member"},
				},
			}},
		},
	}
}
