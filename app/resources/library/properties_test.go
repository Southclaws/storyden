package library

import (
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPropertySchema_Split(t *testing.T) {
	t.Run("all_existing", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		f0 := xid.New()
		f1 := xid.New()
		f2 := xid.New()

		schema := PropertySchema{
			Fields: PropertySchemaFields{
				{ID: f0, Name: "name", Type: PropertyTypeEnumText, Sort: "1"},
				{ID: f1, Name: "age", Type: PropertyTypeEnumNumber, Sort: "2"},
				{ID: f2, Name: "dob", Type: PropertyTypeEnumTimestamp, Sort: "3"},
			},
		}
		pml := PropertyMutationList{
			{ID: opt.New(f0), Name: "name", Value: "John"},
			{ID: opt.New(f1), Name: "age", Value: "25"},
			{ID: opt.New(f2), Name: "dob", Value: "2025-01-01T12:59:21Z"},
		}

		mut, err := schema.Split(pml)
		r.NoError(err)
		a.Len(mut.NewProps, 0)
		a.Len(mut.ExistingProps, 3)
		a.Len(mut.RemovedProps, 0)
	})

	t.Run("some_existing_one_removed", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		f0 := xid.New()
		f1 := xid.New()
		f2 := xid.New()

		schema := PropertySchema{
			Fields: PropertySchemaFields{
				{ID: f0, Name: "name", Type: PropertyTypeEnumText, Sort: "1"},
				{ID: f1, Name: "age", Type: PropertyTypeEnumNumber, Sort: "2"},
				{ID: f2, Name: "dob", Type: PropertyTypeEnumTimestamp, Sort: "3"},
			},
		}
		pml := PropertyMutationList{
			{ID: opt.New(f0), Name: "name", Value: "John"},
			{ID: opt.New(f1), Name: "age", Value: "25"},
		}

		mut, err := schema.Split(pml)
		r.NoError(err)
		a.Len(mut.NewProps, 0)
		a.Len(mut.ExistingProps, 2)
		a.Len(mut.RemovedProps, 1)
	})

	t.Run("all_removed", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		schema := PropertySchema{
			Fields: PropertySchemaFields{
				{Name: "name", Type: PropertyTypeEnumText, Sort: "1"},
				{Name: "age", Type: PropertyTypeEnumNumber, Sort: "2"},
				{Name: "dob", Type: PropertyTypeEnumTimestamp, Sort: "3"},
			},
		}
		pml := PropertyMutationList{}

		mut, err := schema.Split(pml)
		r.NoError(err)
		a.Len(mut.NewProps, 0)
		a.Len(mut.ExistingProps, 0)
		a.Len(mut.RemovedProps, 3)
	})

	t.Run("some_new", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		f0 := xid.New()
		f1 := xid.New()
		f2 := xid.New()

		schema := PropertySchema{
			Fields: PropertySchemaFields{
				{ID: f0, Name: "name", Type: PropertyTypeEnumText, Sort: "1"},
				{ID: f1, Name: "age", Type: PropertyTypeEnumNumber, Sort: "2"},
				{ID: f2, Name: "dob", Type: PropertyTypeEnumTimestamp, Sort: "3"},
			},
		}
		pml := PropertyMutationList{
			{ID: opt.New(f0), Name: "name", Value: "John"},
			{ID: opt.New(f1), Name: "age", Value: "25"},
			{ID: opt.New(f2), Name: "dob", Value: "2025-05-05T12:13:15Z"},
			{Name: "strength", Value: "69"},
		}

		mut, err := schema.Split(pml)
		r.NoError(err)
		a.Len(mut.NewProps, 1)
		a.Len(mut.ExistingProps, 3)
		a.Len(mut.RemovedProps, 0)
	})

	t.Run("all_new_replace_existing", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		schema := PropertySchema{
			Fields: PropertySchemaFields{
				{Name: "name", Type: PropertyTypeEnumText, Sort: "1"},
				{Name: "age", Type: PropertyTypeEnumNumber, Sort: "2"},
				{Name: "dob", Type: PropertyTypeEnumTimestamp, Sort: "3"},
			},
		}
		pml := PropertyMutationList{
			{Name: "strength", Value: "69"},
			{Name: "rarity", Value: "legendary"},
			{Name: "damage", Value: "420"},
		}

		mut, err := schema.Split(pml)
		r.NoError(err)
		a.Len(mut.NewProps, 3)
		a.Len(mut.ExistingProps, 0)
		a.Len(mut.RemovedProps, 3)
	})
}
