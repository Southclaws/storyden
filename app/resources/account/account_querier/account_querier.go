package account_querier

import "github.com/Southclaws/storyden/app/resources/account/account_repo"

type Querier struct {
	*account_repo.Repository
}

func New(repo *account_repo.Repository) *Querier {
	return &Querier{Repository: repo}
}
