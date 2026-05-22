package transform

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/cmd/import-mybb/loader"
	"github.com/Southclaws/storyden/cmd/import-mybb/writer"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/rs/xid"
)

func ImportInteractions(ctx context.Context, w *writer.Writer, data *loader.MyBBData) error {
	if err := importReacts(ctx, w, data); err != nil {
		return err
	}

	if err := importLikes(ctx, w, data); err != nil {
		return err
	}

	if err := importReads(ctx, w, data); err != nil {
		return err
	}

	if err := importReports(ctx, w, data); err != nil {
		return err
	}

	return nil
}

func importReacts(ctx context.Context, w *writer.Writer, data *loader.MyBBData) error {
	if len(data.Reputation) == 0 {
		log.Println("No reputation to import")
		return nil
	}

	builders := make([]*ent.ReactCreate, 0)
	reactSet := make(map[string]bool)

	for _, rep := range data.Reputation {
		accountID, ok := w.AccountIDMap[rep.AddUID]
		if !ok {
			continue
		}

		postID, ok := w.PostIDMap[rep.PID]
		if !ok {
			continue
		}

		emoji := "👍"
		reactKey := fmt.Sprintf("%s:%s:%s", accountID.String(), postID.String(), emoji)
		if reactSet[reactKey] {
			continue
		}
		reactSet[reactKey] = true

		createdAt := time.Unix(rep.DateAdded, 0)

		builder := w.Client().React.Create().
			SetID(xid.New()).
			SetEmoji(emoji).
			SetAccountID(accountID).
			SetPostID(postID).
			SetCreatedAt(createdAt)

		builders = append(builders, builder)
	}

	if len(builders) == 0 {
		log.Println("No reacts to import after filtering")
		return nil
	}

	reacts, err := w.CreateReacts(ctx, builders)
	if err != nil {
		return fmt.Errorf("create reacts: %w", err)
	}

	log.Printf("Imported %d reacts", len(reacts))
	return nil
}

func importLikes(ctx context.Context, w *writer.Writer, data *loader.MyBBData) error {
	if len(data.ThreadRatings) == 0 {
		log.Println("No thread ratings to import")
		return nil
	}

	builders := make([]*ent.LikePostCreate, 0)
	likeSet := make(map[string]bool)

	for _, rating := range data.ThreadRatings {
		accountID, ok := w.AccountIDMap[rating.UID]
		if !ok {
			continue
		}

		postID, ok := w.PostIDMap[rating.TID]
		if !ok {
			continue
		}

		likeKey := fmt.Sprintf("%s:%s", accountID.String(), postID.String())
		if likeSet[likeKey] {
			continue
		}
		likeSet[likeKey] = true

		builder := w.Client().LikePost.Create().
			SetID(xid.New()).
			SetAccountID(accountID).
			SetPostID(postID)

		builders = append(builders, builder)
	}

	if len(builders) == 0 {
		log.Println("No likes to import after filtering")
		return nil
	}

	likes, err := w.CreateLikePosts(ctx, builders)
	if err != nil {
		return fmt.Errorf("create likes: %w", err)
	}

	log.Printf("Imported %d likes", len(likes))
	return nil
}

func importReads(ctx context.Context, w *writer.Writer, data *loader.MyBBData) error {
	if len(data.ThreadsRead) == 0 {
		log.Println("No thread reads to import")
		return nil
	}

	builders := make([]*ent.PostReadCreate, 0)
	readSet := make(map[string]bool)

	for _, read := range data.ThreadsRead {
		accountID, ok := w.AccountIDMap[read.UID]
		if !ok {
			continue
		}

		postID, ok := w.PostIDMap[read.TID]
		if !ok {
			continue
		}

		readKey := fmt.Sprintf("%s:%s", postID.String(), accountID.String())
		if readSet[readKey] {
			continue
		}
		readSet[readKey] = true

		createdAt := time.Unix(read.DateLine, 0)

		builder := w.Client().PostRead.Create().
			SetID(xid.New()).
			SetAccountID(accountID).
			SetRootPostID(postID).
			SetLastSeenAt(createdAt)

		builders = append(builders, builder)
	}

	if len(builders) == 0 {
		log.Println("No reads to import after filtering")
		return nil
	}

	reads, err := w.CreatePostReads(ctx, builders)
	if err != nil {
		return fmt.Errorf("create reads: %w", err)
	}

	log.Printf("Imported %d reads", len(reads))
	return nil
}

func importReports(ctx context.Context, w *writer.Writer, data *loader.MyBBData) error {
	if len(data.ReportedContent) == 0 {
		log.Println("No reported content to import")
		return nil
	}

	builders := make([]*ent.ReportCreate, 0)

	for _, report := range data.ReportedContent {
		if report.Type != "profile" {
			continue
		}

		targetAccountID, ok := w.AccountIDMap[report.ID]
		if !ok {
			continue
		}

		accountID, ok := w.AccountIDMap[report.UID]
		if !ok {
			continue
		}

		createdAt := time.Unix(report.DateLine, 0)

		builder := w.Client().Report.Create().
			SetID(xid.New()).
			SetReason(report.Reason).
			SetReportedByID(accountID).
			SetTargetID(targetAccountID).
			SetTargetKind(datagraph.KindProfile.String()).
			SetCreatedAt(createdAt)

		builders = append(builders, builder)
	}

	if len(builders) == 0 {
		log.Println("No reports to import after filtering")
		return nil
	}

	reports, err := w.CreateReports(ctx, builders)
	if err != nil {
		return fmt.Errorf("create reports: %w", err)
	}

	log.Printf("Imported %d reports", len(reports))
	return nil
}

func ImportAssets(ctx context.Context, w *writer.Writer, data *loader.MyBBData) error {
	log.Println("Asset import not yet implemented - requires filesystem access")
	return nil
}
