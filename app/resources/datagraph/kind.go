package datagraph

//go:generate go run -mod=mod github.com/Southclaws/enumerator

type kindEnum string

const (
	kindThread  kindEnum = "thread"
	kindReply   kindEnum = "reply"
	kindCluster kindEnum = "cluster"
	kindItem    kindEnum = "item"
	kindLink    kindEnum = "link"
)
