package mark

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
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
