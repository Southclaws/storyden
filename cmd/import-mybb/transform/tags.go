package transform

import (
	"context"
	"fmt"
	"log"

	"github.com/Southclaws/storyden/cmd/import-mybb/loader"
	"github.com/Southclaws/storyden/cmd/import-mybb/logger"
	"github.com/Southclaws/storyden/cmd/import-mybb/writer"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/rs/xid"
)

func ImportTags(ctx context.Context, w *writer.Writer, data *loader.MyBBData) error {
	if len(data.ThreadPrefixes) == 0 {
		log.Println("No thread prefixes to import")
		return nil
	}

	builders := make([]*ent.TagCreate, 0, len(data.ThreadPrefixes))

	for _, prefix := range data.ThreadPrefixes {
		if prefix.Prefix == "" {
			continue
		}

		id := xid.New()
		w.TagIDMap[prefix.PID] = id

		builder := w.Client().Tag.Create().
			SetID(id).
			SetName(prefix.Prefix)

		builders = append(builders, builder)

		logger.Tag(prefix.PID, prefix.Prefix)
	}

	tags, err := w.CreateTags(ctx, builders)
	if err != nil {
		return fmt.Errorf("create tags: %w", err)
	}

	log.Printf("Imported %d tags", len(tags))
	return nil
}
