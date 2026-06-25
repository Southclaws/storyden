package model_ref

import (
	"fmt"
	"strings"
	"time"

	"github.com/Southclaws/opt"
)

type Provider struct{ string }

func NewProvider(s string) Provider {
	return Provider{string: s}
}

func (p Provider) String() string { return p.string }

type Model struct{ string }

func NewModel(s string) Model {
	return Model{string: s}
}

func (m Model) String() string { return m.string }

type ModelRef struct {
	Provider Provider
	Model    Model
}

func ParseID(s string) (ModelRef, error) {
	provider, name, found := strings.Cut(s, "/")
	if !found || provider == "" || name == "" {
		return ModelRef{}, fmt.Errorf("invalid model ModelRef %q: expected `provider/model-name`", s)
	}
	return ModelRef{Provider: Provider{provider}, Model: Model{name}}, nil
}

func (id ModelRef) String() string {
	return id.Provider.String() + "/" + id.Model.String()
}

type Info struct {
	Ref        ModelRef
	Raw        map[string]any
	LastSeenAt time.Time
}

func (m Info) Provider() Provider {
	return m.Ref.Provider
}

func (m Info) Model() Model {
	return m.Ref.Model
}

func (m Info) String() string {
	return m.Ref.String()
}

type CacheStatus struct {
	Provider        Provider
	LastRefreshedAt opt.Optional[time.Time]
	LastError       opt.Optional[string]
}
