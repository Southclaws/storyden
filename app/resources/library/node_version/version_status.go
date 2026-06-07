package node_version

//go:generate go run github.com/Southclaws/enumerator

type versionStatusEnum string

const (
	versionStatusDraft   versionStatusEnum = "draft"
	versionStatusApplied versionStatusEnum = "applied"
)
