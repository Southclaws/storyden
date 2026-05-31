package mark

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMark(t *testing.T) {
	a := assert.New(t)

	type args struct {
		in   string
		slug string
		id   string
	}

	check := func(args args) {
		m := NewQueryKey(args.in)

		a.Equal(args.in, m.raw)
		a.Equal(args.slug, m.slug.String())
		a.Equal(args.id, m.id.String())
	}

	t.Run("just an xid", func(t *testing.T) {
		check(args{in: "cpvf89ifunp0qr2aqkog", id: "cpvf89ifunp0qr2aqkog"})
	})

	t.Run("xid and slug", func(t *testing.T) {
		check(args{in: "cpvf89ifunp0qr2aqkog-some-super-neat-post", id: "cpvf89ifunp0qr2aqkog", slug: "some-super-neat-post"})
	})

	t.Run("just a slug", func(t *testing.T) {
		check(args{in: "some-super-neat-post", slug: "some-super-neat-post"})
	})

	t.Run("looks like an xid but is not", func(t *testing.T) {
		check(args{in: "thisisnotanidfrnocap", slug: "thisisnotanidfrnocap"})
	})

	t.Run("non xid trailing hyphen", func(t *testing.T) {
		check(args{in: "thisisnotanidfrnocap-", slug: "thisisnotanidfrnocap-"})
	})

	t.Run("almost xid missing one char", func(t *testing.T) {
		check(args{in: "cpvf89ifunp0qr2aqko", slug: "cpvf89ifunp0qr2aqko"})
	})

	t.Run("valid xid trailing hyphen", func(t *testing.T) {
		check(args{in: "cpvf89ifunp0qr2aqkog-", id: "cpvf89ifunp0qr2aqkog"})
	})

	t.Run("invalid xid but looks like a well formed mark", func(t *testing.T) {
		check(args{in: "thisisnotanidfrnocap-some-slug", slug: "thisisnotanidfrnocap-some-slug"})
	})

	t.Run("long string without hyphens", func(t *testing.T) {
		check(args{in: "thisisareallylongstringwithoutahyphen", slug: "thisisareallylongstringwithoutahyphen"})
	})

	t.Run("cyrillic slug with 20 bytes but 10 runes", func(t *testing.T) {
		check(args{in: "единорогов", slug: "единорогов"})
	})

	t.Run("fuzz plain xid", func(t *testing.T) {
		for range 10000 {
			id := xid.New().String()
			check(args{in: id, id: id})
		}
	})

	t.Run("fuzz xid with slug", func(t *testing.T) {
		for range 10000 {
			id := xid.New().String()
			slug := xid.New().String()
			check(args{in: id + "-" + slug, id: id, slug: slug})
		}
	})

	t.Run("fuzz badly formed xid", func(t *testing.T) {
		for range 10000 {
			id := xid.New().String()
			slug := xid.New().String()
			input := id + slug
			check(args{in: input, slug: input})
		}
	})
}

func TestQueryableEqual(t *testing.T) {
	r := require.New(t)

	idA := xid.New()
	idB := xid.New()

	t.Run("same id, same slug — equal", func(t *testing.T) {
		a := NewQueryKey(idA.String() + "-foo")
		b := NewQueryKey(idA.String() + "-foo")
		r.True(a.Equal(b))
	})

	t.Run("same id, different slug — equal (id is canonical)", func(t *testing.T) {
		a := NewQueryKey(idA.String() + "-foo")
		b := NewQueryKey(idA.String() + "-bar")
		r.True(a.Equal(b))
	})

	t.Run("different ids — not equal", func(t *testing.T) {
		a := NewQueryKeyID(idA)
		b := NewQueryKeyID(idB)
		r.False(a.Equal(b))
	})

	t.Run("same slug only — equal", func(t *testing.T) {
		a := NewQueryKey("hello")
		b := NewQueryKey("hello")
		r.True(a.Equal(b))
	})

	t.Run("different slugs only — not equal", func(t *testing.T) {
		a := NewQueryKey("hello")
		b := NewQueryKey("world")
		r.False(a.Equal(b))
	})

	t.Run("id-only vs slug-only — not equal (cannot establish)", func(t *testing.T) {
		// Previously this returned true and caused 'cannot relate a node to
		// itself' errors when moving a node by slug to a parent referenced by
		// xid (or vice versa).
		a := NewQueryKeyID(idA)
		b := NewQueryKey("some-slug")
		r.False(a.Equal(b))
		r.False(b.Equal(a))
	})
}
