package asset

import "github.com/rs/xid"

//go:generate go run github.com/Southclaws/enumerator

type contentFillRuleEnum string

const (
	contentFillRuleEnumNone contentFillRuleEnum = "none"
	contentFillRulePrepend  contentFillRuleEnum = "prepend"
	contentFillRuleAppend   contentFillRuleEnum = "append"
	contentFillRuleReplace  contentFillRuleEnum = "replace"
)

type ContentFillCommand struct {
	TargetNodeID xid.ID
	FillRule     ContentFillRule
}
