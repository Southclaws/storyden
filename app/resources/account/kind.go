package account

//go:generate go run github.com/Southclaws/enumerator

type accountKindEnum string

const (
	accountKindHuman accountKindEnum = "human"
	accountKindBot   accountKindEnum = "bot"
)
