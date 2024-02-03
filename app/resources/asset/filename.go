package asset

import (
	"fmt"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/rs/xid"
)

type Filename struct {
	id    xid.ID
	name  string
	hasID bool
}

func (f Filename) IsKnown() bool { return f.hasID }

func (a Filename) String() string {
	return a.name
}

func NewFilename(name string) Filename {
	return NewExistingFilename(xid.New(), name)
}

func NewExistingFilename(id xid.ID, name string) Filename {
	return Filename{
		id:    id,
		name:  formatFilename(id, name),
		hasID: true,
	}
}

func NewFilepathFilename(name string) Filename {
	return Filename{name: name}
}

func ParseAssetFilename(s string) (*Filename, error) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return nil, errInvalidFormat
	}

	id, err := xid.FromString(parts[0])
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Filename{
		id:   id,
		name: fmt.Sprintf("%s-%s", parts[0], parts[1]),
	}, nil
}

func formatFilename(id xid.ID, name string) string {
	if name == "" {
		return id.String()
	}
	return fmt.Sprintf("%s-%s", id.String(), name)
}
