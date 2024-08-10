package weaviate_semdexer

import (
	"fmt"
	"testing"

	"github.com/Southclaws/storyden/internal/utils"
	"github.com/rs/xid"
)

func TestGetWeaviateID(t *testing.T) {
	fmt.Println(GetWeaviateID(utils.Must(xid.FromString("cn2h3gfljatbqvjqctdg"))))
	fmt.Println(GetWeaviateID(utils.Must(xid.FromString("cn2h3gfljatbqvjqctdg"))))
	fmt.Println(GetWeaviateID(utils.Must(xid.FromString("cn2h3gfljatbqvjqctdg"))))
}
