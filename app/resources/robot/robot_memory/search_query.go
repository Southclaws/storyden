package robot_memory

import (
	"fmt"
	"strings"

	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/robot"
)

const (
	searchRobotRefMaxBytes        = 256
	searchQueryMaxBytes           = 512
	searchTermLimit               = 16
	searchStructuredValueMaxBytes = 256
)

const searchSubtreeCTE = `with recursive subtree (id) as (
    select
        id
    from
        robot_memories
    where
        robot_ref = ?
        and id = ?
union
    select
        memory.id
    from
        robot_memories memory
        inner join subtree parent on memory.parent_id = parent.id
    where
        memory.robot_ref = ?
)
`

const searchStatement = `select
    id,
    content
from
    robot_memories
where
    %s
order by
    updated_at desc,
    id asc
limit ?
`

type builtSearchQuery struct {
	SQL   string
	Args  []any
	Terms []string
	Limit int
}

type structuredSearchColumn uint8

const (
	structuredSearchSubject structuredSearchColumn = iota
	structuredSearchPredicate
	structuredSearchObject
)

func buildSearchQuery(
	robotRef string,
	filter SearchFilter,
	parentID opt.Optional[robot.MemoryID],
	limit int,
) (builtSearchQuery, error) {
	if robotRef == "" || len(robotRef) > searchRobotRefMaxBytes {
		return builtSearchQuery{}, fmt.Errorf("%w: invalid Robot reference", ErrInvalid)
	}
	if len(filter.Query) > searchQueryMaxBytes {
		return builtSearchQuery{}, fmt.Errorf("%w: query exceeds %d bytes", ErrInvalid, searchQueryMaxBytes)
	}

	terms := strings.Fields(strings.ToLower(strings.TrimSpace(filter.Query)))
	if len(terms) > searchTermLimit {
		return builtSearchQuery{}, fmt.Errorf("%w: query exceeds %d terms", ErrInvalid, searchTermLimit)
	}
	if limit <= 0 || limit > SearchLimit {
		limit = SearchLimit
	}

	predicates := []string{
		"robot_ref = ?",
		"state = 'active'",
	}
	args := make([]any, 0, 8+len(terms))
	cte := ""
	if parentID, ok := parentID.Get(); ok {
		cte = searchSubtreeCTE
		args = append(args, robotRef, parentID.String(), robotRef)
		predicates = append(predicates, "id in (select id from subtree)")
	}
	args = append(args, robotRef)

	for _, term := range terms {
		predicates = append(predicates, "lower(content) like ? escape '\\'")
		args = append(args, "%"+escapeLike(term)+"%")
	}

	structured := []struct {
		column    structuredSearchColumn
		value     opt.Optional[string]
		predicate bool
	}{
		{column: structuredSearchSubject, value: filter.Subject},
		{column: structuredSearchPredicate, value: filter.Predicate, predicate: true},
		{column: structuredSearchObject, value: filter.Object},
	}
	for _, field := range structured {
		if err := appendStructuredSearchFilter(&predicates, &args, field.column, field.value, field.predicate); err != nil {
			return builtSearchQuery{}, err
		}
	}

	args = append(args, limit+1)
	return builtSearchQuery{
		SQL:   cte + fmt.Sprintf(searchStatement, strings.Join(predicates, "\n    and ")),
		Args:  args,
		Terms: terms,
		Limit: limit,
	}, nil
}

func appendStructuredSearchFilter(predicates *[]string, args *[]any, column structuredSearchColumn, value opt.Optional[string], predicate bool) error {
	pattern, ok := value.Get()
	if !ok {
		return nil
	}
	if len(pattern) > searchStructuredValueMaxBytes {
		return fmt.Errorf("%w: structured filter exceeds %d bytes", ErrInvalid, searchStructuredValueMaxBytes)
	}
	pattern = normalisePattern(pattern, predicate)
	if pattern == "" {
		return fmt.Errorf("%w: structured filter is empty", ErrInvalid)
	}
	exactSQL, wildcardSQL, err := structuredSearchPredicates(column)
	if err != nil {
		return err
	}
	if strings.Contains(pattern, "*") {
		*predicates = append(*predicates, wildcardSQL)
		*args = append(*args, strings.ReplaceAll(escapeLike(pattern), "*", "%"))
	} else {
		*predicates = append(*predicates, exactSQL)
		*args = append(*args, pattern)
	}
	return nil
}

func structuredSearchPredicates(column structuredSearchColumn) (string, string, error) {
	switch column {
	case structuredSearchSubject:
		return "subject = ?", "subject like ? escape '\\'", nil
	case structuredSearchPredicate:
		return "predicate = ?", "predicate like ? escape '\\'", nil
	case structuredSearchObject:
		return "object = ?", "object like ? escape '\\'", nil
	default:
		return "", "", fmt.Errorf("%w: invalid structured search column", ErrInvalid)
	}
}

func normalisePattern(value string, predicate bool) string {
	parts := strings.Split(value, "*")
	for i, part := range parts {
		parts[i] = normaliseValue(part, predicate)
	}
	return strings.Join(parts, "*")
}

func escapeLike(value string) string {
	return strings.NewReplacer("\\", "\\\\", "%", "\\%", "_", "\\_").Replace(value)
}
