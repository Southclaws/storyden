package node_writer

import (
	"context"
	"database/sql/driver"
	"fmt"
	"slices"
	"strings"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/tag/tag_writer"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/nodeversion"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

type CreateInput struct {
	ID      library.NodeID
	Name    string
	Slug    mark.Slug
	Options []Option
	Tags    opt.Optional[tag_ref.Names]
}

type UpdateInput struct {
	ID        library.NodeID
	UpdatedAt time.Time
	Versioned bool
	Options   []Option
	Tags      opt.Optional[tag_ref.Names]
}

func (w *Writer) CreateMany(ctx context.Context, owner account.AccountID, items []CreateInput) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := w.db.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	builders := make([]*ent.NodeCreate, len(items))
	tags := make(map[library.NodeID]tag_ref.Names)
	for i, item := range items {
		builder := tx.Node.Create().SetID(xid.ID(item.ID)).SetOwnerID(xid.ID(owner)).SetName(item.Name).SetSlug(item.Slug.String())
		for _, option := range item.Options {
			option(builder.Mutation())
		}

		builders[i] = builder
		if names, ok := item.Tags.Get(); ok {
			tags[item.ID] = names
		}
	}

	if err := setBatchSortKeys(ctx, tx.Client(), builders); err != nil {
		return err
	}

	// A single INSERT lets new nodes refer to other IDs in the same hierarchy.
	if err := tx.Node.CreateBulk(builders...).Exec(ctx); err != nil {
		return err
	}

	if err := w.replaceBatchTags(ctx, tx.Client(), tags); err != nil {
		return err
	}

	return tx.Commit()
}

func (w *Writer) UpdateMany(ctx context.Context, items []UpdateInput) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := w.db.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	predicates := make([]predicate.Node, len(items))
	mutations := make(map[xid.ID]*ent.NodeMutation, len(items))
	tags := make(map[library.NodeID]tag_ref.Names)
	for i, item := range items {
		id := xid.ID(item.ID)
		condition := node.And(node.ID(id), node.UpdatedAtEQ(item.UpdatedAt))
		if item.Versioned {
			condition = node.And(condition, node.Not(node.HasVersionsWith(nodeversion.StatusEQ(nodeversion.StatusDraft))))
		}

		predicates[i] = condition
		mutation := tx.Node.Update().Mutation()
		for _, option := range item.Options {
			option(mutation)
		}

		mutations[id] = mutation
		if names, ok := item.Tags.Get(); ok {
			tags[item.ID] = names
		}
	}

	changed, err := tx.Node.Update().Where(node.Or(predicates...)).Modify(batchUpdate(mutations)).Save(ctx)
	if err != nil {
		return err
	}

	if changed != len(items) {
		return fmt.Errorf("pages changed while preparing the batch; no pages were updated")
	}

	if err := w.replaceBatchTags(ctx, tx.Client(), tags); err != nil {
		return err
	}

	return tx.Commit()
}

func batchUpdate(mutations map[xid.ID]*ent.NodeMutation) func(*sql.UpdateBuilder) {
	return func(update *sql.UpdateBuilder) {
		columns := make(map[string]map[xid.ID]any)
		for id, mutation := range mutations {
			for _, field := range mutation.Fields() {
				value, _ := mutation.Field(field)
				if columns[field] == nil {
					columns[field] = make(map[xid.ID]any)
				}

				columns[field][id] = value
			}

			for _, field := range mutation.ClearedFields() {
				if columns[field] == nil {
					columns[field] = make(map[xid.ID]any)
				}

				columns[field][id] = nil
			}
		}

		fields := make([]string, 0, len(columns))
		for field := range columns {
			fields = append(fields, field)
		}
		slices.Sort(fields)

		for _, field := range fields {
			values := columns[field]
			ids := make([]xid.ID, 0, len(values))
			for id := range values {
				ids = append(ids, id)
			}
			slices.SortFunc(ids, func(a, b xid.ID) int { return a.Compare(b) })

			update.Set(field, sql.ExprFunc(func(b *sql.Builder) {
				b.WriteString("CASE ").Ident(node.FieldID)
				for _, id := range ids {
					b.WriteString(" WHEN ").Arg(id).WriteString(" THEN ").Arg(values[id])
				}

				b.WriteString(" ELSE ").Ident(field).WriteString(" END")
			}))
		}
	}
}

func (w *Writer) replaceBatchTags(ctx context.Context, db *ent.Client, replacements map[library.NodeID]tag_ref.Names) error {
	if len(replacements) == 0 {
		return nil
	}

	names := tag_ref.Names{}
	ids := make([]driver.Value, 0, len(replacements))
	for id, tags := range replacements {
		ids = append(ids, xid.ID(id))
		names = append(names, tags...)
	}
	slices.SortFunc(names, func(a, b tag_ref.Name) int { return strings.Compare(a.String(), b.String()) })
	names = slices.Compact(names)

	tags, err := tag_writer.New(db).Add(ctx, names...)
	if err != nil {
		return err
	}

	byName := make(map[tag_ref.Name]xid.ID, len(tags))
	for _, tag := range tags {
		byName[tag.Name] = xid.ID(tag.ID)
	}

	builder := sql.Dialect(dialect.SQLite)
	if w.raw.DriverName() == "pgx" {
		builder = sql.Dialect(dialect.Postgres)
	}

	query, args := builder.Delete(node.TagsTable).Where(sql.InValues(node.TagsPrimaryKey[1], ids...)).Query()
	if _, err := db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	insert := builder.Insert(node.TagsTable).Columns(node.TagsPrimaryKey...)
	count := 0
	for id, names := range replacements {
		seen := make(map[tag_ref.Name]bool, len(names))
		for _, name := range names {
			if seen[name] {
				continue
			}

			insert.Values(byName[name], xid.ID(id))
			seen[name] = true
			count++
		}
	}

	if count == 0 {
		return nil
	}

	query, args = insert.Query()
	_, err = db.ExecContext(ctx, query, args...)
	return err
}
