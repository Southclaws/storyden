package datagraph

import "github.com/rs/xid"

type Match struct {
	ID          xid.ID
	Kind        Kind
	Slug        string
	Name        string
	Description string
}

type MatchList []Match
