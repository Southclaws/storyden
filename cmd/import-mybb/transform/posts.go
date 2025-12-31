package transform

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/cmd/import-mybb/loader"
	"github.com/Southclaws/storyden/cmd/import-mybb/writer"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/rs/xid"
)

func ImportThreads(ctx context.Context, w *writer.Writer, data *loader.MyBBData) error {
	if len(data.Threads) == 0 {
		log.Println("No threads to import")
		return nil
	}

	// Create a map of PID -> MyBBPost for quick lookup
	postMap := make(map[int]loader.MyBBPost)
	for _, p := range data.Posts {
		postMap[p.PID] = p
	}

	// Track which posts are used as first posts (root posts)
	firstPostPIDs := make(map[int]bool)

	builders := make([]*ent.PostCreate, 0, len(data.Threads))

	for _, thread := range data.Threads {
		// Find the first post for this thread
		firstPost, ok := postMap[thread.FirstPost]
		if !ok {
			log.Printf("Skipping thread %d: first post PID %d not found", thread.TID, thread.FirstPost)
			continue
		}

		// Mark this post as a first post so we don't import it as a reply later
		firstPostPIDs[thread.FirstPost] = true

		accountID, ok := w.AccountIDMap[firstPost.UID]
		if !ok {
			log.Printf("Skipping thread %d: author UID %d not found", thread.TID, firstPost.UID)
			continue
		}

		categoryID, ok := w.CategoryIDMap[thread.FID]
		if !ok {
			log.Printf("Skipping thread %d: category FID %d not found", thread.TID, thread.FID)
			continue
		}

		rootPostID := xid.New()
		w.PostIDMap[thread.TID] = rootPostID
		w.PostIDMap[thread.FirstPost] = rootPostID // Map the first post PID too

		createdAt := time.Unix(firstPost.DateLine, 0)
		lastReplyAt := time.Unix(thread.LastPost, 0)
		if lastReplyAt.IsZero() {
			lastReplyAt = createdAt
		}

		visibility := mapVisibility(thread.Visible)

		// Format slug as "<id>-<slug>" for thread_mark service compatibility
		threadSlug := fmt.Sprintf("%s-%s", rootPostID.String(), mark.Slugify(thread.Subject))

		// Use the first post's body as the thread body, converting BBCode to HTML
		content := convertBBCodeToHTML(firstPost.Message)

		builder := w.Client().Post.Create().
			SetID(rootPostID).
			SetTitle(thread.Subject).
			SetSlug(threadSlug).
			SetBody(content.HTML()).
			SetShort(content.Short()).
			SetAccountPosts(accountID).
			SetCategoryID(categoryID).
			SetPinned(thread.Sticky == 1).
			SetVisibility(visibility).
			SetCreatedAt(createdAt).
			SetUpdatedAt(lastReplyAt).
			SetLastReplyAt(lastReplyAt)

		if thread.Prefix > 0 {
			if tagID, ok := w.TagIDMap[thread.Prefix]; ok {
				builder.AddTagIDs(tagID)
			}
		}

		if thread.DeleteTime > 0 {
			deletedAt := time.Unix(thread.DeleteTime, 0)
			builder.SetDeletedAt(deletedAt)
		}

		builders = append(builders, builder)
	}

	// Store firstPostPIDs in writer for use in ImportPosts
	w.FirstPostPIDs = firstPostPIDs

	posts, err := w.CreatePosts(ctx, builders)
	if err != nil {
		return fmt.Errorf("create thread posts: %w", err)
	}

	log.Printf("Imported %d thread posts", len(posts))
	return nil
}

func ImportPosts(ctx context.Context, w *writer.Writer, data *loader.MyBBData) error {
	if len(data.Posts) == 0 {
		log.Println("No posts to import")
		return nil
	}

	builders := make([]*ent.PostCreate, 0, len(data.Posts))
	skippedFirstPosts := 0

	for _, p := range data.Posts {
		// Skip posts that are first posts (already imported as root posts)
		if w.FirstPostPIDs[p.PID] {
			skippedFirstPosts++
			continue
		}

		accountID, ok := w.AccountIDMap[p.UID]
		if !ok {
			log.Printf("Skipping post %d: author UID %d not found", p.PID, p.UID)
			continue
		}

		rootPostID, ok := w.PostIDMap[p.TID]
		if !ok {
			log.Printf("Skipping post %d: thread TID %d not found", p.PID, p.TID)
			continue
		}

		replyID := xid.New()
		w.PostIDMap[p.PID] = replyID

		createdAt := time.Unix(p.DateLine, 0)

		visibility := mapVisibility(p.Visible)

		// Convert BBCode to HTML
		content := convertBBCodeToHTML(p.Message)

		builder := w.Client().Post.Create().
			SetID(replyID).
			SetBody(content.HTML()).
			SetShort(content.Short()).
			SetAccountPosts(accountID).
			SetVisibility(visibility).
			SetCreatedAt(createdAt).
			SetUpdatedAt(createdAt).
			SetLastReplyAt(createdAt)

		if rootPostID != xid.NilID() {
			builder.SetRootPostID(rootPostID)
		}

		// MyBB doesn't support nested replies, so skip ReplyToPostID

		builders = append(builders, builder)
	}

	posts, err := w.CreatePosts(ctx, builders)
	if err != nil {
		return fmt.Errorf("create reply posts: %w", err)
	}

	log.Printf("Imported %d reply posts (skipped %d first posts)", len(posts), skippedFirstPosts)
	return nil
}

func mapVisibility(visible int) post.Visibility {
	switch visible {
	case 1:
		return post.VisibilityPublished
	case 0:
		return post.VisibilityDraft
	case -1:
		return post.VisibilityReview // Unapproved/awaiting moderation
	case -2:
		return post.VisibilityDraft // Soft deleted in MyBB
	default:
		return post.VisibilityPublished
	}
}
