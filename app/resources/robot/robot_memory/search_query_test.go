package robot_memory

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/robot"
)

func TestBuildSearchQueryGlobal(t *testing.T) {
	t.Parallel()

	built, err := buildSearchQuery(
		"denbot",
		SearchFilter{Query: "  Southclaws   100%_Literal  "},
		opt.NewEmpty[robot.MemoryID](),
		7,
	)
	require.NoError(t, err)

	expectedSQL := `select
    id,
    content
from
    robot_memories
where
    robot_ref = ?
    and state = 'active'
    and lower(content) like ? escape '\'
    and lower(content) like ? escape '\'
order by
    updated_at desc,
    id asc
limit ?
`
	assert.Equal(t, expectedSQL, built.SQL)
	assert.Equal(t, []any{"denbot", "%southclaws%", "%100\\%\\_literal%", 8}, built.Args)
	assert.Equal(t, []string{"southclaws", "100%_literal"}, built.Terms)
	assert.Equal(t, 7, built.Limit)
}

func TestBuildSearchQuerySubtreeAndStructuredFilters(t *testing.T) {
	t.Parallel()

	parentID := robot.MemoryID(xid.ID{1, 2, 3})
	built, err := buildSearchQuery(
		"denbot",
		SearchFilter{
			Subject:   opt.New("  Southclaws  "),
			Predicate: opt.New(" Listens - To* "),
			Object:    opt.New(" Atmospheric   Ambient "),
		},
		opt.New(parentID),
		SearchLimit+100,
	)
	require.NoError(t, err)

	expectedSQL := `with recursive subtree (id) as (
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
select
    id,
    content
from
    robot_memories
where
    robot_ref = ?
    and state = 'active'
    and id in (select id from subtree)
    and subject = ?
    and predicate like ? escape '\'
    and object = ?
order by
    updated_at desc,
    id asc
limit ?
`
	assert.Equal(t, expectedSQL, built.SQL)
	assert.Equal(t, []any{
		"denbot",
		parentID.String(),
		"denbot",
		"denbot",
		"southclaws",
		"listens\\_to%",
		"atmospheric ambient",
		SearchLimit + 1,
	}, built.Args)
	assert.Empty(t, built.Terms)
	assert.Equal(t, SearchLimit, built.Limit)
}

func TestBuildSearchQueryNeverInterpolatesInput(t *testing.T) {
	t.Parallel()

	robotRef := "robot_payload'; DROP TABLE robot_memories; -- ? $1"
	filter := SearchFilter{
		Query:     "query_payload' OR 1=1; DROP TABLE accounts; -- %_\\ ? $2",
		Subject:   opt.New("subject_payload' UNION SELECT password FROM accounts --"),
		Predicate: opt.New("predicate-payload*' OR TRUE --"),
		Object:    opt.New("object_payload'); DELETE FROM robot_memories; --"),
	}
	built, err := buildSearchQuery(robotRef, filter, opt.New(robot.MemoryID(xid.ID{9, 8, 7})), SearchLimit)
	require.NoError(t, err)

	boundValues := strings.NewReplacer("\\", "", "%", "").Replace(strings.ToLower(fmt.Sprint(built.Args)))
	for _, marker := range []string{
		"robot_payload",
		"query_payload",
		"subject_payload",
		"predicate_payload",
		"object_payload",
		"drop table",
		"delete from",
		"union select",
		"or 1=1",
	} {
		assert.NotContains(t, strings.ToLower(built.SQL), marker)
		assert.Contains(t, boundValues, marker)
	}
	assert.Equal(t, strings.Count(built.SQL, "?"), len(built.Args))

	postgresSQL := sqlx.Rebind(sqlx.DOLLAR, built.SQL)
	assert.NotContains(t, postgresSQL, robotRef)
	assert.Contains(t, postgresSQL, "$1")
	assert.Contains(t, postgresSQL, fmt.Sprintf("$%d", len(built.Args)))
}

func TestBuildSearchQueryRejectsResourceExhaustionInputs(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		robotRef string
		filter   SearchFilter
	}{
		"empty Robot reference": {
			filter: SearchFilter{},
		},
		"oversized Robot reference": {
			robotRef: strings.Repeat("r", searchRobotRefMaxBytes+1),
			filter:   SearchFilter{},
		},
		"oversized text query": {
			robotRef: "denbot",
			filter:   SearchFilter{Query: strings.Repeat("q", searchQueryMaxBytes+1)},
		},
		"too many text terms": {
			robotRef: "denbot",
			filter:   SearchFilter{Query: strings.TrimSpace(strings.Repeat("term ", searchTermLimit+1))},
		},
		"oversized subject": {
			robotRef: "denbot",
			filter:   SearchFilter{Subject: opt.New(strings.Repeat("s", searchStructuredValueMaxBytes+1))},
		},
		"oversized predicate": {
			robotRef: "denbot",
			filter:   SearchFilter{Predicate: opt.New(strings.Repeat("p", searchStructuredValueMaxBytes+1))},
		},
		"oversized object": {
			robotRef: "denbot",
			filter:   SearchFilter{Object: opt.New(strings.Repeat("o", searchStructuredValueMaxBytes+1))},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := buildSearchQuery(test.robotRef, test.filter, opt.NewEmpty[robot.MemoryID](), SearchLimit)
			require.ErrorIs(t, err, ErrInvalid)
		})
	}
}

func TestBuildSearchQueryAcceptsBoundaries(t *testing.T) {
	t.Parallel()

	query := strings.Repeat("q", searchQueryMaxBytes)
	structured := strings.Repeat("s", searchStructuredValueMaxBytes)
	built, err := buildSearchQuery(strings.Repeat("r", searchRobotRefMaxBytes), SearchFilter{
		Query: query, Subject: opt.New(structured), Predicate: opt.New(structured), Object: opt.New(structured),
	}, opt.NewEmpty[robot.MemoryID](), 0)
	require.NoError(t, err)
	assert.Equal(t, SearchLimit, built.Limit)
	assert.Equal(t, 6, len(built.Args))
	assert.Equal(t, strings.Count(built.SQL, "?"), len(built.Args))
}

func TestStructuredSearchPredicatesRejectUnknownColumn(t *testing.T) {
	t.Parallel()

	exact, wildcard, err := structuredSearchPredicates(structuredSearchColumn(255))
	require.ErrorIs(t, err, ErrInvalid)
	assert.Empty(t, exact)
	assert.Empty(t, wildcard)
}

func FuzzBuildSearchQuerySQLShape(f *testing.F) {
	f.Add(
		"robot'; drop table robot_memories; --",
		"query' or 1=1 -- %_\\",
		"subject' union select",
		"predicate-*",
		"object'); delete from robot_memories; --",
		uint8(15),
	)
	f.Add("denbot", "normal query", "", "", "", uint8(0))

	f.Fuzz(func(t *testing.T, robotRef, query, subject, predicate, object string, present uint8) {
		filter := SearchFilter{Query: query}
		if present&1 != 0 {
			filter.Subject = opt.New(subject)
		}
		if present&2 != 0 {
			filter.Predicate = opt.New(predicate)
		}
		if present&4 != 0 {
			filter.Object = opt.New(object)
		}
		parentID := opt.NewEmpty[robot.MemoryID]()
		if present&8 != 0 {
			parentID = opt.New(robot.MemoryID(xid.ID{1, 3, 3, 7}))
		}

		built, err := buildSearchQuery(robotRef, filter, parentID, SearchLimit)
		if err != nil {
			return
		}

		safeFilter := SearchFilter{Query: strings.TrimSpace(strings.Repeat("safe ", len(built.Terms)))}
		safeFilter.Subject = equivalentSearchPattern(filter.Subject)
		safeFilter.Predicate = equivalentSearchPattern(filter.Predicate)
		safeFilter.Object = equivalentSearchPattern(filter.Object)
		safe, err := buildSearchQuery("safe-robot", safeFilter, parentID, SearchLimit)
		require.NoError(t, err)

		assert.Equal(t, safe.SQL, built.SQL)
		assert.Equal(t, len(safe.Args), len(built.Args))
		assert.Equal(t, strings.Count(built.SQL, "?"), len(built.Args))
	})
}

func equivalentSearchPattern(value opt.Optional[string]) opt.Optional[string] {
	pattern, ok := value.Get()
	if !ok {
		return opt.NewEmpty[string]()
	}
	if strings.Contains(pattern, "*") {
		return opt.New("safe*value")
	}
	return opt.New("safe")
}
