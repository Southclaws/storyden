package question

import (
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/rs/xid"
)

type Question struct {
	ID     xid.ID
	Slug   string
	Query  string
	Result datagraph.Content
	Author account.Account
}

func Map(in *ent.Question) (*Question, error) {
	authorEdge, err := in.Edges.AuthorOrErr()
	if err != nil {
		return nil, err
	}

	result, err := datagraph.NewRichText(in.Result)
	if err != nil {
		return nil, err
	}

	author, err := account.MapAccount(authorEdge)
	if err != nil {
		return nil, err
	}

	return &Question{
		ID:     in.ID,
		Slug:   in.Slug,
		Query:  in.Query,
		Result: result,
		Author: *author,
	}, nil
}
