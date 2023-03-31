package settings

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"

	"github.com/Southclaws/storyden/internal/ent"
)

func Test_fromEnt(t *testing.T) {
	in := []*ent.Setting{
		{ID: "Title", Value: "Storyden"},
		{ID: "Description", Value: "A forum for the modern age."},
	}

	out, err := fromEnt(in)

	fmt.Println(err)
	pretty.Println(out)
}
