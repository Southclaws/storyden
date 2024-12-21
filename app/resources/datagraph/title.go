package datagraph

import "github.com/rs/xid"

//go:generate go run github.com/Southclaws/enumerator

type titleFillRuleEnum string

const (
	titleFillRuleQuery   titleFillRuleEnum = "query"
	titleFillRuleReplace titleFillRuleEnum = "replace"
)

type TitleFillCommand struct {
	TargetNodeID xid.ID
	FillRule     TitleFillRule
}
