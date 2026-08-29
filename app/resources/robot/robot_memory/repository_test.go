package robot_memory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_memory"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
)

func TestRepositoryHierarchySearchStateAndAccess(t *testing.T) {
	t.Parallel()
	robotRef := "memory-hierarchy-" + xid.New().String()
	otherRobotRef := "memory-hierarchy-other-" + xid.New().String()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		db *ent.Client,
		raw *sqlx.DB,
		repo *robot_memory.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			people, err := repo.Create(ctx, robotRef, opt.NewEmpty[robot.MemoryID](), "Shared information about community members")
			require.NoError(t, err)
			projects, err := repo.Create(ctx, robotRef, opt.NewEmpty[robot.MemoryID](), "Durable project context")
			require.NoError(t, err)
			person, err := repo.Create(ctx, robotRef, opt.New(people.ID), "SOUTHCLAWS uses the 100%_literal token for portable search tests")
			require.NoError(t, err)
			listening, err := repo.Create(ctx, robotRef, opt.New(person.ID), "Southclaws listens to atmospheric ambient music")
			require.NoError(t, err)
			other, err := repo.Create(ctx, otherRobotRef, opt.NewEmpty[robot.MemoryID](), "Southclaws belongs to another Robot")
			require.NoError(t, err)

			freshRepo := robot_memory.New(db, raw)
			rootItems, hasMore, err := freshRepo.List(ctx, robotRef, opt.NewEmpty[robot.MemoryID](), robot_memory.ListLimit)
			require.NoError(t, err)
			assert.False(t, hasMore)
			require.Len(t, rootItems, 2)
			assert.ElementsMatch(t, []string{"Shared information about community members", "Durable project context"}, summaryExcerpts(rootItems))
			for _, item := range rootItems {
				if item.Memory.ID == people.ID {
					assert.Equal(t, 1, item.Children)
				}
			}

			_, err = repo.Create(ctx, otherRobotRef, opt.New(people.ID), "Cross-Robot parent")
			require.Error(t, err)
			_, err = repo.Move(ctx, robotRef, people.ID, opt.New(listening.ID))
			require.ErrorIs(t, err, robot_memory.ErrHierarchyCycle)
			_, err = repo.Move(ctx, otherRobotRef, projects.ID, opt.NewEmpty[robot.MemoryID]())
			require.Error(t, err)

			moved, err := repo.Move(ctx, robotRef, projects.ID, opt.New(people.ID))
			require.NoError(t, err)
			assert.Equal(t, people.ID, moved.ParentID.OrZero())

			results, hasMore, err := repo.Search(ctx, robotRef, robot_memory.SearchFilter{Query: "southclaws token"}, opt.NewEmpty[robot.MemoryID](), robot_memory.SearchLimit)
			require.NoError(t, err)
			assert.False(t, hasMore)
			require.Len(t, results, 1)
			assert.Equal(t, person.ID, results[0].Memory.ID)
			assert.Equal(t, []robot.MemoryID{people.ID, person.ID}, itemIDs(results[0].Path))

			results, _, err = repo.Search(ctx, robotRef, robot_memory.SearchFilter{Query: "%_literal"}, opt.NewEmpty[robot.MemoryID](), robot_memory.SearchLimit)
			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, person.ID, results[0].Memory.ID)

			results, _, err = repo.Search(ctx, robotRef, robot_memory.SearchFilter{Query: "southclaws"}, opt.New(person.ID), robot_memory.SearchLimit)
			require.NoError(t, err)
			assert.ElementsMatch(t, []robot.MemoryID{person.ID, listening.ID}, searchIDs(results))

			results, _, err = repo.Search(ctx, otherRobotRef, robot_memory.SearchFilter{Query: "southclaws"}, opt.NewEmpty[robot.MemoryID](), robot_memory.SearchLimit)
			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, other.ID, results[0].Memory.ID)

			accessed, err := repo.Create(ctx, robotRef, opt.NewEmpty[robot.MemoryID](), "atomic access counter probe")
			require.NoError(t, err)
			assert.Equal(t, uint64(0), accessed.AccessCount)
			results, _, err = repo.Search(ctx, robotRef, robot_memory.SearchFilter{Query: "atomic probe"}, opt.NewEmpty[robot.MemoryID](), robot_memory.SearchLimit)
			require.NoError(t, err)
			require.Len(t, results, 1)
			opened, err := repo.Open(ctx, robotRef, accessed.ID)
			require.NoError(t, err)
			assert.Equal(t, uint64(2), opened.Memory.AccessCount)

			_, err = repo.Update(ctx, robotRef, person.ID, robot_memory.WithState(robot.MemoryStateArchived))
			require.NoError(t, err)
			results, _, err = repo.Search(ctx, robotRef, robot_memory.SearchFilter{Query: "literal token"}, opt.NewEmpty[robot.MemoryID](), robot_memory.SearchLimit)
			require.NoError(t, err)
			assert.Empty(t, results)

			listed, _, err := repo.List(ctx, robotRef, opt.New(people.ID), robot_memory.ListLimit)
			require.NoError(t, err)
			for _, item := range listed {
				if item.Memory.ID == person.ID {
					assert.Equal(t, robot.MemoryStateArchived, item.Memory.State)
				}
			}

			results, _, err = repo.Search(ctx, robotRef+"'; drop table robot_memories; --", robot_memory.SearchFilter{
				Query:   "query' or 1=1; drop table robot_memories; -- %_\\ ?",
				Subject: opt.New("*subject' union select password from accounts --*"),
			}, opt.NewEmpty[robot.MemoryID](), robot_memory.SearchLimit)
			require.NoError(t, err)
			assert.Empty(t, results)

			_, err = repo.Create(ctx, robotRef, opt.NewEmpty[robot.MemoryID](), "The memory table remains available after hostile search input.")
			require.NoError(t, err)
		}))
	}))
}

func TestRepositoryStructuredFieldsAndLimits(t *testing.T) {
	t.Parallel()
	robotRef := "memory-structured-" + xid.New().String()
	limitRobotRef := "memory-limit-" + xid.New().String()
	listRobotRef := "memory-list-limit-" + xid.New().String()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		repo *robot_memory.Repository,
	) {
		lc.Append(fx.StartHook(func() {
			memory, err := repo.Create(
				ctx,
				robotRef,
				opt.NewEmpty[robot.MemoryID](),
				"Southclaws listens to atmospheric ambient music.",
				robot_memory.WithFact("  Southclaws  ", " Listens - To ", " Atmospheric   Ambient "),
			)
			require.NoError(t, err)
			fact, ok := memory.Fact.Get()
			require.True(t, ok)
			assert.Equal(t, robot.MemoryFact{Subject: "southclaws", Predicate: "listens_to", Object: "atmospheric ambient"}, fact)

			duplicate, err := repo.Create(
				ctx,
				robotRef,
				opt.NewEmpty[robot.MemoryID](),
				"A second piece of evidence for the same assertion.",
				robot_memory.WithFact("SOUTHCLAWS", "listens-to", "atmospheric ambient"),
			)
			require.NoError(t, err)
			assert.NotEqual(t, memory.ID, duplicate.ID)

			results, hasMore, err := repo.Search(ctx, robotRef, robot_memory.SearchFilter{
				Subject: opt.New("*claw*"), Predicate: opt.New("listens to"),
			}, opt.NewEmpty[robot.MemoryID](), robot_memory.SearchLimit)
			require.NoError(t, err)
			assert.False(t, hasMore)
			assert.ElementsMatch(t, []robot.MemoryID{memory.ID, duplicate.ID}, searchIDs(results))

			literal, err := repo.Create(
				ctx,
				robotRef,
				opt.NewEmpty[robot.MemoryID](),
				"Literal structured wildcard probe.",
				robot_memory.WithFact("100%_literal", "is-a", "search probe"),
			)
			require.NoError(t, err)
			results, _, err = repo.Search(ctx, robotRef, robot_memory.SearchFilter{Subject: opt.New("*%_*")}, opt.NewEmpty[robot.MemoryID](), robot_memory.SearchLimit)
			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, literal.ID, results[0].Memory.ID)

			cleared, err := repo.Update(ctx, robotRef, literal.ID, robot_memory.WithoutFact())
			require.NoError(t, err)
			assert.False(t, cleared.Fact.Ok())

			_, err = repo.Update(ctx, robotRef, memory.ID, robot_memory.WithState(robot.MemoryStateArchived))
			require.NoError(t, err)
			results, _, err = repo.Search(ctx, robotRef, robot_memory.SearchFilter{Subject: opt.New("southclaws")}, opt.NewEmpty[robot.MemoryID](), robot_memory.SearchLimit)
			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, duplicate.ID, results[0].Memory.ID)

			for i := 0; i < robot_memory.SearchLimit+1; i++ {
				_, err := repo.Create(ctx, limitRobotRef, opt.NewEmpty[robot.MemoryID](), fmt.Sprintf("bounded result probe %02d", i))
				require.NoError(t, err)
			}
			results, hasMore, err = repo.Search(ctx, limitRobotRef, robot_memory.SearchFilter{Query: "bounded probe"}, opt.NewEmpty[robot.MemoryID](), robot_memory.SearchLimit)
			require.NoError(t, err)
			assert.Len(t, results, robot_memory.SearchLimit)
			assert.True(t, hasMore)
			for i := 1; i < len(results); i++ {
				previous, current := results[i-1], results[i]
				assert.True(t, previous.Memory.UpdatedAt.After(current.Memory.UpdatedAt) || previous.Memory.UpdatedAt.Equal(current.Memory.UpdatedAt) && previous.Memory.ID.String() < current.Memory.ID.String())
			}

			for i := 0; i < 26; i++ {
				_, err := repo.Create(ctx, listRobotRef, opt.NewEmpty[robot.MemoryID](), fmt.Sprintf("list result probe %02d", i))
				require.NoError(t, err)
			}
			listed, listHasMore, err := repo.List(ctx, listRobotRef, opt.NewEmpty[robot.MemoryID](), robot_memory.ListLimit)
			require.NoError(t, err)
			assert.Len(t, listed, 25)
			assert.True(t, listHasMore)
		}))
	}))
}

func summaryExcerpts(items []*robot_memory.Item) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = item.Excerpt
	}
	return result
}

func itemIDs(items []*robot_memory.Item) []robot.MemoryID {
	result := make([]robot.MemoryID, len(items))
	for i, item := range items {
		result[i] = item.Memory.ID
	}
	return result
}

func searchIDs(items []*robot_memory.SearchResult) []robot.MemoryID {
	result := make([]robot.MemoryID, len(items))
	for i, item := range items {
		result[i] = item.Memory.ID
	}
	return result
}
