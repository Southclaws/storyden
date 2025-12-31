package transform

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/cmd/import-mybb/loader"
	"github.com/Southclaws/storyden/cmd/import-mybb/writer"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/rs/xid"
)

func ImportCategories(ctx context.Context, w *writer.Writer, data *loader.MyBBData) error {
	if len(data.Forums) == 0 {
		log.Println("No forums to import")
		return nil
	}

	builders := make([]*ent.CategoryCreate, 0, len(data.Forums))

	for _, forum := range data.Forums {
		if forum.Active == 0 {
			log.Printf("Skipping inactive forum: %s", forum.Name)
			continue
		}

		id := xid.New()
		w.CategoryIDMap[forum.FID] = id

		parentID := parseParentFromList(forum.ParentList, w)

		builder := w.Client().Category.Create().
			SetID(id).
			SetName(forum.Name).
			SetSlug(mark.Slugify(forum.Name)).
			SetSort(forum.DispOrder)

		if forum.Description != "" {
			builder.SetDescription(forum.Description)
		}

		if !parentID.IsNil() {
			builder.SetParentCategoryID(parentID)
		}

		builders = append(builders, builder)
	}

	categories, err := w.CreateCategories(ctx, builders)
	if err != nil {
		return fmt.Errorf("create categories: %w", err)
	}

	log.Printf("Imported %d categories", len(categories))
	return nil
}

func parseParentFromList(parentList string, w *writer.Writer) xid.ID {
	if parentList == "" {
		return xid.NilID()
	}

	parts := strings.Split(parentList, ",")
	if len(parts) < 2 {
		return xid.NilID()
	}

	parentIDStr := strings.TrimSpace(parts[len(parts)-2])
	parentFID := 0
	fmt.Sscanf(parentIDStr, "%d", &parentFID)

	if parentID, ok := w.CategoryIDMap[parentFID]; ok {
		return parentID
	}

	return xid.NilID()
}
