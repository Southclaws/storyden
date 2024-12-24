package asset

import (
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
)

//go:generate go run github.com/Southclaws/enumerator

type contentFillRuleEnum string

const (
	contentFillRuleQuery   contentFillRuleEnum = "query"
	contentFillRuleCreate  contentFillRuleEnum = "create"
	contentFillRulePrepend contentFillRuleEnum = "prepend"
	contentFillRuleAppend  contentFillRuleEnum = "append"
	contentFillRuleReplace contentFillRuleEnum = "replace"
)

type fillSourceEnum string

const (
	fillSourceURL     fillSourceEnum = "url"
	fillSourceContent fillSourceEnum = "content"
)

type ContentFillCommand struct {
	TargetNodeID opt.Optional[xid.ID]
	SourceType   opt.Optional[FillSource]
	FillRule     ContentFillRule
}
