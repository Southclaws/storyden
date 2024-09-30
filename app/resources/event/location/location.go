package location

import (
	"net/url"

	"github.com/Southclaws/opt"
)

//go:generate go run github.com/Southclaws/enumerator

type locationTypeEnum string

const (
	locationTypePhysical locationTypeEnum = `physical`
	locationTypeVirtual  locationTypeEnum = `virtual`
)

type Location interface {
	Type() LocationType
}

type Physical struct {
	Name      string
	Address   opt.Optional[string]
	Latitude  opt.Optional[float64]
	Longitude opt.Optional[float64]
	URL       opt.Optional[url.URL]
}

func (p *Physical) Type() LocationType {
	return LocationTypePhysical
}

type Virtual struct {
	Name string
	URL  opt.Optional[url.URL]
}

func (v *Virtual) Type() LocationType {
	return LocationTypeVirtual
}
