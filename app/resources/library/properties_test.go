package library

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPropertySchema_Split(t *testing.T) {
	t.Run("all_existing", func(t *testing.T) {
		a := assert.New(t)

		schema := PropertySchema{
			Fields: PropertySchemaFields{
				{Name: "name", Type: "string", Sort: "1"},
				{Name: "age", Type: "number", Sort: "2"},
				{Name: "dob", Type: "timestamp", Sort: "3"},
			},
		}
		pml := PropertyMutationList{
			{Name: "name", Value: "John"},
			{Name: "age", Value: "25"},
			{Name: "dob", Value: "2025-01-01T12:59:21Z"},
		}

		n, e, r := schema.Split(pml)
		a.Len(n, 0)
		a.Len(e, 3)
		a.Len(r, 0)
	})

	t.Run("some_existing_one_removed", func(t *testing.T) {
		a := assert.New(t)

		schema := PropertySchema{
			Fields: PropertySchemaFields{
				{Name: "name", Type: "string", Sort: "1"},
				{Name: "age", Type: "number", Sort: "2"},
				{Name: "dob", Type: "timestamp", Sort: "3"},
			},
		}
		pml := PropertyMutationList{
			{Name: "name", Value: "John"},
			{Name: "age", Value: "25"},
		}

		n, e, r := schema.Split(pml)
		a.Len(n, 0)
		a.Len(e, 2)
		a.Len(r, 1)
	})

	t.Run("all_removed", func(t *testing.T) {
		a := assert.New(t)

		schema := PropertySchema{
			Fields: PropertySchemaFields{
				{Name: "name", Type: "string", Sort: "1"},
				{Name: "age", Type: "number", Sort: "2"},
				{Name: "dob", Type: "timestamp", Sort: "3"},
			},
		}
		pml := PropertyMutationList{}

		n, e, r := schema.Split(pml)
		a.Len(n, 0)
		a.Len(e, 0)
		a.Len(r, 3)
	})

	t.Run("some_new", func(t *testing.T) {
		a := assert.New(t)

		schema := PropertySchema{
			Fields: PropertySchemaFields{
				{Name: "name", Type: "string", Sort: "1"},
				{Name: "age", Type: "number", Sort: "2"},
				{Name: "dob", Type: "timestamp", Sort: "3"},
			},
		}
		pml := PropertyMutationList{
			{Name: "name", Value: "John"},
			{Name: "age", Value: "25"},
			{Name: "dob", Value: "2025-05-05T12:13:15Z"},
			{Name: "strength", Value: "69"},
		}

		n, e, r := schema.Split(pml)
		a.Len(n, 1)
		a.Len(e, 3)
		a.Len(r, 0)
	})

	t.Run("all_new_replace_existing", func(t *testing.T) {
		a := assert.New(t)

		schema := PropertySchema{
			Fields: PropertySchemaFields{
				{Name: "name", Type: "string", Sort: "1"},
				{Name: "age", Type: "number", Sort: "2"},
				{Name: "dob", Type: "timestamp", Sort: "3"},
			},
		}
		pml := PropertyMutationList{
			{Name: "strength", Value: "69"},
			{Name: "rarity", Value: "legendary"},
			{Name: "damage", Value: "420"},
		}

		n, e, r := schema.Split(pml)
		a.Len(n, 3)
		a.Len(e, 0)
		a.Len(r, 3)
	})
}
