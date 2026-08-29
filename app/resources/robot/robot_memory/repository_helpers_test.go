package robot_memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormaliseValue(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "southclaws", normaliseValue("  SOUTHCLAWS  ", false))
	assert.Equal(t, "freyja is owned by southclaws", normaliseValue("\tFreyja\n is   owned by\r\nSouthclaws  ", false))
	assert.Equal(t, "freyja - southclaws", normaliseValue(" Freyja  -  Southclaws ", false))
	assert.Equal(t, "100%_literal", normaliseValue(" 100%_LITERAL ", false))
	assert.Equal(t, "freyja’s café", normaliseValue("  FREYJA’S CAFÉ  ", false))
	assert.Equal(t, "", normaliseValue(" \t\n ", false))
}

func TestNormalisePredicateValue(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "owned_by", normaliseValue("  Owned By  ", true))
	assert.Equal(t, "owned_by", normaliseValue("  Owned - By  ", true))
	assert.Equal(t, "listens_to", normaliseValue("-- listens---to --", true))
	assert.Equal(t, "is_already_normal", normaliseValue(" IS_ALREADY-NORMAL ", true))
	assert.Equal(t, "lives_with_southclaws", normaliseValue(" lives - with   Southclaws ", true))
	assert.Equal(t, "", normaliseValue(" - -- - ", true))
}

func TestNormalisePattern(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "south*claws", normalisePattern("  SOUTH * CLAWS  ", false))
	assert.Equal(t, "*southclaws*", normalisePattern("* SOUTHCLAWS *", false))
	assert.Equal(t, "*listens_to*", normalisePattern("* Listens - To *", true))
	assert.Equal(t, "owned_by***lives_with", normalisePattern(" Owned By *** Lives-With ", true))
	assert.Equal(t, "100%_literal*", normalisePattern(" 100%_LITERAL * ", false))
	assert.Equal(t, "*", normalisePattern("*", true))
}

func TestEscapeLike(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "southclaws", escapeLike("southclaws"))
	assert.Equal(t, `100\%\_literal`, escapeLike("100%_literal"))
	assert.Equal(t, `path\\memory`, escapeLike(`path\memory`))
	assert.Equal(t, `*100\%\_literal*`, escapeLike("*100%_literal*"))
	assert.Equal(t, `\%\_\\`, escapeLike(`%_\`))
	assert.Equal(t, "", escapeLike(""))
}

func TestMakeExcerpt(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "Freyja is owned by Southclaws.", makeExcerpt("  Freyja\n is   owned by Southclaws.  ", nil, 100))
	assert.Equal(t, "abcdefghij…", makeExcerpt("abcdefghijklmnopqrstuvwxyz", nil, 10))
	assert.Equal(t, "…89 TARGET …", makeExcerpt("0123456789 TARGET abcdefghij", []string{"target"}, 10))
	assert.Equal(t, "…δεζηθι…", makeExcerpt("αβγδεζηθικλμν", []string{"ζη"}, 6))
	assert.Equal(t, "…جيم دال هاء …", makeExcerpt("ألف باء جيم دال هاء واو", []string{"دال"}, 12))
	assert.Equal(t, "…记忆芙蕾雅属于南…", makeExcerpt("前言记忆芙蕾雅属于南爪后记", []string{"芙蕾雅"}, 8))
	assert.Equal(t, "…ヤは猫です後書き", makeExcerpt("序文フレイヤは猫です後書き", []string{"猫"}, 8))
}
