package pinecone_semdexer

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func Test_generateChunkID(t *testing.T) {
	a := assert.New(t)

	id := xid.New()

	id1 := generateChunkID(id, "chunk number one")
	id2 := generateChunkID(id, "chunk number one")
	id3 := generateChunkID(id, "chunk number two")

	a.Equal(id1, id2)
	a.NotEqual(id1, id3)
	a.NotEqual(id2, id3)
}

func Test_generateID(t *testing.T) {
	a := assert.New(t)

	chunk := `org

**Name-** Fake News Inference Dataset

**Link-** [https://ieee-dataport.org/open-access/fnid-fake-news-inference-dataset](https://ieee-dataport.org/open-access/fnid-fake-news-inference-dataset)

This database is provided for the Fake News Detection task`

	id, _ := xid.FromString("cth0hcifunp6ib5ivvug")

	id1 := generateChunkID(id, chunk)

	a.Equal("cth0hcifunp6ib5ivvug/97d726b5-5092-49af-8eb7-fa00d1be9b8d", id1)
}
