package datagraph

//go:generate go run -mod=mod github.com/Southclaws/enumerator

type kindEnum string

const (
	kindPost kindEnum = "post"
	kindNode kindEnum = "node"
)
