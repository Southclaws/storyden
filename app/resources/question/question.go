package question

import (
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

type Question struct {
	ID     xid.ID
	Slug   string
	Query  string
	Result datagraph.Content
	Author opt.Optional[account.Account]
}

func Map(in *ent.Question) (*Question, error) {
	authorEdge := opt.NewPtr(in.Edges.Author)

	result, err := datagraph.NewRichText(in.Result)
	if err != nil {
		return nil, err
	}

	author, err := opt.MapErr(authorEdge, func(a ent.Account) (account.Account, error) {
		acc, err := account.MapAccount(&a)
		if err != nil {
			return account.Account{}, err
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
