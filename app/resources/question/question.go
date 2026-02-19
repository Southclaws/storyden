package question

import (
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
)

type Question struct {
	ID     xid.ID
	Slug   string
	Query  string
	Result datagraph.Content
	Author opt.Optional[profile.Ref]
}

func Map(in *ent.Question, roleHydratorFn func(accID xid.ID) (held.Roles, error)) (*Question, error) {
	profileMapper := profile.RefMapper(roleHydratorFn)

	authorEdge := opt.NewPtr(in.Edges.Author)

	result, err := datagraph.NewRichText(in.Result)
	if err != nil {
		return nil, err
	}

	author, err := opt.MapErr(authorEdge, func(a ent.Account) (profile.Ref, error) {
		acc, err := profileMapper(&a)
		if err != nil {
			return profile.Ref{}, err
		}
		return *acc, nil
	})
	if err != nil {
		return nil, err
	}

	return &Question{
		ID:     in.ID,
		Slug:   in.Slug,
		Query:  in.Query,
		Result: result,
		Author: author,
	}, nil
}
