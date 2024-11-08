package tag

import "github.com/rs/xid"

//go:generate go run github.com/Southclaws/enumerator

type tagFillRuleEnum string

const (
	tagFillRuleEnumNone tagFillRuleEnum = "none"
	tagFillRuleQuery    tagFillRuleEnum = "query"
	tagFillRuleReplace  tagFillRuleEnum = "replace"
)

type TagFillCommand struct {
	TargetNodeID xid.ID
	FillRule     TagFillRule
}
