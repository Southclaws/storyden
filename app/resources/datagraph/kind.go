package datagraph

//go:generate go run -mod=mod github.com/Southclaws/enumerator

type kindEnum string

const (
	kindPost       kindEnum = "post"
	kindThread     kindEnum = "thread"
	kindReply      kindEnum = "reply"
	kindNode       kindEnum = "node"
	kindCollection kindEnum = "collection"
	kindProfile    kindEnum = "profile"
	kindEvent      kindEnum = "event"
)
