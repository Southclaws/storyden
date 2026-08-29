package robot_memory

import (
	"context"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot_memory "github.com/Southclaws/storyden/internal/ent/robotmemory"
)

const (
	ListLimit           = 25
	SearchLimit         = 20
	RootLimit           = 20
	searchExcerptRunes  = 240
	summaryExcerptRunes = 160
)

var (
	ErrNotFound       = fault.New("memory not found", ftag.With(ftag.NotFound))
	ErrInvalid        = fault.New("invalid memory input", ftag.With(ftag.InvalidArgument))
	ErrHierarchyCycle = fault.New("memory move would create a hierarchy cycle", ftag.With(ftag.InvalidArgument))
)

type Item struct {
	Memory   *robot.Memory
	Excerpt  string
	Children int
}

type Detail struct {
	Memory   *robot.Memory
	Path     []*Item
	Children []*Item
}

type SearchResult struct {
	Memory  *robot.Memory
	Path    []*Item
	Excerpt string
}

type SearchFilter struct {
	Query     string
	Subject   opt.Optional[string]
	Predicate opt.Optional[string]
	Object    opt.Optional[string]
}

type Option func(*ent.RobotMemoryMutation)

func WithContent(value string) Option {
	return func(m *ent.RobotMemoryMutation) {
		m.SetContent(strings.TrimSpace(value))
	}
}

func WithFact(subject, predicate, object string) Option {
	return func(m *ent.RobotMemoryMutation) {
		m.SetSubject(normaliseValue(subject, false))
		m.SetPredicate(normaliseValue(predicate, true))
		m.SetObject(normaliseValue(object, false))
	}
}

func WithoutFact() Option {
	return func(m *ent.RobotMemoryMutation) {
		m.ClearSubject()
		m.ClearPredicate()
		m.ClearObject()
	}
}

func WithState(value robot.MemoryState) Option {
	return func(m *ent.RobotMemoryMutation) {
		m.SetState(ent_robot_memory.State(value.String()))
	}
}

type Repository struct {
	db  *ent.Client
	raw *sqlx.DB
}

func New(db *ent.Client, raw *sqlx.DB) *Repository {
	return &Repository{db: db, raw: raw}
}

func (r *Repository) Create(
	ctx context.Context,
	robotRef string,
	parentID opt.Optional[robot.MemoryID],
	content string,
	opts ...Option,
) (*robot.Memory, error) {
	robotRef = strings.TrimSpace(robotRef)
	content = strings.TrimSpace(content)
	if robotRef == "" || content == "" {
		return nil, fault.Wrap(ErrInvalid, fctx.With(ctx))
	}
	if parentID, ok := parentID.Get(); ok {
		if _, err := r.getRow(ctx, robotRef, parentID); err != nil {
			return nil, err
		}
	}

	create := r.db.RobotMemory.Create().
		SetRobotRef(robotRef).
		SetContent(content)
	parentID.Call(func(id robot.MemoryID) { create.SetParentID(xid.ID(id)) })
	for _, option := range opts {
		option(create.Mutation())
	}
	if !validFactMutation(create.Mutation()) {
		return nil, fault.Wrap(ErrInvalid, fctx.With(ctx))
	}
	row, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return robot.MapMemory(row), nil
}

func (r *Repository) Open(ctx context.Context, robotRef string, id robot.MemoryID) (*Detail, error) {
	updated, err := r.db.RobotMemory.Update().
		Where(ent_robot_memory.IDEQ(xid.ID(id)), ent_robot_memory.RobotRefEQ(robotRef)).
		SetLastAccessedAt(time.Now().UTC()).
		AddAccessCount(1).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if updated == 0 {
		return nil, fault.Wrap(ErrNotFound, fctx.With(ctx))
	}
	return r.detail(ctx, robotRef, id)
}

func (r *Repository) List(ctx context.Context, robotRef string, parentID opt.Optional[robot.MemoryID], limit int) ([]*Item, bool, error) {
	if limit <= 0 || limit > ListLimit {
		limit = ListLimit
	}
	if parentID, ok := parentID.Get(); ok {
		if _, err := r.getRow(ctx, robotRef, parentID); err != nil {
			return nil, false, err
		}
	}

	query := r.db.RobotMemory.Query().Where(ent_robot_memory.RobotRefEQ(robotRef))
	if parentID, ok := parentID.Get(); ok {
		query.Where(ent_robot_memory.ParentIDEQ(xid.ID(parentID)))
	} else {
		query.Where(ent_robot_memory.ParentIDIsNil())
	}
	rows, err := query.
		Order(ent_robot_memory.ByUpdatedAt(sql.OrderDesc()), ent_robot_memory.ByID(sql.OrderAsc())).
		Limit(limit + 1).
		All(ctx)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	items, err := r.items(ctx, rows)
	return items, hasMore, err
}

func (r *Repository) ActiveRoot(ctx context.Context, robotRef string, limit int) ([]*Item, bool, error) {
	if limit <= 0 || limit > RootLimit {
		limit = RootLimit
	}
	rows, err := r.db.RobotMemory.Query().
		Where(
			ent_robot_memory.RobotRefEQ(robotRef),
			ent_robot_memory.ParentIDIsNil(),
			ent_robot_memory.StateEQ(ent_robot_memory.StateActive),
		).
		Order(ent_robot_memory.ByUpdatedAt(sql.OrderDesc()), ent_robot_memory.ByID(sql.OrderAsc())).
		Limit(limit + 1).
		All(ctx)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}
	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}
	items, err := r.items(ctx, rows)
	return items, hasMore, err
}

func (r *Repository) Update(ctx context.Context, robotRef string, id robot.MemoryID, opts ...Option) (*robot.Memory, error) {
	if len(opts) == 0 {
		return nil, fault.Wrap(ErrInvalid, fctx.With(ctx))
	}
	update := r.db.RobotMemory.Update().Where(
		ent_robot_memory.IDEQ(xid.ID(id)),
		ent_robot_memory.RobotRefEQ(robotRef),
	)
	for _, option := range opts {
		option(update.Mutation())
	}
	if content, exists := update.Mutation().Content(); exists && content == "" {
		return nil, fault.Wrap(ErrInvalid, fctx.With(ctx))
	}
	if !validFactMutation(update.Mutation()) {
		return nil, fault.Wrap(ErrInvalid, fctx.With(ctx))
	}
	count, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if count == 0 {
		return nil, fault.Wrap(ErrNotFound, fctx.With(ctx))
	}
	row, err := r.getRow(ctx, robotRef, id)
	if err != nil {
		return nil, err
	}
	return robot.MapMemory(row), nil
}

func validFactMutation(m *ent.RobotMemoryMutation) bool {
	subject, subjectSet := m.Subject()
	predicate, predicateSet := m.Predicate()
	object, objectSet := m.Object()
	if subjectSet || predicateSet || objectSet {
		return subjectSet && predicateSet && objectSet && subject != "" && predicate != "" && object != ""
	}
	subjectCleared := m.SubjectCleared()
	predicateCleared := m.PredicateCleared()
	objectCleared := m.ObjectCleared()
	return subjectCleared == predicateCleared && predicateCleared == objectCleared
}

func (r *Repository) Move(ctx context.Context, robotRef string, id robot.MemoryID, parentID opt.Optional[robot.MemoryID]) (*robot.Memory, error) {
	if _, err := r.getRow(ctx, robotRef, id); err != nil {
		return nil, err
	}
	if parentID, ok := parentID.Get(); ok {
		if parentID == id {
			return nil, fault.Wrap(ErrHierarchyCycle, fctx.With(ctx))
		}
		current := parentID
		seen := map[robot.MemoryID]struct{}{}
		for {
			if current == id {
				return nil, fault.Wrap(ErrHierarchyCycle, fctx.With(ctx))
			}
			if _, ok := seen[current]; ok {
				return nil, fault.Wrap(ErrHierarchyCycle, fctx.With(ctx))
			}
			seen[current] = struct{}{}
			row, err := r.getRow(ctx, robotRef, current)
			if err != nil {
				return nil, err
			}
			if row.ParentID == nil {
				break
			}
			current = robot.MemoryID(*row.ParentID)
		}
	}

	update := r.db.RobotMemory.Update().Where(
		ent_robot_memory.IDEQ(xid.ID(id)),
		ent_robot_memory.RobotRefEQ(robotRef),
	)
	if parentID, ok := parentID.Get(); ok {
		update.SetParentID(xid.ID(parentID))
	} else {
		update.ClearParentID()
	}
	count, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if count == 0 {
		return nil, fault.Wrap(ErrNotFound, fctx.With(ctx))
	}
	row, err := r.getRow(ctx, robotRef, id)
	if err != nil {
		return nil, err
	}
	return robot.MapMemory(row), nil
}

type searchRow struct {
	ID      string `db:"id"`
	Content string `db:"content"`
}

func (r *Repository) Search(
	ctx context.Context,
	robotRef string,
	filter SearchFilter,
	parentID opt.Optional[robot.MemoryID],
	limit int,
) ([]*SearchResult, bool, error) {
	built, err := buildSearchQuery(robotRef, filter, parentID, limit)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}
	if parentID, ok := parentID.Get(); ok {
		if _, err := r.getRow(ctx, robotRef, parentID); err != nil {
			return nil, false, err
		}
	}

	var rows []searchRow
	if err := r.raw.SelectContext(ctx, &rows, r.raw.Rebind(built.SQL), built.Args...); err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}
	hasMore := len(rows) > built.Limit
	if hasMore {
		rows = rows[:built.Limit]
	}
	results := make([]*SearchResult, 0, len(rows))
	ids := make([]xid.ID, 0, len(rows))
	for _, row := range rows {
		parsed, err := xid.FromString(row.ID)
		if err != nil {
			return nil, false, fault.Wrap(err, fctx.With(ctx))
		}
		id := robot.MemoryID(parsed)
		memoryRow, err := r.getRow(ctx, robotRef, id)
		if err != nil {
			return nil, false, err
		}
		path, err := r.path(ctx, robotRef, id)
		if err != nil {
			return nil, false, err
		}
		results = append(results, &SearchResult{
			Memory:  robot.MapMemory(memoryRow),
			Path:    path,
			Excerpt: makeExcerpt(row.Content, built.Terms, searchExcerptRunes),
		})
		ids = append(ids, parsed)
	}
	if len(ids) > 0 {
		if _, err := r.db.RobotMemory.Update().
			Where(ent_robot_memory.RobotRefEQ(robotRef), ent_robot_memory.IDIn(ids...)).
			SetLastAccessedAt(time.Now().UTC()).
			AddAccessCount(1).
			Save(ctx); err != nil {
			return nil, false, fault.Wrap(err, fctx.With(ctx))
		}
	}
	return results, hasMore, nil
}

func (r *Repository) detail(ctx context.Context, robotRef string, id robot.MemoryID) (*Detail, error) {
	row, err := r.getRow(ctx, robotRef, id)
	if err != nil {
		return nil, err
	}
	children, _, err := r.List(ctx, robotRef, opt.New(id), ListLimit)
	if err != nil {
		return nil, err
	}
	path, err := r.path(ctx, robotRef, id)
	if err != nil {
		return nil, err
	}
	return &Detail{Memory: robot.MapMemory(row), Path: path, Children: children}, nil
}

func (r *Repository) getRow(ctx context.Context, robotRef string, id robot.MemoryID) (*ent.RobotMemory, error) {
	row, err := r.db.RobotMemory.Query().Where(
		ent_robot_memory.IDEQ(xid.ID(id)),
		ent_robot_memory.RobotRefEQ(robotRef),
	).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(ErrNotFound, fctx.With(ctx))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return row, nil
}

func (r *Repository) path(ctx context.Context, robotRef string, id robot.MemoryID) ([]*Item, error) {
	rows := []*ent.RobotMemory{}
	seen := map[robot.MemoryID]struct{}{}
	current := id
	for {
		if _, ok := seen[current]; ok {
			return nil, fault.Wrap(ErrHierarchyCycle, fctx.With(ctx))
		}
		seen[current] = struct{}{}
		row, err := r.getRow(ctx, robotRef, current)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
		if row.ParentID == nil {
			break
		}
		current = robot.MemoryID(*row.ParentID)
	}
	slices.Reverse(rows)
	return r.items(ctx, rows)
}

func (r *Repository) items(ctx context.Context, rows []*ent.RobotMemory) ([]*Item, error) {
	result := make([]*Item, 0, len(rows))
	for _, row := range rows {
		children, err := row.QueryChildren().Count(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		memory := robot.MapMemory(row)
		result = append(result, &Item{
			Memory:   memory,
			Excerpt:  makeExcerpt(memory.Content, nil, summaryExcerptRunes),
			Children: children,
		})
	}
	return result, nil
}

func normaliseValue(value string, predicate bool) string {
	value = strings.ToLower(strings.Join(strings.Fields(value), " "))
	if !predicate {
		return value
	}
	var result strings.Builder
	separator := false
	for _, r := range value {
		if r == '-' || r == ' ' {
			separator = true
			continue
		}
		if separator && result.Len() > 0 {
			result.WriteByte('_')
		}
		separator = false
		result.WriteRune(r)
	}
	return strings.Trim(result.String(), "_")
}

func makeExcerpt(content string, terms []string, limit int) string {
	text := strings.Join(strings.Fields(content), " ")
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	lower := strings.ToLower(text)
	position := 0
	for _, term := range terms {
		if index := strings.Index(lower, term); index >= 0 {
			position = utf8.RuneCountInString(lower[:index])
			break
		}
	}
	start := max(0, position-limit/3)
	end := min(len(runes), start+limit)
	if end-start < limit {
		start = max(0, end-limit)
	}
	excerpt := string(runes[start:end])
	if start > 0 {
		excerpt = "…" + excerpt
	}
	if end < len(runes) {
		excerpt += "…"
	}
	return excerpt
}
