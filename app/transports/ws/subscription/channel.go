package subscription

//go:generate go run github.com/Southclaws/enumerator

type channelEnum string

const (
	channelNone   channelEnum = "none"
	channelThread channelEnum = "thread"
)
